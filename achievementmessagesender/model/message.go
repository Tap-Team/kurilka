package model

import "github.com/Tap-Team/kurilka/internal/model/achievementmodel"

type AchievementMessageData struct {
	achievementType achievementmodel.AchievementType
}

func (a AchievementMessageData) AchievementType() achievementmodel.AchievementType {
	return a.achievementType
}

func NewAchievementMessageData(achtype achievementmodel.AchievementType) AchievementMessageData {
	return AchievementMessageData{achievementType: achtype}
}

type MessageData struct {
	achievementMessageData AchievementMessageData
	deleted                bool
	userId                 int64
}

func (m *MessageData) MarkAsDeleted() {
	m.deleted = true
}

func NewMessageData(userId int64, achievementMessageData AchievementMessageData) *MessageData {
	return &MessageData{
		achievementMessageData: achievementMessageData,
		userId:                 userId,
	}
}

func (m *MessageData) IsDeleted() bool {
	return m.deleted
}
func (m *MessageData) UserId() int64 {
	return m.userId
}

func (m *MessageData) AchievementMessageData() AchievementMessageData {
	return m.achievementMessageData
}
