package userworker_test

import (
	"bytes"
	"database/sql"
	"fmt"

	"github.com/Tap-Team/kurilka/internal/sqlmodel/usersql"
)

func insertUsers(db *sql.DB, amount int64) error {
	values := &bytes.Buffer{}
	for i := int64(0); i < amount; i++ {
		values.WriteString(fmt.Sprintf("(%d, '', 0,0,0)", i))
		if i < amount-1 {
			values.WriteRune(',')
		}
	}
	query := fmt.Sprintf(`
		INSERT INTO %s (%s, %s, %s, %s, %s) VALUES %s
	`,
		usersql.Table,

		usersql.ID,
		usersql.Name,
		usersql.CigaretteDayAmount,
		usersql.CigarettePackAmount,
		usersql.PackPrice,

		values.String(),
	)
	_, err := db.Exec(query)
	return err
}
