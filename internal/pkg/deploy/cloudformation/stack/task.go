package stack

import (
	"fmt"
	"github.com/aws/amazon-ecs-cli-v2/internal/pkg/template"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"strconv"
)

const (
	taskTemplatePath = "task/task.yml"
	taskLogRetention = "30"

	TaskNameParamKey           = "TaskName"
	TaskContainerImageParamKey = "ContainerImage"
	TaskCPUParamKey            = "TaskCPU"
	TaskMemoryParamKey         = "TaskMemory"
	TaskCountParamKey          = "TaskCount"
	TaskLogRetentionParamKey   = "LogRetention"
)

type task struct {
	name string

	count  int
	cpu    int
	memory int

	rc     RuntimeConfig
	parser template.ReadParser
}

func NewTaskStackConfig() *task {
	return &task{
		parser: template.New(),
	}
}

func (t *task) StackName() string {
	return NameForTask(t.name)
}

func (t *task) Template() (string, error) {
	content, err := t.parser.Read(taskTemplatePath)
	if err != nil {
		return "", fmt.Errorf("read template for task stack: %w", err)
	}
	return content.String(), nil
}

func (t *task) Parameters() ([]*cloudformation.Parameter, error) {
	return []*cloudformation.Parameter{
		{
			ParameterKey:   aws.String(TaskNameParamKey),
			ParameterValue: aws.String(t.name),
		},
		{
			ParameterKey:   aws.String(TaskContainerImageParamKey),
			ParameterValue: aws.String(fmt.Sprintf("%s:%s", t.rc.ImageRepoURL, t.rc.ImageTag)),
		},
		{
			ParameterKey:   aws.String(TaskCPUParamKey),
			ParameterValue: aws.String(strconv.Itoa(aws.IntValue(&t.cpu))),
		},
		{
			ParameterKey:   aws.String(TaskMemoryParamKey),
			ParameterValue: aws.String(strconv.Itoa(aws.IntValue(&t.memory))),
		},
		{
			ParameterKey:   aws.String(TaskCountParamKey),
			ParameterValue: aws.String(strconv.Itoa(t.count)),
		},
		{
			ParameterKey:   aws.String(TaskLogRetentionParamKey),
			ParameterValue: aws.String(taskLogRetention),
		},
	}, nil
}

func (t *task) Tags() []*cloudformation.Tag {
	return mergeAndFlattenTags(t.rc.AdditionalTags, map[string]string{})
}
