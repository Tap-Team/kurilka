// Code generated by MockGen. DO NOT EDIT.
// Source: privacysetting.go

// Package privacysetting is a generated GoMock package.
package privacysetting

import (
	context "context"
	reflect "reflect"

	usermodel "github.com/Tap-Team/kurilka/internal/model/usermodel"
	gomock "github.com/golang/mock/gomock"
)

// MockPrivacySettingSwitcher is a mock of PrivacySettingSwitcher interface.
type MockPrivacySettingSwitcher struct {
	ctrl     *gomock.Controller
	recorder *MockPrivacySettingSwitcherMockRecorder
}

// MockPrivacySettingSwitcherMockRecorder is the mock recorder for MockPrivacySettingSwitcher.
type MockPrivacySettingSwitcherMockRecorder struct {
	mock *MockPrivacySettingSwitcher
}

// NewMockPrivacySettingSwitcher creates a new mock instance.
func NewMockPrivacySettingSwitcher(ctrl *gomock.Controller) *MockPrivacySettingSwitcher {
	mock := &MockPrivacySettingSwitcher{ctrl: ctrl}
	mock.recorder = &MockPrivacySettingSwitcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPrivacySettingSwitcher) EXPECT() *MockPrivacySettingSwitcherMockRecorder {
	return m.recorder
}

// Switch mocks base method.
func (m *MockPrivacySettingSwitcher) Switch(ctx context.Context, userId int64, setting usermodel.PrivacySetting) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Switch", ctx, userId, setting)
	ret0, _ := ret[0].(error)
	return ret0
}

// Switch indicates an expected call of Switch.
func (mr *MockPrivacySettingSwitcherMockRecorder) Switch(ctx, userId, setting interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Switch", reflect.TypeOf((*MockPrivacySettingSwitcher)(nil).Switch), ctx, userId, setting)
}