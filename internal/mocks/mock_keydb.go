// Code generated by MockGen. DO NOT EDIT.
// Source: location-service/internal/storage/wrapper (interfaces: KeyDB)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockKeyDB is a mock of KeyDB interface
type MockKeyDB struct {
	ctrl     *gomock.Controller
	recorder *MockKeyDBMockRecorder
}

// MockKeyDBMockRecorder is the mock recorder for MockKeyDB
type MockKeyDBMockRecorder struct {
	mock *MockKeyDB
}

// NewMockKeyDB creates a new mock instance
func NewMockKeyDB(ctrl *gomock.Controller) *MockKeyDB {
	mock := &MockKeyDB{ctrl: ctrl}
	mock.recorder = &MockKeyDBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockKeyDB) EXPECT() *MockKeyDBMockRecorder {
	return m.recorder
}

// Clear mocks base method
func (m *MockKeyDB) Clear() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Clear")
	ret0, _ := ret[0].(error)
	return ret0
}

// Clear indicates an expected call of Clear
func (mr *MockKeyDBMockRecorder) Clear() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Clear", reflect.TypeOf((*MockKeyDB)(nil).Clear))
}

// Delete mocks base method
func (m *MockKeyDB) Delete(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockKeyDBMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockKeyDB)(nil).Delete), arg0)
}

// Get mocks base method
func (m *MockKeyDB) Get(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockKeyDBMockRecorder) Get(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockKeyDB)(nil).Get), arg0)
}

// Set mocks base method
func (m *MockKeyDB) Set(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set
func (mr *MockKeyDBMockRecorder) Set(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockKeyDB)(nil).Set), arg0, arg1)
}

// SetIfExists mocks base method
func (m *MockKeyDB) SetIfExists(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetIfExists", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetIfExists indicates an expected call of SetIfExists
func (mr *MockKeyDBMockRecorder) SetIfExists(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetIfExists", reflect.TypeOf((*MockKeyDB)(nil).SetIfExists), arg0, arg1)
}