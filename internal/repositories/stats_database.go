package repositories

import (
	"database/sql"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/vllvll/devops/internal/types"
	"log"
)

type StatsDatabase struct {
	db *sql.DB
}

func NewStatsDatabaseRepository(db *sql.DB) StatsRepository {
	return &StatsDatabase{
		db: db,
	}
}

func (s *StatsDatabase) UpdateGauge(key string, value types.Gauge) {
	id, _ := uuid.NewV4()

	result, err := s.db.Exec(
		"INSERT INTO gauges (id, name, value) VALUES ($1, $2, $3) ON CONFLICT (name) DO UPDATE SET value = excluded.value",
		id.String(),
		key,
		value,
	)
	if err != nil {
		log.Fatalf("Error with update gauge result: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("Error with affected gauge rows: %v", err)
	}

	if rows != 1 {
		log.Fatalf("Error with expected single row affected, got %d rows affected", rows)
	}
}

func (s *StatsDatabase) UpdateCount(key string, value types.Counter) {
	id, _ := uuid.NewV4()

	result, err := s.db.Exec(
		"INSERT INTO counters (id, name, value) VALUES ($1, $2, $3) ON CONFLICT (name) DO UPDATE SET value = counters.value + excluded.value;",
		id.String(),
		key,
		value,
	)
	if err != nil {
		log.Fatalf("Error with update counter result: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("Error with affected counter rows: %v", err)
	}

	if rows != 1 {
		log.Fatalf("Error with expected single row affected, got %d rows affected", rows)
	}
}

func (s *StatsDatabase) GetAll() (map[string]types.Gauge, map[string]types.Counter) {
	return s.getGauges(), s.getCounters()
}

func (s *StatsDatabase) GetGaugeByKey(key string) (types.Gauge, error) {
	var gauge types.Gauge

	row := s.db.QueryRow("SELECT value FROM gauges WHERE name = $1 LIMIT 1", key)
	err := row.Scan(&gauge)
	if err != nil {
		return types.Gauge(0), fmt.Errorf("%s key doesn't exists", key)
	}

	return gauge, nil
}

func (s *StatsDatabase) GetCounterByKey(key string) (types.Counter, error) {
	var counter types.Counter

	row := s.db.QueryRow("SELECT value FROM counters WHERE name = $1 LIMIT 1", key)
	err := row.Scan(&counter)
	if err != nil {
		return types.Counter(0), fmt.Errorf("%s key doesn't exists", key)
	}

	return counter, nil
}

func (s *StatsDatabase) UpdateAll(gauges types.Gauges, counters types.Counters) error {
	tx, err := s.db.Begin()
	if err != nil {
		log.Printf("Error with open transaction: %v\n", err)

		return err
	}

	stmtGauges, err := tx.Prepare("INSERT INTO gauges (id, name, value) VALUES ($1, $2, $3) ON CONFLICT (name) DO UPDATE SET value = excluded.value")
	if err != nil {
		log.Printf("Error with create prepared statement for gauge: %v\n", err)

		return err
	}

	stmtCounters, err := tx.Prepare("INSERT INTO counters (id, name, value) VALUES ($1, $2, $3) ON CONFLICT (name) DO UPDATE SET value = counters.value + excluded.value")
	if err != nil {
		return err
	}

	for key, value := range gauges {
		id, _ := uuid.NewV4()

		if _, err = stmtGauges.Exec(id.String(), key, value); err != nil {
			if err = tx.Rollback(); err != nil {
				log.Fatalf("Error with unable to rollback: %v", err)
			}

			return err
		}
	}

	for key, value := range counters {
		id, _ := uuid.NewV4()

		if _, err = stmtCounters.Exec(id.String(), key, value); err != nil {
			if err = tx.Rollback(); err != nil {
				log.Fatalf("Error with unable to rollback: %v", err)
			}

			return err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("Erro with unable to commit: %v", err)
	}

	return nil
}

func (s *StatsDatabase) getGauges() map[string]types.Gauge {
	var gaugeCount int64

	row := s.db.QueryRow("SELECT COUNT(*) as count FROM gauges")
	err := row.Scan(&gaugeCount)
	if err != nil {
		log.Fatalf("Error with get gauge count: %v", err)
	}

	gauges := make(map[string]types.Gauge, gaugeCount)

	rows, err := s.db.Query("SELECT name, value FROM gauges")
	if err != nil || rows.Err() != nil {
		log.Fatalf("Error with get gauge name and value: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var value types.Gauge

		err = rows.Scan(&name, &value)
		if err != nil {
			log.Fatalf("Error with scan gauge: %v", err)
		}

		gauges[name] = value
	}

	return gauges
}

func (s *StatsDatabase) getCounters() map[string]types.Counter {
	var counterCount int64

	row := s.db.QueryRow("SELECT COUNT(*) as count FROM counters")
	err := row.Scan(&counterCount)
	if err != nil {
		log.Fatalf("Error with get counter count: %v", err)
	}

	counters := make(map[string]types.Counter, counterCount)

	rows, err := s.db.Query("SELECT name, value FROM counters")
	if err != nil || rows.Err() != nil {
		log.Fatalf("Error with get counter name and value: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var value types.Counter

		err = rows.Scan(&name, &value)
		if err != nil {
			log.Fatalf("Error with scan counter: %v", err)
		}

		counters[name] = value
	}

	return counters
}
