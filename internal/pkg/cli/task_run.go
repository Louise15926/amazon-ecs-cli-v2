// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"errors"
	"fmt"

	"github.com/aws/copilot-cli/internal/pkg/aws/ecr"

	awscfn "github.com/aws/copilot-cli/internal/pkg/aws/cloudformation"

	"github.com/aws/copilot-cli/internal/pkg/deploy"
	"github.com/aws/copilot-cli/internal/pkg/term/log"

	"github.com/aws/copilot-cli/internal/pkg/aws/ecs"

	"github.com/aws/copilot-cli/internal/pkg/aws/ec2"
	"github.com/aws/copilot-cli/internal/pkg/aws/session"
	"github.com/aws/copilot-cli/internal/pkg/cli/selector"
	"github.com/aws/copilot-cli/internal/pkg/config"
	"github.com/aws/copilot-cli/internal/pkg/deploy/cloudformation"
	"github.com/aws/copilot-cli/internal/pkg/docker"
	"github.com/aws/copilot-cli/internal/pkg/term/color"
	"github.com/aws/copilot-cli/internal/pkg/term/prompt"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	errNumNotPositive = errors.New("number of tasks must be positive")
	errCpuNotPositive = errors.New("CPU units must be positive")
	errMemNotPositive = errors.New("memory must be positive")

	fmtTaskRunEnvPrompt       = fmt.Sprintf("In which %s would you like to run this %s?", color.Emphasize("environment"), color.Emphasize("task"))
	fmtTaskRunGroupNamePrompt = fmt.Sprintf("What would you like to %s your task group?", color.Emphasize("name"))

	taskRunEnvPromptHelp = fmt.Sprintf("Task will be deployed to the selected environment. "+
		"Select %s to run the task in your default VPC instead of any existing environment.", color.Emphasize(config.EnvNameNone))
	taskRunGroupNamePromptHelp = "The group name of the task. Tasks with the same group name share the same set of resources, including CloudFormation stack, CloudWatch log group, task definition and ECR repository."

	fmtImageURL = "%s:%s"
)

type runTaskVars struct {
	*GlobalOpts
	count  int64
	cpu    int
	memory int

	groupName string

	image          string
	dockerfilePath string
	imageTag       string

	taskRole string

	subnets        []string
	securityGroups []string
	env            string

	envVars map[string]string
	command string
}

type runTaskOpts struct {
	runTaskVars

	// Interfaces to interact with dependencies.
	fs     afero.Fs
	store  store
	parser dockerfileParser
	sel    appEnvWithNoneSelector

	docker dockerService
	ecr    ecrService
	ec2    vpcService
	ecs    *ecs.ECS
	cfn    cloudformation.CloudFormation
}

func newTaskRunOpts(vars runTaskVars) (*runTaskOpts, error) {
	store, err := config.NewStore()
	if err != nil {
		return nil, fmt.Errorf("new config store: %w", err)
	}

	sess, err := session.NewProvider().Default()
	return &runTaskOpts{
		runTaskVars: vars,

		fs:     &afero.Afero{Fs: afero.NewOsFs()},
		store:  store,
		sel:    selector.NewSelect(vars.prompt, store),
		docker: docker.New(),
		ec2:    ec2.New(sess),
		ecs:    ecs.New(sess),
		ecr:    ecr.New(sess),
		cfn:    cloudformation.New(sess),
	}, nil
}

// Validate returns an error if the flag values passed by the user are invalid.
func (o *runTaskOpts) Validate() error {
	if o.count <= 0 {
		return errNumNotPositive
	}

	if o.cpu <= 0 {
		return errCpuNotPositive
	}

	if o.memory <= 0 {
		return errMemNotPositive
	}

	if o.groupName != "" {
		if err := basicNameValidation(o.groupName); err != nil {
			return err
		}
	}

	if o.image != "" && o.dockerfilePath != "" {
		return errors.New("cannot specify both image and Dockerfile path")
	}

	if o.dockerfilePath != "" {
		if _, err := o.fs.Stat(o.dockerfilePath); err != nil {
			return err
		}
	}

	if o.env != "" && (o.subnets != nil || o.securityGroups != nil) {
		return errors.New("neither subnet nor security groups should be specified if environment is specified")
	}

	if o.appName != "" {
		if err := o.validateAppName(); err != nil {
			return err
		}
	}

	if o.env != "" {
		if err := o.validateEnvName(); err != nil {
			return err
		}
	}

	return nil
}

// Ask prompts the user for any required or important fields that are not provided.
func (o *runTaskOpts) Ask() error {
	if err := o.askTaskGroupName(); err != nil {
		return err
	}
	if err := o.askEnvName(); err != nil {
		return err
	}
	return nil
}

// Execute create or update task resources and run the task.
func (o *runTaskOpts) Execute() error {
	if o.image == "" && o.dockerfilePath == "" {
		o.dockerfilePath = "./Dockerfile"
	}

	if o.imageTag == "" {
		o.imageTag = "latest"
	}

	if err := o.getNetworkConfig(); err != nil {
		return err
	}

	if err := o.deployTaskResource(); err != nil {
		return err
	}

	//if image is not provided, then we build the image and push to ECR repo
	if o.image == "" {
		uri, err := o.pushToECRRepo()
		if err != nil {
			return err
		}
		o.image = fmt.Sprintf(fmtImageURL, uri, o.imageTag)

		// update image to stack
		if err := o.deployTaskResource(); err != nil {
			return err
		}
	}

	// TODO: kick off task
	return nil
}

func (o *runTaskOpts) getNetworkConfig() error {
	if o.env != config.EnvNameNone {
		subnets, err := o.ec2.GetSubnetIDsFromAppEnv(o.AppName(), o.env)
		if err != nil {
			return fmt.Errorf("get subnet ids: %w", err)
		}

		securityGroups, err := o.ec2.GetSecurityGroupsFromAppEnv(o.AppName(), o.env)
		if err != nil {
			return fmt.Errorf("get security groups: %w", err)
		}

		o.subnets = subnets
		o.securityGroups = securityGroups

		return nil
	}

	// get default subnet ids if not provided
	if o.subnets == nil {
		subnetIDs, err := o.ec2.GetDefaultSubnetIDs()
		if err != nil {
			return fmt.Errorf("get subnet ids: %w", err)
		}
		o.subnets = subnetIDs
	}
	return nil
}

func (o *runTaskOpts) deployTaskResource() error {
	if err := o.cfn.DeployTask(&deploy.CreateTaskResourcesInput{
		Name:     o.groupName,
		Cpu:      o.cpu,
		Memory:   o.memory,
		Image:    o.image,
		TaskRole: o.taskRole,
		Command:  o.command,
	}); err != nil {
		var errChangeSetEmpty *awscfn.ErrChangeSetEmpty
		if errors.As(err, &errChangeSetEmpty) {
			return nil
		}
		log.Errorf("failed to deploy resources for task: %w", err)
		return fmt.Errorf("deploy task: %w", err)
	}
	return nil
}

func (o *runTaskOpts) pushToECRRepo() (string, error) {
	repoName := fmt.Sprintf("%s-%s", "copilot", o.groupName)

	uri, err := o.ecr.GetRepository(repoName)
	if err != nil {
		return "", fmt.Errorf("get ECR repository URI: %w", err)
	}

	if err := o.docker.Build(uri, o.imageTag, o.dockerfilePath); err != nil {
		return "", fmt.Errorf("build Dockerfile at %s with tag %s: %w", o.dockerfilePath, o.imageTag, err)
	}

	auth, err := o.ecr.GetECRAuth()
	if err != nil {
		return "", fmt.Errorf("get ECR auth data: %w", err)
	}

	if err := o.docker.Login(uri, auth.Username, auth.Password); err != nil {
		return "", fmt.Errorf("login to repo %s: %w", repoName,  err)
	}

	if err := o.docker.Push(uri, o.imageTag); err != nil {
		return "", fmt.Errorf("push to repo: %w", err)
	}
	return uri, nil
}

func (o *runTaskOpts) validateAppName() error {
	if _, err := o.store.GetApplication(o.appName); err != nil {
		return fmt.Errorf("get application: %w", err)
	}
	return nil
}

func (o *runTaskOpts) validateEnvName() error {
	if o.AppName() != "" {
		if _, err := o.store.GetEnvironment(o.AppName(), o.env); err != nil {
			return fmt.Errorf("get environment: %w", err)
		}
	} else {
		return errNoAppInWorkspace
	}

	return nil
}

func (o *runTaskOpts) askTaskGroupName() error {
	if o.groupName != "" {
		return nil
	}

	// TODO: maybe list existing tasks like in ListApplications, ask whether to use existing tasks; require to implement task store first

	groupName, err := o.prompt.Get(
		fmtTaskRunGroupNamePrompt,
		taskRunGroupNamePromptHelp,
		basicNameValidation,
		prompt.WithFinalMessage("Task group name:"))
	if err != nil {
		return fmt.Errorf("prompt get task group name: %w", err)
	}
	o.groupName = groupName
	return nil
}

func (o *runTaskOpts) askEnvName() error {
	if o.env != "" {
		return nil
	}

	if o.AppName() == "" || o.subnets != nil {
		o.env = config.EnvNameNone
		return nil
	}

	env, err := o.sel.EnvironmentWithNone(fmtTaskRunEnvPrompt, taskRunEnvPromptHelp, o.AppName())
	if err != nil {
		return fmt.Errorf("ask for environment: %w", err)
	}
	o.env = env
	return nil
}

// BuildTaskRunCmd build the command for running a new task
func BuildTaskRunCmd() *cobra.Command {
	vars := runTaskVars{
		GlobalOpts: NewGlobalOpts(),
	}
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run a one-off task",
		Long:  `Run a one-off task with configurations such as cpu-units, memory, image, etc.`,
		Example: `
Run a task with default setting.
/code $ copilot task run
Run a task named "db-migrate" in the "test" environment under the current workspace.
/code $ copilot task run -n db-migrate --env test
Starts 4 tasks with 2GB memory, Runs a particular image.
/code $ copilot task run --num 4 --memory 2048 --task-role frontend-exec-role
Run a task with environment variables.
/code $ copilot task run --env-vars name=myName,user=myUser
`,
		RunE: runCmdE(func(cmd *cobra.Command, args []string) error {
			opts, err := newTaskRunOpts(vars)
			if err != nil {
				return err
			}
			if err := opts.Validate(); err != nil { // validate flags
				return err
			}

			if err := opts.Ask(); err != nil {
				return err
			}

			if err := opts.Execute(); err != nil {
				return err
			}
			return nil
		}),
	}

	cmd.Flags().Int64Var(&vars.count, countFlag, 1, countFlagDescription)
	cmd.Flags().IntVar(&vars.cpu, cpuFlag, 256, cpuFlagDescription)
	cmd.Flags().IntVar(&vars.memory, memoryFlag, 512, memoryFlagDescription)

	cmd.Flags().StringVarP(&vars.groupName, taskGroupNameFlag, nameFlagShort, "", taskGroupFlagDescription)

	cmd.Flags().StringVar(&vars.image, imageFlag, "", imageFlagDescription)
	cmd.Flags().StringVar(&vars.dockerfilePath, dockerFileFlag, "", dockerFileFlagDescription)
	cmd.Flags().StringVar(&vars.imageTag, imageTagFlag, "", imageFlagDescription)

	cmd.Flags().StringVar(&vars.taskRole, taskRoleFlag, "", taskRoleFlagDescription)

	cmd.Flags().StringVar(&vars.appName, appFlag, "", appFlagDescription)
	cmd.Flags().StringVar(&vars.env, envFlag, "", envFlagDescription)
	cmd.Flags().StringSliceVar(&vars.subnets, subnetFlag, nil, subnetFlagDescription)
	cmd.Flags().StringSliceVar(&vars.securityGroups, securityGroupsFlag, nil, securityGroupsFlagDescription)

	cmd.Flags().StringToStringVar(&vars.envVars, envVarsFlag, nil, envVarsFlagDescription)
	cmd.Flags().StringVar(&vars.command, commandsFlag, "", commandsFlagDescription)

	return cmd
}
