package helper

import (
	"database/sql"

	"github.com/sirupsen/logrus"
)

// WithTx executes a function within a database transaction.
// It automatically handles Begin, Commit, Rollback, and panic recovery.
func WithTx(db *sql.DB, fn func(*sql.Tx) error, logger *logrus.Logger) error {
	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		if logger != nil {
			logger.Error("Failed to begin transaction: ", err)
		}
		return err
	}

	// Defer rollback and panic recovery
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	// Execute the provided function
	err = fn(tx)
	if err != nil {
		tx.Rollback()
		if logger != nil {
			logger.Error("Transaction failed: ", err)
		}
		return err
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		if logger != nil {
			logger.Error("Failed to commit transaction: ", err)
		}
		return err
	}

	return nil
}
