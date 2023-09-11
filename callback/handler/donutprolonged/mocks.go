// Code generated by MockGen. DO NOT EDIT.
// Source: handler.go

// Package donutprolonged is a generated GoMock package.
package donutprolonged

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockProlongationer is a mock of Prolongationer interface.
type MockProlongationer struct {
	ctrl     *gomock.Controller
	recorder *MockProlongationerMockRecorder
}

// MockProlongationerMockRecorder is the mock recorder for MockProlongationer.
type MockProlongationerMockRecorder struct {
	mock *MockProlongationer
}

// NewMockProlongationer creates a new mock instance.
func NewMockProlongationer(ctrl *gomock.Controller) *MockProlongationer {
	mock := &MockProlongationer{ctrl: ctrl}
	mock.recorder = &MockProlongationerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProlongationer) EXPECT() *MockProlongationerMockRecorder {
	return m.recorder
}

// ProlongSubscription mocks base method.
func (m *MockProlongationer) ProlongSubscription(ctx context.Context, userId int64, amount int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProlongSubscription", ctx, userId, amount)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProlongSubscription indicates an expected call of ProlongSubscription.
func (mr *MockProlongationerMockRecorder) ProlongSubscription(ctx, userId, amount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProlongSubscription", reflect.TypeOf((*MockProlongationer)(nil).ProlongSubscription), ctx, userId, amount)
}