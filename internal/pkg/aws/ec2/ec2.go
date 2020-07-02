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
	filterDefault = "default-for-az"
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
	filters := []*ec2.Filter{
		&ec2.Filter{
			Name:   aws.String(filterDefault),
			Values: aws.StringSlice([]string{"true"}),
		},
	}
	return c.getSubnetIDs(filters)
}

// GetDefaultSubnetIDs finds the subnet IDs associated with the environment of the application
func (c *EC2) GetSubnetIDs(app string, env string) ([]string, error) {
	filters := []*ec2.Filter{
		&ec2.Filter{
			Name:   aws.String(fmtFilterTagApp),
			Values: aws.StringSlice([]string{app}),
		},
		&ec2.Filter{
			Name:   aws.String(fmtFilterTagEnv),
			Values: aws.StringSlice([]string{env}),
		},
	}
	return c.getSubnetIDs(filters)
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

func (c *EC2) getSubnetIDs(filters []*ec2.Filter) ([]string, error) {
	response, err := c.client.DescribeSubnets(&ec2.DescribeSubnetsInput{
		Filters: filters,
	})

	if err != nil {
		return nil, fmt.Errorf("find subnet IDs: %w", err)
	}

	if len(response.Subnets) == 0 {
		return nil, errors.New("no subnet ID found")
	}

	subnetIDs := make([]string, len(response.Subnets))
	for idx, subnet := range response.Subnets {
		subnetIDs[idx] = aws.StringValue(subnet.SubnetId)
	}
	return subnetIDs, nil
}
