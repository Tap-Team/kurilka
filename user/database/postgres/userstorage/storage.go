package userstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"log/slog"

	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/achievementsql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/levelsql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/motivationsql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/subscriptiontypesql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/triggersql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/userachievementsql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/usersql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/usersubscriptionsql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/usertriggersql"
	"github.com/Tap-Team/kurilka/internal/sqlmodel/welcomemotivationsql"
	"github.com/Tap-Team/kurilka/pkg/amidtime"
	"github.com/Tap-Team/kurilka/pkg/exception"
	"github.com/Tap-Team/kurilka/pkg/sqlutils"
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
		case usersql.PrimaryKey:
			return exception.Wrap(usererror.ExceptionUserExist(), cause)
		default:
			return exception.Wrap(err, cause)
		}
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
	triggerIdQuery = func(triggerName usermodel.Trigger) string {
		return fmt.Sprintf(`SELECT %s FROM %s WHERE %s = '%s'`, triggersql.ID, triggersql.Table, triggersql.Name, string(triggerName))
	}
)

type userInsert struct {
	tx         *sql.Tx
	err        error
	userId     int64
	createUser *usermodel.CreateUser
	user       usermodel.UserData
}

func (u userInsert) User() *usermodel.UserData {
	return &u.user
}
func (u userInsert) Err() error {
	return u.err
}

func (u userInsert) InsertUserQuery() string {
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

func (u *userInsert) InsertUser(ctx context.Context) {
	if u.err != nil {
		return
	}
	err := u.tx.QueryRowContext(
		ctx,
		u.InsertUserQuery(),
		u.userId,
		u.createUser.Name,
		u.createUser.CigaretteDayAmount,
		u.createUser.CigarettePackAmount,
		u.createUser.PackPrice,
	).Scan(&u.user.AbstinenceTime)
	if err != nil {
		u.err = Error(err, exception.NewCause("insert into user table", "InsertUser", _PROVIDER))
	}
	u.user.Name = u.createUser.Name
	u.user.PackPrice = u.createUser.PackPrice
	u.user.CigaretteDayAmount = u.createUser.CigaretteDayAmount
	u.user.CigarettePackAmount = u.createUser.CigarettePackAmount
}

func (u *userInsert) InsertSubscriptionQuery() string {
	return fmt.Sprintf(
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
}

func (u *userInsert) InsertSubscription(ctx context.Context, expired amidtime.Timestamp) {
	if u.err != nil {
		return
	}
	_, err := u.tx.ExecContext(ctx, u.InsertSubscriptionQuery(), u.userId, expired)
	if err != nil {
		u.err = Error(err, exception.NewCause("insert into subscription table", "InsertSubscription", _PROVIDER))
	}
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
		triggerIdQuery(usermodel.ENABLE_MESSAGES),
	)
}

func (u *userInsert) InsertTriggers(ctx context.Context) {
	if u.err != nil {
		return
	}
	_, err := u.tx.ExecContext(ctx, u.InsertTriggersQuery(), u.userId)
	if err != nil {
		u.err = Error(err, exception.NewCause("insert into triggers table", "InsertTriggers", _PROVIDER))
	}
	u.user.Triggers = []usermodel.Trigger{
		usermodel.SUPPORT_CIGGARETTE,
		usermodel.SUPPORT_HEALTH,
		usermodel.SUPPORT_TRIAL,
		usermodel.THANK_YOU,
	}
}

func (u *userInsert) SetLevel(ctx context.Context) {
	if u.err != nil {
		return
	}
	query := fmt.Sprintf(
		`SELECT %s,%s FROM %s WHERE %s = 1`,
		levelsql.Rank,
		levelsql.MaxExp,

		levelsql.Table,

		levelsql.Level,
	)
	err := u.tx.QueryRowContext(ctx, query).Scan(&u.user.Level.Rank, &u.user.Level.MaxExp)
	if err != nil {
		u.err = Error(err, exception.NewCause("get init level", "SetLevel", _PROVIDER))
	}
	u.user.Level.Level = 1
}

func (u *userInsert) MotivationsQuery() string {
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

func (u *userInsert) SetMotivations(ctx context.Context) {
	if u.err != nil {
		return
	}
	err := u.tx.QueryRowContext(ctx, u.MotivationsQuery(), u.userId).Scan(&u.user.Motivation, &u.user.WelcomeMotivation)
	if err != nil {
		u.err = exception.Wrap(err, exception.NewCause("get user motivation", "SetMotivations", _PROVIDER))
	}
}

func (s *Storage) InsertUser(ctx context.Context, userId int64, user *usermodel.CreateUser) (*usermodel.UserData, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, Error(err, exception.NewCause("begin tx", "InsertUser", _PROVIDER))
	}
	defer tx.Rollback()
	insert := userInsert{
		createUser: user,
		tx:         tx,
		userId:     userId,
	}
	insert.InsertUser(ctx)
	expired := amidtime.Timestamp{Time: time.Now().Add(s.trialPeriod)}
	insert.InsertSubscription(ctx, expired)
	insert.InsertTriggers(ctx)
	insert.SetLevel(ctx)
	insert.SetMotivations(ctx)
	if insert.Err() != nil {
		return nil, exception.Wrap(err, exception.NewCause("insert user", "InsertUser", _PROVIDER))
	}
	err = tx.Commit()
	if err != nil {
		return nil, Error(err, exception.NewCause("commit tx", "InsertUser", _PROVIDER))
	}
	return insert.User(), nil
}

var userExpQuery = fmt.Sprintf(
	`
	SELECT coalesce(
	(
		SELECT sum(%s) FROM %s
    	INNER JOIN %s ON %s = %s
		WHERE %s IS NOT NULL AND %s = $1
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

	sqlutils.Full(userachievementsql.OpenDate),
	sqlutils.Full(userachievementsql.UserId),
)

func (s *Storage) UserExp(ctx context.Context, userId int64) (int, error) {
	var exp int
	err := s.db.QueryRowContext(ctx, userExpQuery, userId).Scan(&exp)
	return exp, err
}

var selectUserQuery = fmt.Sprintf(
	`
	SELECT %s, array_agg(%s) FROM %s
	INNER JOIN %s ON %s = %s
	INNER JOIN %s ON %s = %s
	INNER JOIN %s ON %s = %s
	INNER JOIN %s ON %s = %s
	INNER JOIN %s ON %s <= $2 AND %s >= $2
	WHERE %s = $1 AND NOT %s
	GROUP BY %s
	`,
	sqlutils.Full(
		// userdata
		usersql.Name,
		usersql.CigaretteDayAmount,
		usersql.CigarettePackAmount,
		usersql.PackPrice,
		usersql.AbstinenceTime,
		// level data
		levelsql.Level,
		levelsql.Rank,
		levelsql.MinExp,
		levelsql.MaxExp,

		motivationsql.Motivation,
		welcomemotivationsql.Motivation,
	),
	sqlutils.Full(triggersql.Name),

	usersql.Table,

	usertriggersql.Table,
	sqlutils.Full(usersql.ID),
	sqlutils.Full(usertriggersql.UserId),

	triggersql.Table,
	sqlutils.Full(usertriggersql.TriggerId),
	sqlutils.Full(triggersql.ID),

	motivationsql.Table,
	sqlutils.Full(usersql.MotivationId),
	sqlutils.Full(motivationsql.ID),

	welcomemotivationsql.Table,
	sqlutils.Full(usersql.WelcomeMotivationId),
	sqlutils.Full(welcomemotivationsql.ID),

	levelsql.Table,
	sqlutils.Full(levelsql.MinExp),
	sqlutils.Full(levelsql.MaxExp),

	sqlutils.Full(usersql.ID),
	sqlutils.Full(usersql.Deleted),
	// group by
	sqlutils.Full(
		// userdata
		usersql.Name,
		usersql.CigaretteDayAmount,
		usersql.CigarettePackAmount,
		usersql.PackPrice,
		usersql.AbstinenceTime,
		// level data
		levelsql.Level,
		levelsql.Rank,
		levelsql.MinExp,
		levelsql.MaxExp,

		motivationsql.Motivation,
		welcomemotivationsql.Motivation,
	),
)

var userLevelQuery = fmt.Sprintf(
	`SELECT %s,%s,%s,%s FROM %s WHERE %s <= $1 AND %s >= $1`,
	levelsql.Level,
	levelsql.MinExp,
	levelsql.MaxExp,
	levelsql.Rank,

	levelsql.Table,

	levelsql.MinExp,
	levelsql.MaxExp,
)

func (s *Storage) UserLevel(ctx context.Context, userId int64) (*usermodel.LevelInfo, error) {
	exp, err := s.UserExp(ctx, userId)
	if err != nil {
		return nil, Error(err, exception.NewCause("get user exp", "UserLevel", _PROVIDER))
	}
	level := usermodel.LevelInfo{Exp: exp}
	err = s.db.QueryRowContext(ctx, userLevelQuery, exp).Scan(
		&level.Level,
		&level.MinExp,
		&level.MaxExp,
		&level.Rank,
	)
	if err != nil {
		return nil, Error(err, exception.NewCause("scan user level", "UserLevel", _PROVIDER))
	}
	return &level, nil
}

func (s *Storage) User(ctx context.Context, userId int64) (*usermodel.UserData, error) {
	exp, err := s.UserExp(ctx, userId)
	if err != nil {
		return nil, Error(err, exception.NewCause("get user exp", "User", _PROVIDER))
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
		&userData.AbstinenceTime,

		&userData.Level.Level,
		&userData.Level.Rank,
		&userData.Level.MinExp,
		&userData.Level.MaxExp,

		&userData.Motivation,
		&userData.WelcomeMotivation,
		pq.Array(&userData.Triggers),
	)
	if err != nil {
		return nil, Error(err, exception.NewCause("select user from database", "User", _PROVIDER))
	}
	return &userData, nil
}

var existsQuery = fmt.Sprintf(
	`
	SELECT %s FROM %s WHERE %s = ANY($1) AND NOT %s ORDER BY %s ASC
`,
	usersql.ID,
	usersql.Table,
	usersql.ID,
	usersql.Deleted,
	usersql.ID,
)

func (s *Storage) Exists(ctx context.Context, userIds []int64) []int64 {
	rows, err := s.db.QueryContext(ctx, existsQuery, pq.Int64Array(userIds))
	if err != nil {
		slog.Error(err.Error())
		return make([]int64, 0)
	}
	defer rows.Close()
	existsUsers := make([]int64, 0)
	for rows.Next() {
		var id int64
		err := rows.Scan(&id)
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		existsUsers = append(existsUsers, id)
	}
	return existsUsers
}

func (s *Storage) UserDeleted(ctx context.Context, userId int64) (bool, error) {
	query := fmt.Sprintf(
		`SELECT %s FROM %s WHERE %s = $1`,
		usersql.Deleted,
		usersql.Table,
		usersql.ID,
	)
	deleted := true
	err := s.db.QueryRowContext(ctx, query, userId).Scan(&deleted)
	if err != nil {
		return false, Error(err, exception.NewCause("get user deleted", "UserDeleted", _PROVIDER))
	}
	return deleted, nil
}
