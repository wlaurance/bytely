package database

import (
	"database/sql"
	_ "github.com/lib/pq"
)

var DB_URL string = "user=rocconicosia dbname=linkthing sslmode=disable"

// WithDatabase is a small context manager that helps to remove all of the
// bothersome boilerplate database connection code from all of the functions
// that have to deal with database queries.
// A connection is opened, passed into the function provided, and closed
// after your function returns.
func WithDatabase(fn func(*sql.DB) error) error {
	db, err := sql.Open("postgres", DB_URL)
	if err != nil {
		return err
	}
	defer db.Close()

	return fn(db)
}
