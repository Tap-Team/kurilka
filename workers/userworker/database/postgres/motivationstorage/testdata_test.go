package motivationstorage_test

import (
	"database/sql"
	"fmt"

	"github.com/Tap-Team/kurilka/internal/sqlmodel/usersql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/welcomemotivationsql"
	"github.com/Tap-Team/kurilka/pkg/sqlutils"
	"github.com/Tap-Team/kurilka/workers/userworker/model"
)

func userMotivation(db *sql.DB, userId int64) (model.Motivation, error) {
	query := fmt.Sprintf(
		`SELECT %s, %s FROM %s INNER JOIN %s ON %s = %s WHERE %s = $1`,
		sqlutils.Full(welcomemotivationsql.ID),
		sqlutils.Full(welcomemotivationsql.Motivation),

		welcomemotivationsql.Table,

		usersql.Table,

		sqlutils.Full(welcomemotivationsql.ID),
		sqlutils.Full(usersql.MotivationId),

		sqlutils.Full(usersql.ID),
	)

	var motivation model.Motivation
	err := db.QueryRow(query, userId).Scan(&motivation.ID, &motivation.Motivation)
	if err != nil {
		return motivation, err
	}
	return motivation, nil
}

func insertUser(db *sql.DB, userId int64, motivationId int) error {
	query := fmt.Sprintf(
		`INSERT INTO %s (%s,%s,%s,%s,%s,%s) VALUES($1,$2,$3,$4,$5,$6)`,
		usersql.Table,

		usersql.ID,
		usersql.Name,
		usersql.CigaretteDayAmount,
		usersql.CigarettePackAmount,
		usersql.PackPrice,
		usersql.MotivationId,
	)
	_, err := db.Exec(query, userId, "dima", 1, 1, 1.0, motivationId)
	return err
}
