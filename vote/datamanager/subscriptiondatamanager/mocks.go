// Code generated by MockGen. DO NOT EDIT.
// Source: datamanager.go

// Package subscriptiondatamanager is a generated GoMock package.
package subscriptiondatamanager

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

// UpdateUserSubscription mocks base method.
func (m *MockSubscriptionStorage) UpdateUserSubscription(ctx context.Context, userId int64, subscription usermodel.Subscription) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserSubscription", ctx, userId, subscription)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserSubscription indicates an expected call of UpdateUserSubscription.
func (mr *MockSubscriptionStorageMockRecorder) UpdateUserSubscription(ctx, userId, subscription interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserSubscription", reflect.TypeOf((*MockSubscriptionStorage)(nil).UpdateUserSubscription), ctx, userId, subscription)
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

// MockSubscriptionCache is a mock of SubscriptionCache interface.
type MockSubscriptionCache struct {
	ctrl     *gomock.Controller
	recorder *MockSubscriptionCacheMockRecorder
}

// MockSubscriptionCacheMockRecorder is the mock recorder for MockSubscriptionCache.
type MockSubscriptionCacheMockRecorder struct {
	mock *MockSubscriptionCache
}

// NewMockSubscriptionCache creates a new mock instance.
func NewMockSubscriptionCache(ctrl *gomock.Controller) *MockSubscriptionCache {
	mock := &MockSubscriptionCache{ctrl: ctrl}
	mock.recorder = &MockSubscriptionCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSubscriptionCache) EXPECT() *MockSubscriptionCacheMockRecorder {
	return m.recorder
}

// RemoveUserSubscription mocks base method.
func (m *MockSubscriptionCache) RemoveUserSubscription(ctx context.Context, userId int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveUserSubscription", ctx, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveUserSubscription indicates an expected call of RemoveUserSubscription.
func (mr *MockSubscriptionCacheMockRecorder) RemoveUserSubscription(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveUserSubscription", reflect.TypeOf((*MockSubscriptionCache)(nil).RemoveUserSubscription), ctx, userId)
}

// UpdateUserSubscription mocks base method.
func (m *MockSubscriptionCache) UpdateUserSubscription(ctx context.Context, userId int64, subscription usermodel.Subscription) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserSubscription", ctx, userId, subscription)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserSubscription indicates an expected call of UpdateUserSubscription.
func (mr *MockSubscriptionCacheMockRecorder) UpdateUserSubscription(ctx, userId, subscription interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserSubscription", reflect.TypeOf((*MockSubscriptionCache)(nil).UpdateUserSubscription), ctx, userId, subscription)
}

// UserSubscription mocks base method.
func (m *MockSubscriptionCache) UserSubscription(ctx context.Context, userId int64) (usermodel.Subscription, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserSubscription", ctx, userId)
	ret0, _ := ret[0].(usermodel.Subscription)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserSubscription indicates an expected call of UserSubscription.
func (mr *MockSubscriptionCacheMockRecorder) UserSubscription(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserSubscription", reflect.TypeOf((*MockSubscriptionCache)(nil).UserSubscription), ctx, userId)
}

// MockSubscriptionManager is a mock of SubscriptionManager interface.
type MockSubscriptionManager struct {
	ctrl     *gomock.Controller
	recorder *MockSubscriptionManagerMockRecorder
}

// MockSubscriptionManagerMockRecorder is the mock recorder for MockSubscriptionManager.
type MockSubscriptionManagerMockRecorder struct {
	mock *MockSubscriptionManager
}

// NewMockSubscriptionManager creates a new mock instance.
func NewMockSubscriptionManager(ctrl *gomock.Controller) *MockSubscriptionManager {
	mock := &MockSubscriptionManager{ctrl: ctrl}
	mock.recorder = &MockSubscriptionManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSubscriptionManager) EXPECT() *MockSubscriptionManagerMockRecorder {
	return m.recorder
}

// SetUserSubscription mocks base method.
func (m *MockSubscriptionManager) SetUserSubscription(ctx context.Context, userId int64, subscription usermodel.Subscription) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetUserSubscription", ctx, userId, subscription)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetUserSubscription indicates an expected call of SetUserSubscription.
func (mr *MockSubscriptionManagerMockRecorder) SetUserSubscription(ctx, userId, subscription interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetUserSubscription", reflect.TypeOf((*MockSubscriptionManager)(nil).SetUserSubscription), ctx, userId, subscription)
}

// UserSubscription mocks base method.
func (m *MockSubscriptionManager) UserSubscription(ctx context.Context, userid int64) (usermodel.Subscription, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserSubscription", ctx, userid)
	ret0, _ := ret[0].(usermodel.Subscription)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserSubscription indicates an expected call of UserSubscription.
func (mr *MockSubscriptionManagerMockRecorder) UserSubscription(ctx, userid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserSubscription", reflect.TypeOf((*MockSubscriptionManager)(nil).UserSubscription), ctx, userid)
}
