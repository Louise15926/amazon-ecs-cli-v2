// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/pkg/task/logs.go

// Package mocks is a generated GoMock package.
package mocks

import (
	cloudwatchlogs "github.com/aws/copilot-cli/internal/pkg/aws/cloudwatchlogs"
	ecs "github.com/aws/copilot-cli/internal/pkg/aws/ecs"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockTaskDescriber is a mock of TaskDescriber interface
type MockTaskDescriber struct {
	ctrl     *gomock.Controller
	recorder *MockTaskDescriberMockRecorder
}

// MockTaskDescriberMockRecorder is the mock recorder for MockTaskDescriber
type MockTaskDescriberMockRecorder struct {
	mock *MockTaskDescriber
}

// NewMockTaskDescriber creates a new mock instance
func NewMockTaskDescriber(ctrl *gomock.Controller) *MockTaskDescriber {
	mock := &MockTaskDescriber{ctrl: ctrl}
	mock.recorder = &MockTaskDescriberMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTaskDescriber) EXPECT() *MockTaskDescriberMockRecorder {
	return m.recorder
}

// DescribeTasks mocks base method
func (m *MockTaskDescriber) DescribeTasks(cluster string, taskARNs []string) ([]*ecs.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DescribeTasks", cluster, taskARNs)
	ret0, _ := ret[0].([]*ecs.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeTasks indicates an expected call of DescribeTasks
func (mr *MockTaskDescriberMockRecorder) DescribeTasks(cluster, taskARNs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeTasks", reflect.TypeOf((*MockTaskDescriber)(nil).DescribeTasks), cluster, taskARNs)
}

// MockCWLogService is a mock of CWLogService interface
type MockCWLogService struct {
	ctrl     *gomock.Controller
	recorder *MockCWLogServiceMockRecorder
}

// MockCWLogServiceMockRecorder is the mock recorder for MockCWLogService
type MockCWLogServiceMockRecorder struct {
	mock *MockCWLogService
}

// NewMockCWLogService creates a new mock instance
func NewMockCWLogService(ctrl *gomock.Controller) *MockCWLogService {
	mock := &MockCWLogService{ctrl: ctrl}
	mock.recorder = &MockCWLogServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCWLogService) EXPECT() *MockCWLogServiceMockRecorder {
	return m.recorder
}

// TaskLogEvents mocks base method
func (m *MockCWLogService) TaskLogEvents(logGroupName string, streamLastEventTime map[string]int64, opts ...cloudwatchlogs.GetLogEventsOpts) (*cloudwatchlogs.LogEventsOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{logGroupName, streamLastEventTime}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "TaskLogEvents", varargs...)
	ret0, _ := ret[0].(*cloudwatchlogs.LogEventsOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TaskLogEvents indicates an expected call of TaskLogEvents
func (mr *MockCWLogServiceMockRecorder) TaskLogEvents(logGroupName, streamLastEventTime interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{logGroupName, streamLastEventTime}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TaskLogEvents", reflect.TypeOf((*MockCWLogService)(nil).TaskLogEvents), varargs...)
}