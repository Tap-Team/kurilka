package achievement_test

import (
	"math/rand"
	"time"

	"slices"

	"github.com/Tap-Team/kurilka/internal/model/achievementmodel"
	"github.com/Tap-Team/kurilka/pkg/amidtime"
)

func compareAchievements(a1, a2 *achievementmodel.Achievement) bool {
	if a1.ID != a2.ID {
		return false
	}
	if a1.Type != a2.Type {
		return false
	}
	if a1.Exp != a2.Exp {
		return false
	}
	if a1.Level != a2.Level {
		return false
	}
	if a1.Shown != a2.Shown {
		return false
	}
	return true
}

type achievementGenerator struct {
	availableTypes []achievementmodel.AchievementType
	state          map[achievementmodel.AchievementType][]int
}

func NewAchievementGenerator() *achievementGenerator {
	achtypeList := []achievementmodel.AchievementType{
		achievementmodel.DURATION,
		achievementmodel.CIGARETTE,
		achievementmodel.HEALTH,
		achievementmodel.WELL_BEING,
		achievementmodel.SAVING,
	}
	state := make(map[achievementmodel.AchievementType][]int, len(achtypeList))
	for _, tp := range achtypeList {
		state[tp] = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	}
	return &achievementGenerator{
		availableTypes: achtypeList,
		state:          state,
	}
}

func (a *achievementGenerator) Achievement(
	id int64,
	openDate time.Time,
	reachDate time.Time,
	shown bool,
) *achievementmodel.Achievement {
	if len(a.availableTypes) == 0 {
		return nil
	}

	achIndex := rand.Intn(len(a.availableTypes))
	achtype := a.availableTypes[achIndex]

	levelIndex := rand.Intn(len(a.state[achtype]))
	level := a.state[achtype][levelIndex]

	a.state[achtype] = slices.Delete(a.state[achtype], levelIndex, levelIndex+1)

	if len(a.state[achtype]) == 0 {
		a.availableTypes = slices.Delete(a.availableTypes, achIndex, achIndex+1)
	}
	achievement := achievementmodel.NewAchievement(
		id,
		achtype,
		level,
		20,
		amidtime.Timestamp{Time: openDate},
		amidtime.Timestamp{Time: reachDate},
		shown,
		0,
	)

	return achievement
}

func generateRandomAchievementList(size int, opts ...func(*achievementmodel.Achievement)) []*achievementmodel.Achievement {

	achievements := make([]*achievementmodel.Achievement, 0, size)

	achGen := NewAchievementGenerator()
	size = min(size, 50)
	// Генерируем случайные достижения
	for i := 0; i < size; i++ {
		achievement := achGen.Achievement(int64(i+1), time.Time{}, time.Time{}, false)
		for _, opt := range opts {
			opt(achievement)
		}
		achievements = append(achievements, achievement)
	}

	return achievements
}
