package achievementstorage

// var userExpQuery = fmt.Sprintf(
// 	`
// 	SELECT coalesce(
// 		(
// 			SELECT sum(%s) FROM %s
// 			INNER JOIN %s ON %s = %s
// 			WHERE %s IS NOT NULL AND %s = $1
// 		),
// 		0
// 	)
// 	`,
// 	sqlutils.Full(achievementsql.Exp),
// 	achievementsql.Table,

// 	// inner join
// 	userachievementsql.Table,
// 	sqlutils.Full(achievementsql.ID),
// 	sqlutils.Full(userachievementsql.AchievementId),

// 	sqlutils.Full(userachievementsql.OpenDate),
// 	sqlutils.Full(userachievementsql.UserId),
// )

// var userLevelQuery = fmt.Sprintf(
// 	`
// 		SELECT %s,%s,%s,%s FROM %s WHERE %s <= $1 AND %s >= $1
// 	`,

// 	levelsql.Level,
// 	levelsql.Rank,
// 	levelsql.MinExp,
// 	levelsql.MaxExp,

// 	levelsql.Table,

// 	// where exp in range [minExp, maxExp]
// 	levelsql.MinExp,
// 	levelsql.MaxExp,
// )

// func (s *Storage) level(tx *sql.Tx, ctx context.Context, userId int64) (usermodel.LevelInfo, error) {
// 	var err error
// 	var level usermodel.LevelInfo
// 	err = tx.QueryRowContext(ctx, userExpQuery, userId).Scan(&level.Exp)
// 	row := tx.QueryRowContext(ctx, userLevelQuery, level.Exp)
// 	err = row.Scan(
// 		&level.Level,
// 		&level.Rank,
// 		&level.MinExp,
// 		&level.MaxExp,
// 	)
// 	if err != nil {
// 		return level, Error(err, exception.NewCause("scan user level", "UserLevel", _PROVIDER))
// 	}
// 	return level, nil
// }
