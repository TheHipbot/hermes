// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/TheHipbot/hermes/pkg/storage (interfaces: Storage)

// Package mock is a generated GoMock package.
package mock

import (
	storage "github.com/TheHipbot/hermes/pkg/storage"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockStorage is a mock of Storage interface
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// AddRemote mocks base method
func (m *MockStorage) AddRemote(arg0, arg1, arg2, arg3 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddRemote", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddRemote indicates an expected call of AddRemote
func (mr *MockStorageMockRecorder) AddRemote(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddRemote", reflect.TypeOf((*MockStorage)(nil).AddRemote), arg0, arg1, arg2, arg3)
}

// AddRepository mocks base method
func (m *MockStorage) AddRepository(arg0 *storage.Repository) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddRepository", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddRepository indicates an expected call of AddRepository
func (mr *MockStorageMockRecorder) AddRepository(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddRepository", reflect.TypeOf((*MockStorage)(nil).AddRepository), arg0)
}

// Close mocks base method
func (m *MockStorage) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockStorageMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockStorage)(nil).Close))
}

// ListRemotes mocks base method
func (m *MockStorage) ListRemotes() []*storage.Remote {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListRemotes")
	ret0, _ := ret[0].([]*storage.Remote)
	return ret0
}

// ListRemotes indicates an expected call of ListRemotes
func (mr *MockStorageMockRecorder) ListRemotes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListRemotes", reflect.TypeOf((*MockStorage)(nil).ListRemotes))
}

// Open mocks base method
func (m *MockStorage) Open() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Open")
}

// Open indicates an expected call of Open
func (mr *MockStorageMockRecorder) Open() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Open", reflect.TypeOf((*MockStorage)(nil).Open))
}

// RemoveRepository mocks base method
func (m *MockStorage) RemoveRepository(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveRepository", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveRepository indicates an expected call of RemoveRepository
func (mr *MockStorageMockRecorder) RemoveRepository(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveRepository", reflect.TypeOf((*MockStorage)(nil).RemoveRepository), arg0)
}

// Save mocks base method
func (m *MockStorage) Save() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save")
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save
func (mr *MockStorageMockRecorder) Save() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockStorage)(nil).Save))
}

// SearchRemote mocks base method
func (m *MockStorage) SearchRemote(arg0 string) (*storage.Remote, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchRemote", arg0)
	ret0, _ := ret[0].(*storage.Remote)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// SearchRemote indicates an expected call of SearchRemote
func (mr *MockStorageMockRecorder) SearchRemote(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchRemote", reflect.TypeOf((*MockStorage)(nil).SearchRemote), arg0)
}

// SearchRepositories mocks base method
func (m *MockStorage) SearchRepositories(arg0 string) []storage.Repository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchRepositories", arg0)
	ret0, _ := ret[0].([]storage.Repository)
	return ret0
}

// SearchRepositories indicates an expected call of SearchRepositories
func (mr *MockStorageMockRecorder) SearchRepositories(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchRepositories", reflect.TypeOf((*MockStorage)(nil).SearchRepositories), arg0)
}
