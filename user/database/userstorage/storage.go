package userstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/achievementsql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/levelsql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/subscriptiontypesql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/userachievementsql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/userachievementstatussql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/usersql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/usersubscriptionsql"
	"github.com/Tap-Team/kurilka/pkg/exception"
	"github.com/Tap-Team/kurilka/pkg/sqlutils"
	"github.com/Tap-Team/kurilka/user/model/usermodel"
	"github.com/lib/pq"
)

const _PROVIDER = "user/database/userstorage"

type Storage struct {
	db          *sql.DB
	trialPeriod time.Duration
}

func New(db *sql.DB, trialPeriod time.Duration) *Storage {
	return &Storage{db: db, trialPeriod: trialPeriod}
}

func Error(err error, cause exception.Cause) error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Constraint {
		case usersubscriptionsql.ForeignKeyUsers:
			return exception.Wrap(usererror.ExceptionUserNotFound(), cause)
		}
		return exception.Wrap(err, cause)
	}
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return exception.Wrap(usererror.ExceptionUserNotFound(), cause)

	default:
		return exception.Wrap(err, cause)
	}
}

var (
	insertUserQuery = fmt.Sprintf(
		`INSERT INTO %s (%s, %s, %s, %s, %s) VALUES ($1, $2, $3, $4, $5)`,
		// insert into users
		usersql.Table,
		usersql.ID,
		usersql.Name,
		usersql.CigaretteDayAmount,
		usersql.CigarettePackAmount,
		usersql.PackPrice,
	)
	insertSubscriptionQuery = fmt.Sprintf(
		`
		INSERT INTO %s (%s, %s, %s) VALUES (
			$1,
			(SELECT %s FROM %s WHERE %s = 'TRIAL'),
			$2
		)
		`,
		// insert into user_subscriptions
		usersubscriptionsql.Table,
		usersubscriptionsql.UserId,
		usersubscriptionsql.TypeId,
		usersubscriptionsql.Expired,

		// select subscription type id
		subscriptiontypesql.ID,
		subscriptiontypesql.Table,
		subscriptiontypesql.Type,
	)
	selectInitLevel = fmt.Sprintf(
		`SELECT %s,%s FROM %s WHERE %s = $1`,
		levelsql.Rank,
		levelsql.MaxExp,

		levelsql.Table,

		levelsql.Level,
	)
)

func (s *Storage) InsertUser(ctx context.Context, userId int64, user *usermodel.CreateUser) (usermodel.UserData, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return usermodel.UserData{}, Error(err, exception.NewCause("begin tx", "InsertUser", _PROVIDER))
	}
	defer tx.Rollback()
	expired := time.Now().Add(s.trialPeriod)
	userData := usermodel.UserData{
		Name:                user.Name,
		CigaretteDayAmount:  user.CigaretteDayAmount,
		CigarettePackAmount: user.CigarettePackAmount,
		PackPrice:           user.PackPrice,
		Subscription:        usermodel.NewSubscription(usermodel.TRIAL, expired),
		Level: usermodel.LevelInfo{
			Level: usermodel.One,
		},
	}
	_, err = tx.ExecContext(
		ctx,
		insertUserQuery,
		userId,
		userData.Name,
		userData.CigaretteDayAmount,
		userData.CigarettePackAmount,
		userData.PackPrice,
	)
	if err != nil {
		return usermodel.UserData{}, Error(err, exception.NewCause("exec insert user query", "InsertUser", _PROVIDER))
	}
	_, err = tx.ExecContext(
		ctx,
		insertSubscriptionQuery,
		userId,
		userData.Subscription.Expired,
	)
	if err != nil {
		return usermodel.UserData{}, Error(err, exception.NewCause("exec insert user subscription", "InsertUser", _PROVIDER))
	}
	err = tx.QueryRowContext(ctx, selectInitLevel, userData.Level.Level).Scan(
		&userData.Level.Rank,
		&userData.Level.MaxExp,
	)
	if err != nil {
		return usermodel.UserData{}, Error(err, exception.NewCause("select init level", "InsertUser", _PROVIDER))
	}
	err = tx.Commit()
	if err != nil {
		return usermodel.UserData{}, Error(err, exception.NewCause("commit tx", "InsertUser", _PROVIDER))
	}
	return userData, nil
}

var deleteUserQuery = fmt.Sprintf(
	`
	DELETE FROM %s WHERE %s = $1
	`,
	usersql.Table,
	usersql.ID,
)

func (s *Storage) DeleteUser(ctx context.Context, userId int64) error {
	_, err := s.db.ExecContext(ctx, deleteUserQuery, userId)
	if err != nil {
		return Error(err, exception.NewCause("delete user", "DeleteUser", _PROVIDER))
	}
	return nil
}

var userExpQuery = fmt.Sprintf(
	`
	SELECT coalesce(
	(
		SELECT sum(%s) FROM %s
    	INNER JOIN %s ON %s = %s
    	INNER JOIN %s ON %s = %s AND %s = 'OPENED'
		WHERE %s = $1
    	GROUP BY %s
	),
	0
)
	`,
	sqlutils.Full(achievementsql.Exp),
	achievementsql.Table,

	// inner join
	userachievementsql.Table,
	sqlutils.Full(achievementsql.ID),
	sqlutils.Full(userachievementsql.AchievementId),

	// inner join
	userachievementstatussql.Table,
	sqlutils.Full(userachievementstatussql.ID),
	sqlutils.Full(userachievementsql.StatusId),
	// and user achievement status = OPENED
	sqlutils.Full(userachievementstatussql.Status),

	sqlutils.Full(userachievementsql.UserId),
	// group by achievement.exp
	sqlutils.Full(achievementsql.Exp),
)

func (s *Storage) UserExp(ctx context.Context, userId int64) (int, error) {
	var exp int
	err := s.db.QueryRowContext(ctx, userExpQuery, userId).Scan(&exp)
	if err != nil {
		return 0, err
	}
	return exp, nil
}

var selectUserQuery = fmt.Sprintf(
	`
	SELECT %s FROM %s
	INNER JOIN %s ON %s = %s
	INNER JOIN %s ON %s = %s
	INNER JOIN %s ON %s <= $2 AND %s > $2
	WHERE %s = $1
	GROUP BY %s
	`,
	sqlutils.Full(
		// userdata
		usersql.Name,
		usersql.CigaretteDayAmount,
		usersql.CigarettePackAmount,
		usersql.PackPrice,
		// level data
		levelsql.Level,
		levelsql.Rank,
		levelsql.MinExp,
		levelsql.MaxExp,
		// subscription data
		subscriptiontypesql.Type,
		usersubscriptionsql.Expired,
	),
	usersql.Table,

	// inner join usersubscriptions
	usersubscriptionsql.Table,
	sqlutils.Full(usersql.ID),
	sqlutils.Full(usersubscriptionsql.UserId),

	// inner join subscriptiontype
	subscriptiontypesql.Table,
	sqlutils.Full(usersubscriptionsql.TypeId),
	sqlutils.Full(subscriptiontypesql.ID),

	levelsql.Table,
	sqlutils.Full(levelsql.MinExp),
	sqlutils.Full(levelsql.MaxExp),

	sqlutils.Full(usersql.ID),
	// group by
	sqlutils.Full(
		usersql.ID,
		subscriptiontypesql.Type,
		usersubscriptionsql.TypeId,
		levelsql.Level,
		usersubscriptionsql.UserId,
	),
)

func (s *Storage) User(ctx context.Context, userId int64) (usermodel.UserData, error) {
	exp, err := s.UserExp(ctx, userId)
	if err != nil {
		return usermodel.UserData{}, Error(err, exception.NewCause("get user exp", "User", _PROVIDER))
	}
	var userData usermodel.UserData
	userData.Level.Exp = exp
	row := s.db.QueryRowContext(
		ctx,
		selectUserQuery,
		userId,
		exp,
	)
	err = row.Scan(
		&userData.Name,
		&userData.CigaretteDayAmount,
		&userData.CigarettePackAmount,
		&userData.PackPrice,

		&userData.Level.Level,
		&userData.Level.Rank,
		&userData.Level.MinExp,
		&userData.Level.MaxExp,

		&userData.Subscription.Type,
		&userData.Subscription.Expired,
	)
	if err != nil {
		return userData, Error(err, exception.NewCause("select user from database", "User", _PROVIDER))
	}
	return userData, nil
}
