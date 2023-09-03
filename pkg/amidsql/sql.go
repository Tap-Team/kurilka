package amidsql

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type SQLOption func(*sql.DB)

func NewPostgres(dataSourceName string, opts ...SQLOption) (*sql.DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed open connection, %s", err)
	}
	for _, opt := range opts {
		opt(db)
	}
	return db, nil
}
