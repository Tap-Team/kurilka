package app

import (
	"database/sql"
	"log"

	"github.com/Tap-Team/kurilka/internal/config"
	"github.com/Tap-Team/kurilka/pkg/amidsql"
)

func Postgres(
	cnf config.PostgresConfig,
) *sql.DB {
	db, err := amidsql.NewPostgres(
		cnf.URL(),
		func(d *sql.DB) {
			d.SetMaxOpenConns(50)
		},
	)
	if err != nil {
		log.Fatalf("failed connect to postgres, %s", err)
	}
	return db
}
