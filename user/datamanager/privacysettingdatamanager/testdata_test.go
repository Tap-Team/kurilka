package privacysettingdatamanager_test

import (
	"math/rand"
	"slices"

	"github.com/Tap-Team/kurilka/internal/model/usermodel"
)

func randomPrivacySettingList(size int) []usermodel.PrivacySetting {
	privacySettings := []usermodel.PrivacySetting{
		usermodel.STATISTICS_MONEY,
		usermodel.STATISTICS_CIGARETTE,
		usermodel.STATISTICS_LIFE,
		usermodel.STATISTICS_TIME,
		usermodel.ACHIEVEMENTS_DURATION,
		usermodel.ACHIEVEMENTS_HEALTH,
		usermodel.ACHIEVEMENTS_WELL_BEING,
		usermodel.ACHIEVEMENTS_SAVING,
		usermodel.ACHIEVEMENTS_CIGARETTE,
	}

	if len(privacySettings) >= size {
		return privacySettings
	}
	if len(privacySettings) < 0 {
		return []usermodel.PrivacySetting{}
	}
	for len(privacySettings) != size {
		index := rand.Intn(len(privacySettings))
		privacySettings = slices.Delete(privacySettings, index, index+1)
	}
	return privacySettings
}
