package model

import (
	"time"

	"github.com/Tap-Team/kurilka/internal/model/achievementmodel"
	"github.com/Tap-Team/kurilka/pkg/amidtime"
)

type OpenAchievement struct {
	AchievementId int64              `json:"achievementId"`
	OpenTime      amidtime.Timestamp `json:"openTime"`
}

func NewOpenAchievement(
	achievementId int64,
	openTime time.Time,
) OpenAchievement {
	return OpenAchievement{
		AchievementId: achievementId,
		OpenTime:      amidtime.Timestamp{Time: openTime},
	}
}

type OpenAchievementType struct {
	AchievementType achievementmodel.AchievementType `json:"type"`
	OpenTime        amidtime.Timestamp               `json:"openTime"`
}

func NewAchievementType(
	achievementType achievementmodel.AchievementType,
	openTime time.Time,
) OpenAchievementType {
	return OpenAchievementType{
		AchievementType: achievementType,
		OpenTime:        amidtime.Timestamp{Time: openTime},
	}
}

type OpenAchievementResponse struct {
	OpenTime amidtime.Timestamp `json:"openDate"`
}

func NewOpenAchievementResponse(
	openTime time.Time,
) *OpenAchievementResponse {
	return &OpenAchievementResponse{
		OpenTime: amidtime.Timestamp{Time: openTime},
	}
}
