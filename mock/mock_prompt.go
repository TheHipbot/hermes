// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/TheHipbot/hermes/pkg/prompt (interfaces: InputPrompt,SelectPrompt,Factory)

// Package mock is a generated GoMock package.
package mock

import (
	prompt "github.com/TheHipbot/hermes/pkg/prompt"
	gomock "github.com/golang/mock/gomock"
	promptui "github.com/manifoldco/promptui"
	reflect "reflect"
)

// MockInputPrompt is a mock of InputPrompt interface
type MockInputPrompt struct {
	ctrl     *gomock.Controller
	recorder *MockInputPromptMockRecorder
}

// MockInputPromptMockRecorder is the mock recorder for MockInputPrompt
type MockInputPromptMockRecorder struct {
	mock *MockInputPrompt
}

// NewMockInputPrompt creates a new mock instance
func NewMockInputPrompt(ctrl *gomock.Controller) *MockInputPrompt {
	mock := &MockInputPrompt{ctrl: ctrl}
	mock.recorder = &MockInputPromptMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockInputPrompt) EXPECT() *MockInputPromptMockRecorder {
	return m.recorder
}

// Run mocks base method
func (m *MockInputPrompt) Run() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Run")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Run indicates an expected call of Run
func (mr *MockInputPromptMockRecorder) Run() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockInputPrompt)(nil).Run))
}

// MockSelectPrompt is a mock of SelectPrompt interface
type MockSelectPrompt struct {
	ctrl     *gomock.Controller
	recorder *MockSelectPromptMockRecorder
}

// MockSelectPromptMockRecorder is the mock recorder for MockSelectPrompt
type MockSelectPromptMockRecorder struct {
	mock *MockSelectPrompt
}

// NewMockSelectPrompt creates a new mock instance
func NewMockSelectPrompt(ctrl *gomock.Controller) *MockSelectPrompt {
	mock := &MockSelectPrompt{ctrl: ctrl}
	mock.recorder = &MockSelectPromptMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSelectPrompt) EXPECT() *MockSelectPromptMockRecorder {
	return m.recorder
}

// Run mocks base method
func (m *MockSelectPrompt) Run() (int, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Run")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Run indicates an expected call of Run
func (mr *MockSelectPromptMockRecorder) Run() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockSelectPrompt)(nil).Run))
}

// MockFactory is a mock of Factory interface
type MockFactory struct {
	ctrl     *gomock.Controller
	recorder *MockFactoryMockRecorder
}

// MockFactoryMockRecorder is the mock recorder for MockFactory
type MockFactoryMockRecorder struct {
	mock *MockFactory
}

// NewMockFactory creates a new mock instance
func NewMockFactory(ctrl *gomock.Controller) *MockFactory {
	mock := &MockFactory{ctrl: ctrl}
	mock.recorder = &MockFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFactory) EXPECT() *MockFactoryMockRecorder {
	return m.recorder
}

// CreateInputPrompt mocks base method
func (m *MockFactory) CreateInputPrompt(arg0 string) prompt.InputPrompt {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateInputPrompt", arg0)
	ret0, _ := ret[0].(prompt.InputPrompt)
	return ret0
}

// CreateInputPrompt indicates an expected call of CreateInputPrompt
func (mr *MockFactoryMockRecorder) CreateInputPrompt(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateInputPrompt", reflect.TypeOf((*MockFactory)(nil).CreateInputPrompt), arg0)
}

// CreateSelectPrompt mocks base method
func (m *MockFactory) CreateSelectPrompt(arg0 string, arg1 interface{}, arg2 *promptui.SelectTemplates) prompt.SelectPrompt {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSelectPrompt", arg0, arg1, arg2)
	ret0, _ := ret[0].(prompt.SelectPrompt)
	return ret0
}

// CreateSelectPrompt indicates an expected call of CreateSelectPrompt
func (mr *MockFactoryMockRecorder) CreateSelectPrompt(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSelectPrompt", reflect.TypeOf((*MockFactory)(nil).CreateSelectPrompt), arg0, arg1, arg2)
}
