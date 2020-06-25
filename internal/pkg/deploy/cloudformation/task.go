package cloudformation

import (
    "errors"
    "github.com/aws/amazon-ecs-cli-v2/internal/pkg/aws/cloudformation"
)

// DeployTask deploys a task stack and waits until the deployment is done.
// If the task stack doesn't exist, then it creates the stack.
// If the task stack already exists, it updates the stack.
func (cf CloudFormation) DeployTask(conf StackConfiguration, opts ...cloudformation.StackOption) error {
    stack, err := toStack(conf)
    if err != nil {
        return err
    }

    for _, opt := range opts {
        opt(stack)
    }

    if err := cf.cfnClient.CreateAndWait(stack); err != nil {
        var errAlreadyExists *cloudformation.ErrStackAlreadyExists
        if !errors.As(err, &errAlreadyExists) {
            return err
        }
        return cf.cfnClient.UpdateAndWait(stack)
    }
    return nil
}
