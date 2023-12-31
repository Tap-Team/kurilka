package achievementstorage_test

import (
	"cmp"
	"database/sql"
	"fmt"
	"math/rand"
	"sort"
	"strings"

	"slices"

	"github.com/Tap-Team/kurilka/internal/model/achievementmodel"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/achievementsql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/achievementtypesql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/userachievementsql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/usersql"
	"github.com/Tap-Team/kurilka/pkg/sqlutils"
)

func insertUser(db *sql.DB, userId int64, createUser usermodel.CreateUser) error {
	query := fmt.Sprintf(
		`INSERT INTO %s (%s, %s, %s, %s, %s) VALUES ($1,$2,$3,$4,$5)`,
		// insert into users
		usersql.Table,
		usersql.ID,
		usersql.Name,
		usersql.CigaretteDayAmount,
		usersql.CigarettePackAmount,
		usersql.PackPrice,
	)
	_, err := db.Exec(query,
		userId,
		createUser.Name,
		createUser.CigaretteDayAmount,
		createUser.CigarettePackAmount,
		createUser.PackPrice,
	)
	return err
}

type insertAchievement struct {
	level   int
	achType achievementmodel.AchievementType
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

func (a *achievementGenerator) Achievement() (insertAchievement, bool) {
	if len(a.availableTypes) == 0 {
		return insertAchievement{}, false
	}

	achIndex := rand.Intn(len(a.availableTypes))
	achtype := a.availableTypes[achIndex]

	levelIndex := rand.Intn(len(a.state[achtype]))
	level := a.state[achtype][levelIndex]

	a.state[achtype] = slices.Delete(a.state[achtype], levelIndex, levelIndex+1)

	if len(a.state[achtype]) == 0 {
		a.availableTypes = slices.Delete(a.availableTypes, achIndex, achIndex+1)
	}

	achievement := insertAchievement{level: level, achType: achtype}
	return achievement, true
}

func generateRandomAchievementList(size int) []*insertAchievement {

	achievements := make([]*insertAchievement, 0, size)

	achGen := NewAchievementGenerator()

	// Генерируем случайные достижения
	for i := 0; i < size; i++ {
		achievement, ok := achGen.Achievement()
		if !ok {
			break
		}
		achievements = append(achievements, &achievement)
	}

	return achievements
}

func userAchieve(db *sql.DB, userId int64, ach *insertAchievement) error {
	query := fmt.Sprintf(
		`
		WITH achievement_select as (
			SELECT %s as id FROM %s 
			INNER JOIN %s ON %s = %s 
			WHERE %s = $2 AND %s = $3
			GROUP BY %s
		)
		INSERT INTO %s (%s,%s) VALUES ($1, (SELECT id FROM achievement_select)) 
		`,
		sqlutils.Full(achievementsql.ID),
		achievementsql.Table,

		// inner join achievement type
		achievementtypesql.Table,
		sqlutils.Full(achievementsql.TypeId),
		sqlutils.Full(achievementtypesql.ID),

		// where level and type eq $2 and $3
		sqlutils.Full(achievementtypesql.Type),
		sqlutils.Full(achievementsql.Level),

		// group by
		sqlutils.Full(achievementsql.ID),

		userachievementsql.Table,
		userachievementsql.UserId,
		userachievementsql.AchievementId,
	)
	_, err := db.Exec(query, userId, ach.achType, ach.level)
	return err
}

func compareInsertWithAchievement(iach *insertAchievement, ach *achievementmodel.Achievement) (string, bool) {
	var fields strings.Builder
	if iach.achType != ach.Type {
		fields.WriteString("Type ")
	}
	if iach.level != ach.Level {
		fields.WriteString("Level ")
	}
	return fields.String(), fields.Len() == 0
}

func sortAchievements(achList []*achievementmodel.Achievement) {
	sort.SliceStable(achList, func(i, j int) bool {
		levelCompare := cmp.Compare(achList[i].Level, achList[j].Level)
		typeCompare := cmp.Compare(achList[i].Type, achList[j].Type)
		if levelCompare == 0 {
			return typeCompare == 1
		}
		return levelCompare == 1
	})
}

func sortInsertAchievements(iachList []*insertAchievement) {
	sort.SliceStable(iachList, func(i, j int) bool {
		levelCompare := cmp.Compare(iachList[i].level, iachList[j].level)
		typeCompare := cmp.Compare(iachList[i].achType, iachList[j].achType)
		if levelCompare == 0 {
			return typeCompare == 1
		}
		return levelCompare == 1
	})
}

func userReachedAchievements(db *sql.DB, userId int64) ([]*achievementmodel.Achievement, error) {
	achievements := make([]*achievementmodel.Achievement, 0)
	query := fmt.Sprintf(`
		SELECT %s FROM %s 
		INNER JOIN %s ON %s = %s
		INNER JOIN %s ON %s = %s
		WHERE %s = $1
		GROUP BY %s
		ORDER BY %s
	`,
		sqlutils.Full(
			achievementsql.ID,
			achievementtypesql.Type,
			achievementsql.Exp,
			achievementsql.Level,
			userachievementsql.OpenDate,
			userachievementsql.ReachDate,
			userachievementsql.Shown,
		),
		achievementsql.Table,
		// inner join
		achievementtypesql.Table,
		sqlutils.Full(achievementsql.TypeId),
		sqlutils.Full(achievementtypesql.ID),

		// left join
		userachievementsql.Table,
		sqlutils.Full(achievementsql.ID),
		sqlutils.Full(userachievementsql.AchievementId),

		// where user id eq $1
		sqlutils.Full(userachievementsql.UserId),

		// group by
		sqlutils.Full(
			achievementsql.ID,
			achievementtypesql.Type,
			achievementsql.Exp,
			achievementsql.Level,
			userachievementsql.OpenDate,
			userachievementsql.ReachDate,
			userachievementsql.Shown,
		),

		// order by
		sqlutils.Full(
			achievementsql.TypeId,
			achievementsql.Level,
		),
	)
	rows, err := db.Query(query, userId)
	if err != nil {
		return nil, fmt.Errorf("failed execute query, %s", err)
	}
	for rows.Next() {
		var achievement achievementmodel.Achievement
		err := rows.Scan(
			&achievement.ID,
			&achievement.Type,
			&achievement.Exp,
			&achievement.Level,
			&achievement.OpenDate,
			&achievement.ReachDate,
			&achievement.Shown,
		)
		if err != nil {
			return nil, fmt.Errorf("failed scan achievement, %s", err)
		}
		achievements = append(achievements, &achievement)
	}
	return achievements, nil
}
