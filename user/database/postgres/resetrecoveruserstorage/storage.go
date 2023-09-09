package resetrecoveruserstorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/levelsql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/motivationsql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/triggersql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/userachievementsql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/userprivacysettingsql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/usersql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/usertriggersql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/welcomemotivationsql"
	"github.com/Tap-Team/kurilka/pkg/exception"
	"github.com/Tap-Team/kurilka/pkg/sqlutils"
	"github.com/lib/pq"
)

const _PROVIDER = "user/database/postgres/resetrecoveruserstorage"

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func Error(err error, cause exception.Cause) error {
	return exception.Wrap(err, cause)
}

type clearUser struct {
	userId int64
	tx     *sql.Tx
	err    error
}

func (d *clearUser) MarkDeleted(ctx context.Context) int64 {
	if d.err != nil {
		return 0
	}
	query := fmt.Sprintf(`UPDATE %s SET %s = TRUE WHERE %s = $1`,
		usersql.Table,
		usersql.Deleted,
		usersql.ID,
	)
	res, err := d.tx.ExecContext(ctx, query, d.userId)
	if err != nil {
		d.err = Error(err, exception.NewCause("mark user deleted", "MarkDeleted", _PROVIDER))
		return 0
	}
	r, err := res.RowsAffected()
	if err != nil {
		d.err = Error(err, exception.NewCause("get rows affected", "MarkDeleted", _PROVIDER))
	}
	return r
}

func (d *clearUser) Achievements(ctx context.Context) {
	if d.err != nil {
		return
	}
	query := fmt.Sprintf(`DELETE FROM %s WHERE %s = $1`, userachievementsql.Table, userachievementsql.UserId)
	_, err := d.tx.Exec(query, d.userId)
	if err != nil {
		d.err = Error(err, exception.NewCause("clear achievements", "Achievements", _PROVIDER))
	}
}

func (d *clearUser) PrivacySettings(ctx context.Context) {
	if d.err != nil {
		return
	}
	query := fmt.Sprintf(`DELETE FROM %s WHERE %s = $1`, userprivacysettingsql.Table, userprivacysettingsql.UserId)
	_, err := d.tx.Exec(query, d.userId)
	if err != nil {
		d.err = Error(err, exception.NewCause("clear privacy settings", "PrivacySettings", _PROVIDER))
	}
}

func (s *Storage) ResetUser(ctx context.Context, userId int64) error {
	tx, err := s.db.Begin()
	if err != nil {
		return Error(err, exception.NewCause("begin tx", "ResetUser", _PROVIDER))
	}
	defer tx.Rollback()
	clear := clearUser{userId: userId, tx: tx}
	r := clear.MarkDeleted(ctx)
	if r == 0 {
		return exception.Wrap(usererror.ExceptionUserNotFound(), exception.NewCause("delete user", "ResetUser", _PROVIDER))
	}
	clear.Achievements(ctx)
	clear.PrivacySettings(ctx)
	if clear.err != nil {
		return Error(err, exception.NewCause("clear", "ResetUser", _PROVIDER))
	}
	err = tx.Commit()
	if err != nil {
		return Error(err, exception.NewCause("commit tx", "ResetUser", _PROVIDER))
	}
	return nil
}

type recoverUserData struct {
	userId int64
	user   *usermodel.CreateUser
}

func (r recoverUserData) Query() string {
	return fmt.Sprintf(
		`UPDATE %s SET %s = $2, %s = $3, %s = $4, %s = $5, %s = FALSE, %s = now() WHERE %s = $1 AND %s`,
		usersql.Table,

		usersql.Name,
		usersql.CigaretteDayAmount,
		usersql.CigarettePackAmount,
		usersql.PackPrice,
		usersql.Deleted,
		usersql.AbstinenceTime,

		usersql.ID,
		usersql.Deleted,
	)
}

func (r recoverUserData) Exec(ctx context.Context, tx *sql.Tx) (int64, error) {
	res, err := tx.ExecContext(ctx, r.Query(), r.userId, r.user.Name, r.user.CigaretteDayAmount, r.user.CigarettePackAmount, r.user.PackPrice)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

type user struct {
	userId int64
	tx     *sql.Tx
	err    error
}

func (u user) Err() error {
	return u.err
}

func (u user) LevelQuery() string {
	return fmt.Sprintf(
		`SELECT %s, %s, %s, %s FROM %s WHERE %s = 1`,
		levelsql.Level,
		levelsql.Rank,
		levelsql.MinExp,
		levelsql.MaxExp,
		levelsql.Table,
		levelsql.Level,
	)
}

func (u *user) Level(ctx context.Context) usermodel.LevelInfo {
	if u.err != nil {
		return usermodel.LevelInfo{}
	}
	level := usermodel.LevelInfo{Level: 1}
	err := u.tx.QueryRowContext(ctx, u.LevelQuery()).Scan(
		&level.Level,
		&level.Rank,
		&level.MinExp,
		&level.MaxExp,
	)
	if err != nil {
		u.err = Error(err, exception.NewCause("level", "Level", _PROVIDER))
	}
	return level
}

// func (u user) SubscriptionQuery() string {
// 	return fmt.Sprintf(
// 		`SELECT %s,%s FROM %s INNER JOIN %s ON %s = %s WHERE %s = $1 GROUP BY %s,%s`,
// 		sqlutils.Full(usersubscriptionsql.Expired),
// 		sqlutils.Full(subscriptiontypesql.Type),

// 		usersubscriptionsql.Table,

// 		subscriptiontypesql.Table,
// 		sqlutils.Full(usersubscriptionsql.TypeId),
// 		sqlutils.Full(subscriptiontypesql.ID),

// 		sqlutils.Full(usersubscriptionsql.UserId),

// 		sqlutils.Full(usersubscriptionsql.Expired),
// 		sqlutils.Full(subscriptiontypesql.Type),
// 	)
// }

// func (u *user) Subscription(ctx context.Context) usermodel.Subscription {
// 	if u.err != nil {
// 		return usermodel.Subscription{}
// 	}
// 	var subscription usermodel.Subscription
// 	err := u.tx.QueryRowContext(ctx, u.SubscriptionQuery(), u.userId).Scan(&subscription.Expired, &subscription.Type)
// 	if err != nil {
// 		u.err = Error(err, exception.NewCause("subscription", "Subscription", _PROVIDER))
// 	}
// 	return subscription
// }

func (u *user) TriggersQuery() string {
	return fmt.Sprintf(
		`SELECT array_agg(%s) FROM %s INNER JOIN %s ON %s = %s WHERE %s = $1`,
		sqlutils.Full(triggersql.Name),

		usertriggersql.Table,

		triggersql.Table,
		sqlutils.Full(usertriggersql.TriggerId),
		sqlutils.Full(triggersql.ID),

		sqlutils.Full(usertriggersql.UserId),
	)
}

func (u *user) Triggers(ctx context.Context) []usermodel.Trigger {
	if u.err != nil {
		return nil
	}
	triggers := make([]usermodel.Trigger, 0)
	err := u.tx.QueryRowContext(ctx, u.TriggersQuery(), u.userId).Scan(pq.Array(&triggers))
	if err != nil {
		u.err = Error(err, exception.NewCause("triggers", "Triggers", _PROVIDER))
	}
	return triggers
}

func (u *user) MotivationsQuery() string {
	return fmt.Sprintf(`
		SELECT %s FROM %s 
		INNER JOIN %s ON %s = %s 
		INNER JOIN %s ON %s = %s 
		WHERE %s = $1
		GROUP BY %s
	`,
		sqlutils.Full(
			motivationsql.Motivation,
			welcomemotivationsql.Motivation,
		),
		usersql.Table,

		motivationsql.Table,
		sqlutils.Full(usersql.MotivationId),
		sqlutils.Full(motivationsql.ID),

		welcomemotivationsql.Table,
		sqlutils.Full(usersql.WelcomeMotivationId),
		sqlutils.Full(welcomemotivationsql.ID),

		sqlutils.Full(usersql.ID),

		sqlutils.Full(
			motivationsql.Motivation,
			welcomemotivationsql.Motivation,
		),
	)
}

func (u *user) Motivations(ctx context.Context) (motivation, welcomeMotivation string) {
	if u.err != nil {
		return
	}
	err := u.tx.QueryRowContext(ctx, u.MotivationsQuery(), u.userId).Scan(&motivation, &welcomeMotivation)
	if err != nil {
		u.err = Error(err, exception.NewCause("get user motivation", "Motivations", _PROVIDER))
	}
	return
}

func (s *Storage) RecoverUser(ctx context.Context, userId int64, createUser *usermodel.CreateUser) (*usermodel.UserData, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, Error(err, exception.NewCause("begin tx", "RecoverUser", _PROVIDER))
	}
	defer tx.Rollback()

	r, err := recoverUserData{userId: userId, user: createUser}.Exec(ctx, tx)
	if err != nil {
		return nil, Error(err, exception.NewCause("recover user data", "RecoverUser", _PROVIDER))
	}
	if r != 1 {
		return nil, Error(usererror.ExceptionUserNotFound(), exception.NewCause("many than one rows updated", "RecoverUser", _PROVIDER))
	}
	u := user{userId: userId, tx: tx}
	level := u.Level(ctx)
	triggers := u.Triggers(ctx)
	motivation, welcomeMotivation := u.Motivations(ctx)
	if u.Err() != nil {
		return nil, Error(u.Err(), exception.NewCause("get user data", "RecoverUser", _PROVIDER))
	}
	err = tx.Commit()
	if err != nil {
		return nil, Error(err, exception.NewCause("commit tx", "RecoverUser", _PROVIDER))
	}
	user := usermodel.NewUserData(string(createUser.Name), uint8(createUser.CigaretteDayAmount), uint8(createUser.CigarettePackAmount), float32(createUser.PackPrice), motivation, welcomeMotivation, time.Now(), level, triggers)
	return user, nil
}
