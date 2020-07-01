package ec2

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/aws/copilot-cli/internal/pkg/deploy/cloudformation/stack"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var (
	fmtFilterTagApp = fmt.Sprintf("tag:%s", stack.AppTagKey)
	fmtFilterTagEnv = fmt.Sprintf("tag:%s", stack.EnvTagKey)
)

const (
	fmtFilterDefault = "default-for-az"
)

type api interface {
	DescribeSubnets(*ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error)
	DescribeSecurityGroups(*ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error)
}

type EC2 struct {
	client api
}

func New(s *session.Session) *EC2 {
	return &EC2{
		client: ec2.New(s),
	}
}

// GetDefaultSubnetIDs finds the default subnet IDs
func (c *EC2) GetDefaultSubnetIDs() ([]string, error) {
	response, err := c.client.DescribeSubnets(&ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String(fmtFilterDefault),
				Values: aws.StringSlice([]string{"true"}),
			},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("find default subnet IDs: %w", err)
	}

	if len(response.Subnets) == 0 {
		return nil, errors.New("no default subnet ID found")
	}

	subnetIDs := make([]string, len(response.Subnets))
	for _, subnet := range response.Subnets {
		subnetIDs = append(subnetIDs, aws.StringValue(subnet.SubnetId))
	}
	return subnetIDs, nil
}

// GetDefaultSubnetIDs finds the subnet IDs associated with the environment of the application
func (c *EC2) GetSubnetIDs(app string, env string) ([]string, error) {
	response, err := c.client.DescribeSubnets(&ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String(fmtFilterTagApp),
				Values: aws.StringSlice([]string{app}),
			},
			{
				Name:   aws.String(fmtFilterTagEnv),
				Values: aws.StringSlice([]string{env}),
			},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("get subnet IDs from environment: %w", err)
	}

	if len(response.Subnets) == 0 {
		return nil, fmt.Errorf("no subnet id found for %s app %s env", app, env)
	}

	subnetIDs := make([]string, len(response.Subnets))
	for _, subnet := range response.Subnets {
		subnetIDs = append(subnetIDs, aws.StringValue(subnet.SubnetId))
	}
	return subnetIDs, nil
}

// GetDefaultSubnetIDs finds the security group IDs associated with the environment of the application
func (c *EC2) GetSecurityGroups(app string, env string) ([]string, error) {
	response, err := c.client.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String(fmtFilterTagApp),
				Values: aws.StringSlice([]string{app}),
			},
			{
				Name:   aws.String(fmtFilterTagEnv),
				Values: aws.StringSlice([]string{env}),
			},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("get security groups from environment: %w", err)
	}

	securityGroups := make([]string, len(response.SecurityGroups))
	for _, sg := range response.SecurityGroups {
		securityGroups = append(securityGroups, aws.StringValue(sg.GroupId))
	}
	return securityGroups, nil
}
