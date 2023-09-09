// Code generated by MockGen. DO NOT EDIT.
// Source: subscription.go

// Package subscriptionusecase is a generated GoMock package.
package subscriptionusecase

import (
	context "context"
	reflect "reflect"
	time "time"

	usermodel "github.com/Tap-Team/kurilka/internal/model/usermodel"
	gomock "github.com/golang/mock/gomock"
)

// MockVK_Subscription_Manager is a mock of VK_Subscription_Manager interface.
type MockVK_Subscription_Manager struct {
	ctrl     *gomock.Controller
	recorder *MockVK_Subscription_ManagerMockRecorder
}

// MockVK_Subscription_ManagerMockRecorder is the mock recorder for MockVK_Subscription_Manager.
type MockVK_Subscription_ManagerMockRecorder struct {
	mock *MockVK_Subscription_Manager
}

// NewMockVK_Subscription_Manager creates a new mock instance.
func NewMockVK_Subscription_Manager(ctrl *gomock.Controller) *MockVK_Subscription_Manager {
	mock := &MockVK_Subscription_Manager{ctrl: ctrl}
	mock.recorder = &MockVK_Subscription_ManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVK_Subscription_Manager) EXPECT() *MockVK_Subscription_ManagerMockRecorder {
	return m.recorder
}

// UserSubscriptionById mocks base method.
func (m *MockVK_Subscription_Manager) UserSubscriptionById(ctx context.Context, userId int64) (time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserSubscriptionById", ctx, userId)
	ret0, _ := ret[0].(time.Time)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserSubscriptionById indicates an expected call of UserSubscriptionById.
func (mr *MockVK_Subscription_ManagerMockRecorder) UserSubscriptionById(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserSubscriptionById", reflect.TypeOf((*MockVK_Subscription_Manager)(nil).UserSubscriptionById), ctx, userId)
}

// MockSubscriptionUseCase is a mock of SubscriptionUseCase interface.
type MockSubscriptionUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockSubscriptionUseCaseMockRecorder
}

// MockSubscriptionUseCaseMockRecorder is the mock recorder for MockSubscriptionUseCase.
type MockSubscriptionUseCaseMockRecorder struct {
	mock *MockSubscriptionUseCase
}

// NewMockSubscriptionUseCase creates a new mock instance.
func NewMockSubscriptionUseCase(ctrl *gomock.Controller) *MockSubscriptionUseCase {
	mock := &MockSubscriptionUseCase{ctrl: ctrl}
	mock.recorder = &MockSubscriptionUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSubscriptionUseCase) EXPECT() *MockSubscriptionUseCaseMockRecorder {
	return m.recorder
}

// UserSubscription mocks base method.
func (m *MockSubscriptionUseCase) UserSubscription(ctx context.Context, userId int64) usermodel.SubscriptionType {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserSubscription", ctx, userId)
	ret0, _ := ret[0].(usermodel.SubscriptionType)
	return ret0
}

// UserSubscription indicates an expected call of UserSubscription.
func (mr *MockSubscriptionUseCaseMockRecorder) UserSubscription(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserSubscription", reflect.TypeOf((*MockSubscriptionUseCase)(nil).UserSubscription), ctx, userId)
}
