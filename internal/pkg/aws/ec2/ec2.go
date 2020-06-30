package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type api interface {
	DescribeVpcs(*ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error)
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

// GetDefaultSubnetIDs finds the default subnet ids
func (c *EC2) DefaultSubnetIDs() ([]string, error) {
	response, err := c.client.DescribeSubnets(&ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("default-for-az"),
				Values: aws.StringSlice([]string{"true"}),
			},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("find default subnet ids: %w", err)
	}

	subnetIDs := make([]string, len(response.Subnets))
	for _, subnet := range response.Subnets {
		subnetIDs = append(subnetIDs, aws.StringValue(subnet.SubnetId))
	}
	return subnetIDs, nil
}

func (c *EC2) GetSubnetIDsFromAppEnv(app string, env string) ([]string, error) {
	response, err := c.client.DescribeSubnets(&ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:copilot-application"),
				Values: aws.StringSlice([]string{app}),
			},
			{
				Name:   aws.String("tag:copilot-environment"),
				Values: aws.StringSlice([]string{env}),
			},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("get subnet ids from environment: %w", err)
	}

	subnetIDs := make([]string, len(response.Subnets))
	for _, subnet := range response.Subnets {
		subnetIDs = append(subnetIDs, aws.StringValue(subnet.SubnetId))
	}
	return subnetIDs, nil
}

func (c *EC2) GetSecurityGroupsFromAppEnv(app string, env string) ([]string, error) {
	response, err := c.client.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:copilot-application"),
				Values: aws.StringSlice([]string{app}),
			},
			{
				Name:   aws.String("tag:copilot-environment"),
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
