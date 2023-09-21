package userstorage_test

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/usersql"
	"github.com/Tap-Team/kurilka/pkg/amidtime"
)

func insertUser(db *sql.DB, userId int64, createUser *usermodel.CreateUser, absTime time.Time) error {
	if createUser == nil {
		return nil
	}
	query := fmt.Sprintf(
		`INSERT INTO %s (%s, %s, %s, %s, %s, %s) VALUES ($1,$2,$3,$4,$5,$6)`,
		// insert into users
		usersql.Table,

		usersql.ID,
		usersql.Name,
		usersql.CigaretteDayAmount,
		usersql.CigarettePackAmount,
		usersql.PackPrice,
		usersql.AbstinenceTime,
	)
	_, err := db.Exec(query,
		userId,
		createUser.Name,
		createUser.CigaretteDayAmount,
		createUser.CigarettePackAmount,
		createUser.PackPrice,
		amidtime.Timestamp{Time: absTime},
	)
	return err
}
