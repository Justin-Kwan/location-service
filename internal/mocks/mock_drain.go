// Code generated by MockGen. DO NOT EDIT.
// Source: location-service/internal/types (interfaces: Drain)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockDrain is a mock of Drain interface
type MockDrain struct {
	ctrl     *gomock.Controller
	recorder *MockDrainMockRecorder
}

// MockDrainMockRecorder is the mock recorder for MockDrain
type MockDrainMockRecorder struct {
	mock *MockDrain
}

// NewMockDrain creates a new mock instance
func NewMockDrain(ctrl *gomock.Controller) *MockDrain {
	mock := &MockDrain{ctrl: ctrl}
	mock.recorder = &MockDrainMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDrain) EXPECT() *MockDrainMockRecorder {
	return m.recorder
}

// Close mocks base method
func (m *MockDrain) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close
func (mr *MockDrainMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockDrain)(nil).Close))
}

// GetOutput mocks base method
func (m *MockDrain) GetOutput() <-chan interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOutput")
	ret0, _ := ret[0].(<-chan interface{})
	return ret0
}

// GetOutput indicates an expected call of GetOutput
func (mr *MockDrainMockRecorder) GetOutput() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOutput", reflect.TypeOf((*MockDrain)(nil).GetOutput))
}

// Read mocks base method
func (m *MockDrain) Read() (interface{}, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read")
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// Read indicates an expected call of Read
func (mr *MockDrainMockRecorder) Read() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockDrain)(nil).Read))
}

// Send mocks base method
func (m *MockDrain) Send(arg0 interface{}) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Send indicates an expected call of Send
func (mr *MockDrainMockRecorder) Send(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockDrain)(nil).Send), arg0)
}

// SetInput mocks base method
func (m *MockDrain) SetInput(arg0 <-chan interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetInput", arg0)
}

// SetInput indicates an expected call of SetInput
func (mr *MockDrainMockRecorder) SetInput(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetInput", reflect.TypeOf((*MockDrain)(nil).SetInput), arg0)
}
