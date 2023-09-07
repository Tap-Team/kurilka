// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Tap-Team/kurilka/achievements/datamanager/achievementdatamanager (interfaces: Cache,AchievementCache,AchievementStorage,AchievementDataManager)

// Package mock_achievementdatamanager is a generated GoMock package.
package achievementdatamanager

import (
        context "context"
        reflect "reflect"
        time "time"

        model "github.com/Tap-Team/kurilka/achievements/model"
        achievementmodel "github.com/Tap-Team/kurilka/internal/model/achievementmodel"
        amidtime "github.com/Tap-Team/kurilka/pkg/amidtime"
        gomock "github.com/golang/mock/gomock"
)

// MockCache is a mock of Cache interface.
type MockCache struct {
        ctrl     *gomock.Controller
        recorder *MockCacheMockRecorder
}

// MockCacheMockRecorder is the mock recorder for MockCache.
type MockCacheMockRecorder struct {
        mock *MockCache
}

// NewMockCache creates a new mock instance.
func NewMockCache(ctrl *gomock.Controller) *MockCache {
        mock := &MockCache{ctrl: ctrl}
        mock.recorder = &MockCacheMockRecorder{mock}
        return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCache) EXPECT() *MockCacheMockRecorder {
        return m.recorder
}

// RemoveUserAchievements mocks base method.
func (m *MockCache) RemoveUserAchievements(arg0 context.Context, arg1 int64) error {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "RemoveUserAchievements", arg0, arg1)
        ret0, _ := ret[0].(error)
        return ret0
}

// RemoveUserAchievements indicates an expected call of RemoveUserAchievements.
func (mr *MockCacheMockRecorder) RemoveUserAchievements(arg0, arg1 interface{}) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveUserAchievements", reflect.TypeOf((*MockCache)(nil).RemoveUserAchievements), arg0, arg1)
}

// SaveUserAchievements mocks base method.
func (m *MockCache) SaveUserAchievements(arg0 context.Context, arg1 int64, arg2 []*achievementmodel.Achievement) error {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "SaveUserAchievements", arg0, arg1, arg2)
        ret0, _ := ret[0].(error)
        return ret0
}

// SaveUserAchievements indicates an expected call of SaveUserAchievements.
func (mr *MockCacheMockRecorder) SaveUserAchievements(arg0, arg1, arg2 interface{}) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveUserAchievements", reflect.TypeOf((*MockCache)(nil).SaveUserAchievements), arg0, arg1, arg2)
}

// UserAchievements mocks base method.
func (m *MockCache) UserAchievements(arg0 context.Context, arg1 int64) ([]*achievementmodel.Achievement, error) {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "UserAchievements", arg0, arg1)
        ret0, _ := ret[0].([]*achievementmodel.Achievement)
        ret1, _ := ret[1].(error)
        return ret0, ret1
}

// UserAchievements indicates an expected call of UserAchievements.
func (mr *MockCacheMockRecorder) UserAchievements(arg0, arg1 interface{}) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserAchievements", reflect.TypeOf((*MockCache)(nil).UserAchievements), arg0, arg1)
}

// MockAchievementCache is a mock of AchievementCache interface.
type MockAchievementCache struct {
        ctrl     *gomock.Controller
        recorder *MockAchievementCacheMockRecorder
}

// MockAchievementCacheMockRecorder is the mock recorder for MockAchievementCache.
type MockAchievementCacheMockRecorder struct {
        mock *MockAchievementCache
}

// NewMockAchievementCache creates a new mock instance.
func NewMockAchievementCache(ctrl *gomock.Controller) *MockAchievementCache {
        mock := &MockAchievementCache{ctrl: ctrl}
        mock.recorder = &MockAchievementCacheMockRecorder{mock}
        return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAchievementCache) EXPECT() *MockAchievementCacheMockRecorder {
        return m.recorder
}

// MarkShown mocks base method.
func (m *MockAchievementCache) MarkShown(arg0 context.Context, arg1 int64) {
        m.ctrl.T.Helper()
        m.ctrl.Call(m, "MarkShown", arg0, arg1)
}

// MarkShown indicates an expected call of MarkShown.
func (mr *MockAchievementCacheMockRecorder) MarkShown(arg0, arg1 interface{}) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkShown", reflect.TypeOf((*MockAchievementCache)(nil).MarkShown), arg0, arg1)
}

// OpenAchievements mocks base method.
func (m *MockAchievementCache) OpenAchievements(arg0 context.Context, arg1 int64, arg2 []int64, arg3 time.Time) {
        m.ctrl.T.Helper()
        m.ctrl.Call(m, "OpenAchievements", arg0, arg1, arg2, arg3)
}

// OpenAchievements indicates an expected call of OpenAchievements.
func (mr *MockAchievementCacheMockRecorder) OpenAchievements(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenAchievements", reflect.TypeOf((*MockAchievementCache)(nil).OpenAchievements), arg0, arg1, arg2, arg3)
}

// ReachAchievements mocks base method.
func (m *MockAchievementCache) ReachAchievements(arg0 context.Context, arg1 int64, arg2 amidtime.Timestamp, arg3 []int64) {
        m.ctrl.T.Helper()
        m.ctrl.Call(m, "ReachAchievements", arg0, arg1, arg2, arg3)
}

// ReachAchievements indicates an expected call of ReachAchievements.
func (mr *MockAchievementCacheMockRecorder) ReachAchievements(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReachAchievements", reflect.TypeOf((*MockAchievementCache)(nil).ReachAchievements), arg0, arg1, arg2, arg3)
}

// SaveUserAchievements mocks base method.
func (m *MockAchievementCache) SaveUserAchievements(arg0 context.Context, arg1 int64, arg2 []*achievementmodel.Achievement) error {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "SaveUserAchievements", arg0, arg1, arg2)
        ret0, _ := ret[0].(error)
        return ret0
}

// SaveUserAchievements indicates an expected call of SaveUserAchievements.
func (mr *MockAchievementCacheMockRecorder) SaveUserAchievements(arg0, arg1, arg2 interface{}) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveUserAchievements", reflect.TypeOf((*MockAchievementCache)(nil).SaveUserAchievements), arg0, arg1, arg2)
}

// UserAchievements mocks base method.
func (m *MockAchievementCache) UserAchievements(arg0 context.Context, arg1 int64) ([]*achievementmodel.Achievement, error) {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "UserAchievements", arg0, arg1)
        ret0, _ := ret[0].([]*achievementmodel.Achievement)
        ret1, _ := ret[1].(error)
        return ret0, ret1
}

// UserAchievements indicates an expected call of UserAchievements.
func (mr *MockAchievementCacheMockRecorder) UserAchievements(arg0, arg1 interface{}) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserAchievements", reflect.TypeOf((*MockAchievementCache)(nil).UserAchievements), arg0, arg1)
}

// MockAchievementStorage is a mock of AchievementStorage interface.
type MockAchievementStorage struct {
        ctrl     *gomock.Controller
        recorder *MockAchievementStorageMockRecorder
}

// MockAchievementStorageMockRecorder is the mock recorder for MockAchievementStorage.
type MockAchievementStorageMockRecorder struct {
        mock *MockAchievementStorage
}

// NewMockAchievementStorage creates a new mock instance.
func NewMockAchievementStorage(ctrl *gomock.Controller) *MockAchievementStorage {
        mock := &MockAchievementStorage{ctrl: ctrl}
        mock.recorder = &MockAchievementStorageMockRecorder{mock}
        return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAchievementStorage) EXPECT() *MockAchievementStorageMockRecorder {
        return m.recorder
}

// InsertUserAchievements mocks base method.
func (m *MockAchievementStorage) InsertUserAchievements(arg0 context.Context, arg1 int64, arg2 amidtime.Timestamp, arg3 []int64) error {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "InsertUserAchievements", arg0, arg1, arg2, arg3)
        ret0, _ := ret[0].(error)
        return ret0
}

// InsertUserAchievements indicates an expected call of InsertUserAchievements.
func (mr *MockAchievementStorageMockRecorder) InsertUserAchievements(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertUserAchievements", reflect.TypeOf((*MockAchievementStorage)(nil).InsertUserAchievements), arg0, arg1, arg2, arg3)
}

// MarkShown mocks base method.
func (m *MockAchievementStorage) MarkShown(arg0 context.Context, arg1 int64) error {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "MarkShown", arg0, arg1)
        ret0, _ := ret[0].(error)
        return ret0
}

// MarkShown indicates an expected call of MarkShown.
func (mr *MockAchievementStorageMockRecorder) MarkShown(arg0, arg1 interface{}) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkShown", reflect.TypeOf((*MockAchievementStorage)(nil).MarkShown), arg0, arg1)
}

// OpenSingle mocks base method.
func (m *MockAchievementStorage) OpenSingle(arg0 context.Context, arg1 int64, arg2 model.OpenAchievement) error {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "OpenSingle", arg0, arg1, arg2)
        ret0, _ := ret[0].(error)
        return ret0
}

// OpenSingle indicates an expected call of OpenSingle.
func (mr *MockAchievementStorageMockRecorder) OpenSingle(arg0, arg1, arg2 interface{}) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenSingle", reflect.TypeOf((*MockAchievementStorage)(nil).OpenSingle), arg0, arg1, arg2)
}

// UserAchievements mocks base method.
func (m *MockAchievementStorage) UserAchievements(arg0 context.Context, arg1 int64) ([]*achievementmodel.Achievement, error) {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "UserAchievements", arg0, arg1)
        ret0, _ := ret[0].([]*achievementmodel.Achievement)
        ret1, _ := ret[1].(error)
        return ret0, ret1
}

// UserAchievements indicates an expected call of UserAchievements.
func (mr *MockAchievementStorageMockRecorder) UserAchievements(arg0, arg1 interface{}) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserAchievements", reflect.TypeOf((*MockAchievementStorage)(nil).UserAchievements), arg0, arg1)
}

// MockAchievementDataManager is a mock of AchievementDataManager interface.
type MockAchievementDataManager struct {
        ctrl     *gomock.Controller
        recorder *MockAchievementDataManagerMockRecorder
}

// MockAchievementDataManagerMockRecorder is the mock recorder for MockAchievementDataManager.
type MockAchievementDataManagerMockRecorder struct {
        mock *MockAchievementDataManager
}

// NewMockAchievementManager creates a new mock instance.
func NewMockAchievementManager(ctrl *gomock.Controller) *MockAchievementDataManager {
        mock := &MockAchievementDataManager{ctrl: ctrl}
        mock.recorder = &MockAchievementDataManagerMockRecorder{mock}
        return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAchievementDataManager) EXPECT() *MockAchievementDataManagerMockRecorder {
        return m.recorder
}

// MarkShown mocks base method.
func (m *MockAchievementDataManager) MarkShown(arg0 context.Context, arg1 int64) error {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "MarkShown", arg0, arg1)
        ret0, _ := ret[0].(error)
        return ret0
}

// MarkShown indicates an expected call of MarkShown.
func (mr *MockAchievementDataManagerMockRecorder) MarkShown(arg0, arg1 interface{}) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkShown", reflect.TypeOf((*MockAchievementDataManager)(nil).MarkShown), arg0, arg1)
}

// OpenSingle mocks base method.
func (m *MockAchievementDataManager) OpenSingle(arg0 context.Context, arg1, arg2 int64) (*model.OpenAchievementResponse, error) {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "OpenSingle", arg0, arg1, arg2)
        ret0, _ := ret[0].(*model.OpenAchievementResponse)
        ret1, _ := ret[1].(error)
        return ret0, ret1
}

// OpenSingle indicates an expected call of OpenSingle.
func (mr *MockAchievementDataManagerMockRecorder) OpenSingle(arg0, arg1, arg2 interface{}) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenSingle", reflect.TypeOf((*MockAchievementDataManager)(nil).OpenSingle), arg0, arg1, arg2)
}

// ReachAchievements mocks base method.
func (m *MockAchievementDataManager) ReachAchievements(arg0 context.Context, arg1 int64, arg2 amidtime.Timestamp, arg3 []int64) error {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "ReachAchievements", arg0, arg1, arg2, arg3)
        ret0, _ := ret[0].(error)
        return ret0
}

// ReachAchievements indicates an expected call of ReachAchievements.
func (mr *MockAchievementDataManagerMockRecorder) ReachAchievements(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReachAchievements", reflect.TypeOf((*MockAchievementDataManager)(nil).ReachAchievements), arg0, arg1, arg2, arg3)
}

// UserAchievements mocks base method.
func (m *MockAchievementDataManager) UserAchievements(arg0 context.Context, arg1 int64) ([]*achievementmodel.Achievement, error) {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "UserAchievements", arg0, arg1)
        ret0, _ := ret[0].([]*achievementmodel.Achievement)
        ret1, _ := ret[1].(error)
        return ret0, ret1
}

// UserAchievements indicates an expected call of UserAchievements.
func (mr *MockAchievementDataManagerMockRecorder) UserAchievements(arg0, arg1 interface{}) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserAchievements", reflect.TypeOf((*MockAchievementDataManager)(nil).UserAchievements), arg0, arg1)
}