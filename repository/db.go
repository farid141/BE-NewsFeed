package repository

import "database/sql"

// DBTX adalah interface yang abstrak database query execution.
// Interface ini kompatibel dengan *sql.DB dan *sql.Tx,
// sehingga repository methods bisa menerima keduanya tanpa if-else logic.
type DBTX interface {
	// Query menjalankan SQL query yang mengembalikan multiple rows.
	// Cocok untuk SELECT dengan hasil banyak data (contoh: GetFeed, GetUsers).
	// Mengembalikan *sql.Rows yang harus di-iterate dan di-close.
	Query(query string, args ...any) (*sql.Rows, error)

	// QueryRow menjalankan SQL query untuk single row result.
	// Cocok untuk SELECT dengan hasil satu data (contoh: GetUserByID, GetByUsername).
	// Menggunakan Scan() langsung tanpa iteration.
	QueryRow(query string, args ...any) *sql.Row

	// Exec menjalankan query yang tidak mengembalikan rows (INSERT, UPDATE, DELETE).
	// Mengembalikan sql.Result untuk akses LastInsertId() atau RowsAffected().
	// Cocok untuk mutasi data di database.
	Exec(query string, args ...any) (sql.Result, error)
}
