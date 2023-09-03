// Code generated by MockGen. DO NOT EDIT.
// Source: ./achievements/transport/achievement_handler.go

// Package transport is a generated GoMock package.
package transport

import (
	context "context"
	reflect "reflect"

	model "github.com/Tap-Team/kurilka/achievements/model"
	achievementmodel "github.com/Tap-Team/kurilka/internal/model/achievementmodel"
	gomock "github.com/golang/mock/gomock"
)

// MockAchievementUseCase is a mock of AchievementUseCase interface.
type MockAchievementUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockAchievementUseCaseMockRecorder
}

// MockAchievementUseCaseMockRecorder is the mock recorder for MockAchievementUseCase.
type MockAchievementUseCaseMockRecorder struct {
	mock *MockAchievementUseCase
}

// NewMockAchievementUseCase creates a new mock instance.
func NewMockAchievementUseCase(ctrl *gomock.Controller) *MockAchievementUseCase {
	mock := &MockAchievementUseCase{ctrl: ctrl}
	mock.recorder = &MockAchievementUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAchievementUseCase) EXPECT() *MockAchievementUseCaseMockRecorder {
	return m.recorder
}

// MarkShown mocks base method.
func (m *MockAchievementUseCase) MarkShown(ctx context.Context, userId int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MarkShown", ctx, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// MarkShown indicates an expected call of MarkShown.
func (mr *MockAchievementUseCaseMockRecorder) MarkShown(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkShown", reflect.TypeOf((*MockAchievementUseCase)(nil).MarkShown), ctx, userId)
}

// OpenAll mocks base method.
func (m *MockAchievementUseCase) OpenAll(ctx context.Context, userId int64) (*model.OpenAchievementResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenAll", ctx, userId)
	ret0, _ := ret[0].(*model.OpenAchievementResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenAll indicates an expected call of OpenAll.
func (mr *MockAchievementUseCaseMockRecorder) OpenAll(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenAll", reflect.TypeOf((*MockAchievementUseCase)(nil).OpenAll), ctx, userId)
}

// OpenSingle mocks base method.
func (m *MockAchievementUseCase) OpenSingle(ctx context.Context, userId, achievementId int64) (*model.OpenAchievementResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenSingle", ctx, userId, achievementId)
	ret0, _ := ret[0].(*model.OpenAchievementResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenSingle indicates an expected call of OpenSingle.
func (mr *MockAchievementUseCaseMockRecorder) OpenSingle(ctx, userId, achievementId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenSingle", reflect.TypeOf((*MockAchievementUseCase)(nil).OpenSingle), ctx, userId, achievementId)
}

// OpenType mocks base method.
func (m *MockAchievementUseCase) OpenType(ctx context.Context, userId int64, achtype achievementmodel.AchievementType) (*model.OpenAchievementResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenType", ctx, userId, achtype)
	ret0, _ := ret[0].(*model.OpenAchievementResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenType indicates an expected call of OpenType.
func (mr *MockAchievementUseCaseMockRecorder) OpenType(ctx, userId, achtype interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenType", reflect.TypeOf((*MockAchievementUseCase)(nil).OpenType), ctx, userId, achtype)
}

// UserAchievements mocks base method.
func (m *MockAchievementUseCase) UserAchievements(ctx context.Context, userId int64) ([]*achievementmodel.Achievement, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserAchievements", ctx, userId)
	ret0, _ := ret[0].([]*achievementmodel.Achievement)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserAchievements indicates an expected call of UserAchievements.
func (mr *MockAchievementUseCaseMockRecorder) UserAchievements(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserAchievements", reflect.TypeOf((*MockAchievementUseCase)(nil).UserAchievements), ctx, userId)
}
