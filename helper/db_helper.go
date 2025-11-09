package helper

import "database/sql"

func CoulmnValueExists(db *sql.DB, table string, column string, value any) (bool, error) {
	var exists bool
	err := db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM users WHERE "+column+" = ?)",
		value,
	).Scan(&exists)
	return exists, err
}
