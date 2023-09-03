package triggerstorage_test

import (
	"database/sql"
	"fmt"

	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/triggersql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/usersql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/usertriggersql"
	"github.com/Tap-Team/kurilka/pkg/random"
	"github.com/Tap-Team/kurilka/pkg/sqlutils"
	"github.com/lib/pq"
)

type userInsert struct {
	tx         *sql.Tx
	err        error
	userId     int64
	createUser *usermodel.CreateUser
}

func (u *userInsert) InsertUserQuery() string {
	return fmt.Sprintf(
		`INSERT INTO %s (%s, %s, %s, %s, %s) VALUES ($1, $2, $3, $4, $5) RETURNING %s`,
		// insert into users
		usersql.Table,
		usersql.ID,
		usersql.Name,
		usersql.CigaretteDayAmount,
		usersql.CigarettePackAmount,
		usersql.PackPrice,

		usersql.AbstinenceTime,
	)
}

func (u *userInsert) InsertUser() {
	if u.err != nil {
		return
	}
	_, err := u.tx.Exec(
		u.InsertUserQuery(),
		u.userId,
		u.createUser.Name,
		u.createUser.CigaretteDayAmount,
		u.createUser.CigarettePackAmount,
		u.createUser.PackPrice,
	)
	if err != nil {
		u.err = fmt.Errorf("insert user, %s", err)
	}
}

var triggerIdQuery = func(triggerName usermodel.Trigger) string {
	return fmt.Sprintf(`SELECT %s FROM %s WHERE %s = '%s'`, triggersql.ID, triggersql.Table, triggersql.Name, string(triggerName))
}

func (u *userInsert) InsertTriggersQuery() string {
	return fmt.Sprintf(
		`INSERT INTO %s (%s, %s) VALUES 
			(
				$1,
				(%s)
			),
			(
				$1,
				(%s)
			),
			(
				$1,
				(%s)
			),
			(
				$1,
				(%s)
			)
		`,
		usertriggersql.Table,
		usertriggersql.UserId,
		usertriggersql.TriggerId,
		triggerIdQuery(usermodel.THANK_YOU),
		triggerIdQuery(usermodel.SUPPORT_CIGGARETTE),
		triggerIdQuery(usermodel.SUPPORT_HEALTH),
		triggerIdQuery(usermodel.SUPPORT_TRIAL),
	)
}

func (u *userInsert) InsertTriggers() {
	if u.err != nil {
		return
	}
	_, err := u.tx.Exec(u.InsertTriggersQuery(), u.userId)
	if err != nil {
		u.err = fmt.Errorf("insert triggers, %s", err)
	}
}

func insertUserWithAllTriggers(db *sql.DB, userId int64) error {
	createUser := random.StructTyped[usermodel.CreateUser]()
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx, %s", err)
	}
	defer tx.Rollback()
	insert := &userInsert{tx: tx, userId: userId, createUser: &createUser}
	insert.InsertUser()
	insert.InsertTriggers()
	if insert.err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit tx, %s", err)
	}
	return nil
}

func userTriggers(db *sql.DB, userId int64) ([]usermodel.Trigger, error) {
	query := fmt.Sprintf(
		`SELECT array_agg(%s) FROM %s INNER JOIN %s ON %s = %s WHERE %s = $1`,
		sqlutils.Full(triggersql.Name),

		usertriggersql.Table,

		triggersql.Table,
		sqlutils.Full(usertriggersql.TriggerId),
		sqlutils.Full(triggersql.ID),

		sqlutils.Full(usertriggersql.UserId),
	)
	triggers := make([]usermodel.Trigger, 0)
	err := db.QueryRow(query, userId).Scan(pq.Array(&triggers))
	if err != nil {
		return nil, fmt.Errorf("get user triggers, %s", err)
	}
	return triggers, nil
}
