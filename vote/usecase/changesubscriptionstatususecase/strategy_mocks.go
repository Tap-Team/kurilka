// Code generated by MockGen. DO NOT EDIT.
// Source: strategy.go

// Package changesubscriptionstatususecase is a generated GoMock package.
package changesubscriptionstatususecase

import (
	context "context"
	reflect "reflect"

	subscription "github.com/Tap-Team/kurilka/vote/model/subscription"
	gomock "github.com/golang/mock/gomock"
)

// MockChangeSubscriptionStatusStrategy is a mock of ChangeSubscriptionStatusStrategy interface.
type MockChangeSubscriptionStatusStrategy struct {
	ctrl     *gomock.Controller
	recorder *MockChangeSubscriptionStatusStrategyMockRecorder
}

// MockChangeSubscriptionStatusStrategyMockRecorder is the mock recorder for MockChangeSubscriptionStatusStrategy.
type MockChangeSubscriptionStatusStrategyMockRecorder struct {
	mock *MockChangeSubscriptionStatusStrategy
}

// NewMockChangeSubscriptionStatusStrategy creates a new mock instance.
func NewMockChangeSubscriptionStatusStrategy(ctrl *gomock.Controller) *MockChangeSubscriptionStatusStrategy {
	mock := &MockChangeSubscriptionStatusStrategy{ctrl: ctrl}
	mock.recorder = &MockChangeSubscriptionStatusStrategyMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChangeSubscriptionStatusStrategy) EXPECT() *MockChangeSubscriptionStatusStrategyMockRecorder {
	return m.recorder
}

// Change mocks base method.
func (m *MockChangeSubscriptionStatusStrategy) Change(ctx context.Context, changeSubscriptionStatus subscription.ChangeSubscriptionStatus) (subscription.ChangeSubscriptionStatusResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Change", ctx, changeSubscriptionStatus)
	ret0, _ := ret[0].(subscription.ChangeSubscriptionStatusResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Change indicates an expected call of Change.
func (mr *MockChangeSubscriptionStatusStrategyMockRecorder) Change(ctx, changeSubscriptionStatus interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Change", reflect.TypeOf((*MockChangeSubscriptionStatusStrategy)(nil).Change), ctx, changeSubscriptionStatus)
}

// MockVoteSubscriptionStorage is a mock of VoteSubscriptionStorage interface.
type MockVoteSubscriptionStorage struct {
	ctrl     *gomock.Controller
	recorder *MockVoteSubscriptionStorageMockRecorder
}

// MockVoteSubscriptionStorageMockRecorder is the mock recorder for MockVoteSubscriptionStorage.
type MockVoteSubscriptionStorageMockRecorder struct {
	mock *MockVoteSubscriptionStorage
}

// NewMockVoteSubscriptionStorage creates a new mock instance.
func NewMockVoteSubscriptionStorage(ctrl *gomock.Controller) *MockVoteSubscriptionStorage {
	mock := &MockVoteSubscriptionStorage{ctrl: ctrl}
	mock.recorder = &MockVoteSubscriptionStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVoteSubscriptionStorage) EXPECT() *MockVoteSubscriptionStorageMockRecorder {
	return m.recorder
}

// CreateSubscription mocks base method.
func (m *MockVoteSubscriptionStorage) CreateSubscription(ctx context.Context, subscriptionId, userId int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSubscription", ctx, subscriptionId, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateSubscription indicates an expected call of CreateSubscription.
func (mr *MockVoteSubscriptionStorageMockRecorder) CreateSubscription(ctx, subscriptionId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSubscription", reflect.TypeOf((*MockVoteSubscriptionStorage)(nil).CreateSubscription), ctx, subscriptionId, userId)
}

// DeleteSubscription mocks base method.
func (m *MockVoteSubscriptionStorage) DeleteSubscription(ctx context.Context, subscriptionId int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSubscription", ctx, subscriptionId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSubscription indicates an expected call of DeleteSubscription.
func (mr *MockVoteSubscriptionStorageMockRecorder) DeleteSubscription(ctx, subscriptionId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSubscription", reflect.TypeOf((*MockVoteSubscriptionStorage)(nil).DeleteSubscription), ctx, subscriptionId)
}

// UpdateUserSubscriptionId mocks base method.
func (m *MockVoteSubscriptionStorage) UpdateUserSubscriptionId(ctx context.Context, userId, subscriptionId int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserSubscriptionId", ctx, userId, subscriptionId)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserSubscriptionId indicates an expected call of UpdateUserSubscriptionId.
func (mr *MockVoteSubscriptionStorageMockRecorder) UpdateUserSubscriptionId(ctx, userId, subscriptionId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserSubscriptionId", reflect.TypeOf((*MockVoteSubscriptionStorage)(nil).UpdateUserSubscriptionId), ctx, userId, subscriptionId)
}

// MockSubscriptionItemStroage is a mock of SubscriptionItemStroage interface.
type MockSubscriptionItemStroage struct {
	ctrl     *gomock.Controller
	recorder *MockSubscriptionItemStroageMockRecorder
}

// MockSubscriptionItemStroageMockRecorder is the mock recorder for MockSubscriptionItemStroage.
type MockSubscriptionItemStroageMockRecorder struct {
	mock *MockSubscriptionItemStroage
}

// NewMockSubscriptionItemStroage creates a new mock instance.
func NewMockSubscriptionItemStroage(ctrl *gomock.Controller) *MockSubscriptionItemStroage {
	mock := &MockSubscriptionItemStroage{ctrl: ctrl}
	mock.recorder = &MockSubscriptionItemStroageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSubscriptionItemStroage) EXPECT() *MockSubscriptionItemStroageMockRecorder {
	return m.recorder
}

// Subscription mocks base method.
func (m *MockSubscriptionItemStroage) Subscription(ctx context.Context, subscriptionId string) (subscription.Subscription, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Subscription", ctx, subscriptionId)
	ret0, _ := ret[0].(subscription.Subscription)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Subscription indicates an expected call of Subscription.
func (mr *MockSubscriptionItemStroageMockRecorder) Subscription(ctx, subscriptionId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscription", reflect.TypeOf((*MockSubscriptionItemStroage)(nil).Subscription), ctx, subscriptionId)
}
