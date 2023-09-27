// Code generated by MockGen. DO NOT EDIT.
// Source: friend_provider.go

// Package userusecase is a generated GoMock package.
package userusecase

import (
	context "context"
	reflect "reflect"

	usermodel "github.com/Tap-Team/kurilka/internal/model/usermodel"
	gomock "github.com/golang/mock/gomock"
)

// MockSubscriptionStorage is a mock of SubscriptionStorage interface.
type MockSubscriptionStorage struct {
	ctrl     *gomock.Controller
	recorder *MockSubscriptionStorageMockRecorder
}

// MockSubscriptionStorageMockRecorder is the mock recorder for MockSubscriptionStorage.
type MockSubscriptionStorageMockRecorder struct {
	mock *MockSubscriptionStorage
}

// NewMockSubscriptionStorage creates a new mock instance.
func NewMockSubscriptionStorage(ctrl *gomock.Controller) *MockSubscriptionStorage {
	mock := &MockSubscriptionStorage{ctrl: ctrl}
	mock.recorder = &MockSubscriptionStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSubscriptionStorage) EXPECT() *MockSubscriptionStorageMockRecorder {
	return m.recorder
}

// Clear mocks base method.
func (m *MockSubscriptionStorage) Clear(ctx context.Context, userId int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Clear", ctx, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// Clear indicates an expected call of Clear.
func (mr *MockSubscriptionStorageMockRecorder) Clear(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Clear", reflect.TypeOf((*MockSubscriptionStorage)(nil).Clear), ctx, userId)
}

// UserSubscription mocks base method.
func (m *MockSubscriptionStorage) UserSubscription(ctx context.Context, userId int64) (usermodel.Subscription, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserSubscription", ctx, userId)
	ret0, _ := ret[0].(usermodel.Subscription)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserSubscription indicates an expected call of UserSubscription.
func (mr *MockSubscriptionStorageMockRecorder) UserSubscription(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserSubscription", reflect.TypeOf((*MockSubscriptionStorage)(nil).UserSubscription), ctx, userId)
}
