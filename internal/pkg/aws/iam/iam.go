package iam

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

type api interface {
	GetRole(*iam.GetRoleInput) (*iam.GetRoleOutput, error)
}

type IAM struct {
	client api
}

func New(s *session.Session) IAM {
	return IAM{
		client: iam.New(s),
	}
}

func (c IAM) GetRole(roleName string) {

}