package postgres

import (
	"database/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func ConnectDatabase(dsn string) (*sql.DB, error) {
	var db *sql.DB

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
