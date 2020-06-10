package ec2

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"strings"
)

type api interface {
	DescribeVpcs(*ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error)
	DescribeSubnets(*ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error)
	DescribeSecurityGroups(*ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error)
}

type EC2 struct {
	client api
}

func New(s *session.Session) EC2 {
	return EC2{
		client: ec2.New(s),
	}
}

func (c EC2) getDefaultVpc() (*ec2.Vpc, error) {
	response, err := c.client.DescribeVpcs(&ec2.DescribeVpcsInput{
		Filters: [] *ec2.Filter {
			{
				Name: aws.String("isDefault"),
				Values: aws.StringSlice([]string{"true"}),
			},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("find default VPC: %w", err)
	}

	foundVPCs := response.Vpcs

	if len(foundVPCs) <= 0 {
		return nil, errors.New("no default VPC found")
	}

	return foundVPCs[0], nil
}

// GetDefaultVpcId returns the ID of the default VPC, if it is available
func (c EC2) GetDefaultVpcId() (string, error) {
	vpc, err := c.getDefaultVpc()
	if err != nil {
		return "", fmt.Errorf("find default VPC's ID: %w", err)
	}
	
	if *vpc.State != ec2.VpcStateAvailable {
		return "", fmt.Errorf("find default VPC's ID: default VPC not available")
	}

	return *vpc.VpcId, nil
}

func (c EC2) getSubnetsWithVpc(vpcId string) ([]*ec2.Subnet, error) {
	response, err := c.client.DescribeSubnets(&ec2.DescribeSubnetsInput{
		Filters: [] *ec2.Filter {
			{
				Name: aws.String("vpc-id"),
				Values: aws.StringSlice([]string{vpcId}),
			},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("find subnet ids under vpc %s: %w", vpcId, err)
	}
	return response.Subnets, nil
}

// IsSubnetInVPC checks if the specified subnet is under the specified vpc
func (c EC2) IsSubnetInVPC(subnetId string, vpcId string) (bool, error) {
	subnets, err := c.getSubnetsWithVpc(vpcId)

	if err != nil {
		return false, fmt.Errorf("check if subnet %s is in vpc %s: %w", subnetId, vpcId, err)
	}

	for _, subnet := range subnets {
		if strings.Compare(*subnet.SubnetId, subnetId) == 0 {
			return true, nil
		}
	}
	return false, nil
}

func (c EC2) IsSecurityGroupInVPC(sgId string, vpcId string) (bool, error) {
	response, err := c.client.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("vpc-id"),
				Values: aws.StringSlice([]string{ vpcId }),
			},
			{
				Name: aws.String("group-id"),
				Values: aws.StringSlice([]string{ sgId }),
			},
		},
	})

	if err != nil {
		return false, fmt.Errorf("check if security group is in VPC: %w", err)
	}

	for _, sgroup := range response.SecurityGroups {
		if strings.Compare(*sgroup.GroupId, sgId) != 0 {
			return false, nil
		}
	}

	return true, nil
}