package cloudformation

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/aws/amazon-ecs-cli-v2/internal/pkg/aws/cloudformation"

	"github.com/aws/amazon-ecs-cli-v2/internal/pkg/deploy/cloudformation/mocks"
	"github.com/golang/mock/gomock"
)

func TestCloudFormation_DeployTask(t *testing.T) {
	const (
		stackName     = "copilot-my-task"
		stackTemplate = "my-task template"
	)

	testCases := map[string]struct {
		mockCfnClient func(m *mocks.MockcfnClient)
	}{
		"create a new stack": {
			mockCfnClient: func(m *mocks.MockcfnClient) {
				stack := cloudformation.NewStack(stackName, stackTemplate)
				m.EXPECT().CreateAndWait(stack).Return(nil)
				m.EXPECT().UpdateAndWait(gomock.Any()).Times(0)
			},
		},
		"update the stack": {
			mockCfnClient: func(m *mocks.MockcfnClient) {
				stack := cloudformation.NewStack(stackName, stackTemplate)
				m.EXPECT().CreateAndWait(stack).Return(&cloudformation.ErrStackAlreadyExists{
					Name: "my-task",
				})
				m.EXPECT().UpdateAndWait(stack).Times(1).Return(nil)
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCfnClient := mocks.NewMockcfnClient(ctrl)
			if tc.mockCfnClient != nil {
				tc.mockCfnClient(mockCfnClient)
			}

			cf := CloudFormation{
				cfnClient: mockCfnClient,
			}

			conf := mocks.NewMockStackConfiguration(ctrl)
			conf.EXPECT().Template().Return(stackTemplate, nil)
			conf.EXPECT().StackName().Return(stackName)
			conf.EXPECT().Parameters().Return(nil, nil)
			conf.EXPECT().Tags().Return(nil)

			err := cf.DeployTask(conf)
			require.NoError(t, err)
		})
	}
}
