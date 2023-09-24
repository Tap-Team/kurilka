// Code generated by MockGen. DO NOT EDIT.
// Source: scheduler.go

// Package achievementmessagesender is a generated GoMock package.
package achievementmessagesender

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockAchievementMessageSenderAtTime is a mock of AchievementMessageSenderAtTime interface.
type MockAchievementMessageSenderAtTime struct {
	ctrl     *gomock.Controller
	recorder *MockAchievementMessageSenderAtTimeMockRecorder
}

// MockAchievementMessageSenderAtTimeMockRecorder is the mock recorder for MockAchievementMessageSenderAtTime.
type MockAchievementMessageSenderAtTimeMockRecorder struct {
	mock *MockAchievementMessageSenderAtTime
}

// NewMockAchievementMessageSenderAtTime creates a new mock instance.
func NewMockAchievementMessageSenderAtTime(ctrl *gomock.Controller) *MockAchievementMessageSenderAtTime {
	mock := &MockAchievementMessageSenderAtTime{ctrl: ctrl}
	mock.recorder = &MockAchievementMessageSenderAtTimeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAchievementMessageSenderAtTime) EXPECT() *MockAchievementMessageSenderAtTimeMockRecorder {
	return m.recorder
}

// SendMessageAtTime mocks base method.
func (m *MockAchievementMessageSenderAtTime) SendMessageAtTime(ctx context.Context, userId int64, messageData AchievementMessageData, t time.Time) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SendMessageAtTime", ctx, userId, messageData, t)
}

// SendMessageAtTime indicates an expected call of SendMessageAtTime.
func (mr *MockAchievementMessageSenderAtTimeMockRecorder) SendMessageAtTime(ctx, userId, messageData, t interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessageAtTime", reflect.TypeOf((*MockAchievementMessageSenderAtTime)(nil).SendMessageAtTime), ctx, userId, messageData, t)
}