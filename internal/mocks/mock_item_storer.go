// Code generated by MockGen. DO NOT EDIT.
// Source: location-service/internal/storage/wrapper (interfaces: ItemStorer)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	internal "location-service/internal"
	reflect "reflect"
)

// MockItemStorer is a mock of ItemStorer interface
type MockItemStorer struct {
	ctrl     *gomock.Controller
	recorder *MockItemStorerMockRecorder
}

// MockItemStorerMockRecorder is the mock recorder for MockItemStorer
type MockItemStorerMockRecorder struct {
	mock *MockItemStorer
}

// NewMockItemStorer creates a new mock instance
func NewMockItemStorer(ctrl *gomock.Controller) *MockItemStorer {
	mock := &MockItemStorer{ctrl: ctrl}
	mock.recorder = &MockItemStorerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockItemStorer) EXPECT() *MockItemStorerMockRecorder {
	return m.recorder
}

// addNewItem mocks base method
func (m *MockItemStorer) addNewItem(arg0 *internal.TrackedItem) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "addNewItem", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// addNewItem indicates an expected call of addNewItem
func (mr *MockItemStorerMockRecorder) addNewItem(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "addNewItem", reflect.TypeOf((*MockItemStorer)(nil).addNewItem), arg0)
}

// delete mocks base method
func (m *MockItemStorer) delete(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "delete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// delete indicates an expected call of delete
func (mr *MockItemStorerMockRecorder) delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "delete", reflect.TypeOf((*MockItemStorer)(nil).delete), arg0)
}

// findAllNearbyItemIDs mocks base method
func (m *MockItemStorer) findAllNearbyItemIDs(arg0 *internal.Location, arg1 float64) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "findAllNearbyItemIDs", arg0, arg1)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// findAllNearbyItemIDs indicates an expected call of findAllNearbyItemIDs
func (mr *MockItemStorerMockRecorder) findAllNearbyItemIDs(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "findAllNearbyItemIDs", reflect.TypeOf((*MockItemStorer)(nil).findAllNearbyItemIDs), arg0, arg1)
}

// getItem mocks base method
func (m *MockItemStorer) getItem(arg0 string) (*internal.TrackedItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getItem", arg0)
	ret0, _ := ret[0].(*internal.TrackedItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// getItem indicates an expected call of getItem
func (mr *MockItemStorerMockRecorder) getItem(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getItem", reflect.TypeOf((*MockItemStorer)(nil).getItem), arg0)
}

// getUnmatchedNearby mocks base method
func (m *MockItemStorer) getUnmatchedNearby(arg0 map[string]float64, arg1 float64) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getUnmatchedNearby", arg0, arg1)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// getUnmatchedNearby indicates an expected call of getUnmatchedNearby
func (mr *MockItemStorerMockRecorder) getUnmatchedNearby(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getUnmatchedNearby", reflect.TypeOf((*MockItemStorer)(nil).getUnmatchedNearby), arg0, arg1)
}

// setMatched mocks base method
func (m *MockItemStorer) setMatched(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "setMatched", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// setMatched indicates an expected call of setMatched
func (mr *MockItemStorerMockRecorder) setMatched(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "setMatched", reflect.TypeOf((*MockItemStorer)(nil).setMatched), arg0)
}

// setUnmatched mocks base method
func (m *MockItemStorer) setUnmatched(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "setUnmatched", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// setUnmatched indicates an expected call of setUnmatched
func (mr *MockItemStorerMockRecorder) setUnmatched(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "setUnmatched", reflect.TypeOf((*MockItemStorer)(nil).setUnmatched), arg0)
}

// update mocks base method
func (m *MockItemStorer) update(arg0 *internal.TrackedItem) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "update", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// update indicates an expected call of update
func (mr *MockItemStorerMockRecorder) update(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "update", reflect.TypeOf((*MockItemStorer)(nil).update), arg0)
}
