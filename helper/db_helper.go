package helper

import "database/sql"

type DBTX interface {
	QueryRow(query string, args ...any) *sql.Row
}

func CoulmnValueExists(db DBTX, table string, column string, value any) (bool, error) {
	var exists bool
	err := db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM users WHERE "+column+" = ?)",
		value,
	).Scan(&exists)
	return exists, err
}
