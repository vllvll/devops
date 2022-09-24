// Package postgres Функционал для работы с PostgreSQL
package postgres

import (
	"database/sql"

	_ "github.com/jackc/pgx/v4/stdlib"
)

// ConnectDatabase Инициализация базы данных
func ConnectDatabase(dsn string) (*sql.DB, error) {
	var db *sql.DB

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS counters
			(
				id    uuid    NOT NULL
					CONSTRAINT counters_pk
						PRIMARY KEY,
				name  text    NOT NULL,
				value bigint  NOT NULL
			);
	
		CREATE UNIQUE INDEX IF NOT EXISTS counters_id_uindex
			ON counters (id);
	
		CREATE UNIQUE INDEX IF NOT EXISTS counters_name_uindex
			ON counters (name);
	
		CREATE TABLE IF NOT EXISTS gauges
		(
			id    uuid             NOT NULL
				CONSTRAINT gauges_pk
					PRIMARY KEY,
			name  text             NOT NULL,
			value double precision NOT NULL
		);
	
		CREATE UNIQUE INDEX IF NOT EXISTS gauges_id_uindex
			ON gauges (id);
	
		create unique index if not exists gauges_name_uindex
			ON gauges (name);
	`)

	if err != nil {
		return db, nil
	}

	return db, nil
}
