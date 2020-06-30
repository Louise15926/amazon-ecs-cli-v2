// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/pkg/describe/status.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	cloudwatch "github.com/aws/copilot-cli/internal/pkg/aws/cloudwatch"
	ecs "github.com/aws/copilot-cli/internal/pkg/aws/ecs"
	gomock "github.com/golang/mock/gomock"
)

// MockalarmStatusGetter is a mock of alarmStatusGetter interface
type MockalarmStatusGetter struct {
	ctrl     *gomock.Controller
	recorder *MockalarmStatusGetterMockRecorder
}

// MockalarmStatusGetterMockRecorder is the mock recorder for MockalarmStatusGetter
type MockalarmStatusGetterMockRecorder struct {
	mock *MockalarmStatusGetter
}

// NewMockalarmStatusGetter creates a new mock instance
func NewMockalarmStatusGetter(ctrl *gomock.Controller) *MockalarmStatusGetter {
	mock := &MockalarmStatusGetter{ctrl: ctrl}
	mock.recorder = &MockalarmStatusGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockalarmStatusGetter) EXPECT() *MockalarmStatusGetterMockRecorder {
	return m.recorder
}

// GetAlarmsWithTags mocks base method
func (m *MockalarmStatusGetter) GetAlarmsWithTags(tags map[string]string) ([]cloudwatch.AlarmStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAlarmsWithTags", tags)
	ret0, _ := ret[0].([]cloudwatch.AlarmStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAlarmsWithTags indicates an expected call of GetAlarmsWithTags
func (mr *MockalarmStatusGetterMockRecorder) GetAlarmsWithTags(tags interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAlarmsWithTags", reflect.TypeOf((*MockalarmStatusGetter)(nil).GetAlarmsWithTags), tags)
}

// MockecsServiceGetter is a mock of ecsServiceGetter interface
type MockecsServiceGetter struct {
	ctrl     *gomock.Controller
	recorder *MockecsServiceGetterMockRecorder
}

// MockecsServiceGetterMockRecorder is the mock recorder for MockecsServiceGetter
type MockecsServiceGetterMockRecorder struct {
	mock *MockecsServiceGetter
}

// NewMockecsServiceGetter creates a new mock instance
func NewMockecsServiceGetter(ctrl *gomock.Controller) *MockecsServiceGetter {
	mock := &MockecsServiceGetter{ctrl: ctrl}
	mock.recorder = &MockecsServiceGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockecsServiceGetter) EXPECT() *MockecsServiceGetterMockRecorder {
	return m.recorder
}

// ServiceTasks mocks base method
func (m *MockecsServiceGetter) ServiceTasks(clusterName, serviceName string) ([]*ecs.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ServiceTasks", clusterName, serviceName)
	ret0, _ := ret[0].([]*ecs.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ServiceTasks indicates an expected call of ServiceTasks
func (mr *MockecsServiceGetterMockRecorder) ServiceTasks(clusterName, serviceName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ServiceTasks", reflect.TypeOf((*MockecsServiceGetter)(nil).ServiceTasks), clusterName, serviceName)
}

// Service mocks base method
func (m *MockecsServiceGetter) Service(clusterName, serviceName string) (*ecs.Service, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Service", clusterName, serviceName)
	ret0, _ := ret[0].(*ecs.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Service indicates an expected call of Service
func (mr *MockecsServiceGetterMockRecorder) Service(clusterName, serviceName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Service", reflect.TypeOf((*MockecsServiceGetter)(nil).Service), clusterName, serviceName)
}

// MockserviceArnGetter is a mock of serviceArnGetter interface
type MockserviceArnGetter struct {
	ctrl     *gomock.Controller
	recorder *MockserviceArnGetterMockRecorder
}

// MockserviceArnGetterMockRecorder is the mock recorder for MockserviceArnGetter
type MockserviceArnGetterMockRecorder struct {
	mock *MockserviceArnGetter
}

// NewMockserviceArnGetter creates a new mock instance
func NewMockserviceArnGetter(ctrl *gomock.Controller) *MockserviceArnGetter {
	mock := &MockserviceArnGetter{ctrl: ctrl}
	mock.recorder = &MockserviceArnGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockserviceArnGetter) EXPECT() *MockserviceArnGetterMockRecorder {
	return m.recorder
}

// GetServiceArn mocks base method
func (m *MockserviceArnGetter) GetServiceArn() (*ecs.ServiceArn, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetServiceArn")
	ret0, _ := ret[0].(*ecs.ServiceArn)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetServiceArn indicates an expected call of GetServiceArn
func (mr *MockserviceArnGetterMockRecorder) GetServiceArn() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetServiceArn", reflect.TypeOf((*MockserviceArnGetter)(nil).GetServiceArn))
}
