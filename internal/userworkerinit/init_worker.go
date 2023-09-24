package userworkerinit

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/Tap-Team/kurilka/internal/sqlmodel/usersql"
	"github.com/Tap-Team/kurilka/pkg/amidtime"
	"github.com/Tap-Team/kurilka/workers"
)

var selectUsersQuery = fmt.Sprintf(
	`SELECT %s, %s FROM %s WHERE NOT %s LIMIT $2 OFFSET $1`,
	usersql.ID,
	usersql.AbstinenceTime,
	usersql.Table,
	usersql.Deleted,
)

func addUsers(ctx context.Context, db *sql.DB, worker workers.UserWorker, offset, limit int) bool {
	rows, err := db.Query(selectUsersQuery, offset, limit)
	if err != nil {
		log.Fatalf("failed exec users query, %s", err)
	}
	defer rows.Close()
	count := 0
	for rows.Next() {
		var id int64
		var time amidtime.Timestamp
		err := rows.Scan(&id, &time)
		if err != nil {
			log.Fatalf("failed scan user data, %s", err)
		}
		worker.AddUser(ctx, workers.NewUser(id, time.Time))
		count++
	}
	return count == 0
}

func InitUserWorkerWorker(db *sql.DB, worker workers.UserWorker) {
	ctx := context.Background()
	offset, limit := 0, 10000
	ok := addUsers(ctx, db, worker, offset, limit)
	for !ok {
		offset += limit
		ok = addUsers(ctx, db, worker, offset, limit)
	}
}
