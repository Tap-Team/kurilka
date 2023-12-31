// Code generated by MockGen. DO NOT EDIT.
// Source: subscriptionadmin.go

// Package subscriptionadmin is a generated GoMock package.
package subscriptionadmin

import (
	context "context"
	reflect "reflect"

	usermodel "github.com/Tap-Team/kurilka/internal/model/usermodel"
	gomock "github.com/golang/mock/gomock"
)

// MockUserSubscriptionUpdater is a mock of UserSubscriptionUpdater interface.
type MockUserSubscriptionUpdater struct {
	ctrl     *gomock.Controller
	recorder *MockUserSubscriptionUpdaterMockRecorder
}

// MockUserSubscriptionUpdaterMockRecorder is the mock recorder for MockUserSubscriptionUpdater.
type MockUserSubscriptionUpdaterMockRecorder struct {
	mock *MockUserSubscriptionUpdater
}

// NewMockUserSubscriptionUpdater creates a new mock instance.
func NewMockUserSubscriptionUpdater(ctrl *gomock.Controller) *MockUserSubscriptionUpdater {
	mock := &MockUserSubscriptionUpdater{ctrl: ctrl}
	mock.recorder = &MockUserSubscriptionUpdaterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserSubscriptionUpdater) EXPECT() *MockUserSubscriptionUpdaterMockRecorder {
	return m.recorder
}

// UpdateUserSubscription mocks base method.
func (m *MockUserSubscriptionUpdater) UpdateUserSubscription(ctx context.Context, userId int64, subscription usermodel.Subscription) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserSubscription", ctx, userId, subscription)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserSubscription indicates an expected call of UpdateUserSubscription.
func (mr *MockUserSubscriptionUpdaterMockRecorder) UpdateUserSubscription(ctx, userId, subscription interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserSubscription", reflect.TypeOf((*MockUserSubscriptionUpdater)(nil).UpdateUserSubscription), ctx, userId, subscription)
}
