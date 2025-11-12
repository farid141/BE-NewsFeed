package db

import (
	"database/sql"
	"fmt"

	"github.com/farid141/go-rest-api/config"
	_ "github.com/lib/pq"
)

func NewDB(cfg config.Config) (*sql.DB, error) {
	// create db connection
	db_dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)
	db, err := sql.Open("mysql", db_dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
