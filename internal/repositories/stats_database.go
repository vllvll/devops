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
		log.Fatal(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	if rows != 1 {
		log.Fatalf("expected single row affected, got %d rows affected", rows)
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
		log.Fatal(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	if rows != 1 {
		log.Fatalf("expected single row affected, got %d rows affected", rows)
	}
}

func (s *StatsDatabase) GetAll() (map[string]types.Gauge, map[string]types.Counter) {
	var gaugeCount, counterCount int64

	row := s.db.QueryRow("SELECT COUNT(*) as count FROM gauges")
	err := row.Scan(&gaugeCount)
	if err != nil {
		panic(err)
	}

	gauges := make(map[string]types.Gauge, gaugeCount)

	rows, err := s.db.Query("SELECT name, value FROM gauges")
	if err != nil || rows.Err() != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var value types.Gauge

		err = rows.Scan(&name, &value)
		if err != nil {
			log.Fatal(err)
		}

		gauges[name] = value
	}

	row = s.db.QueryRow("SELECT COUNT(*) as count FROM counters")
	err = row.Scan(&counterCount)
	if err != nil {
		panic(err)
	}

	counters := make(map[string]types.Counter, counterCount)

	rows, err = s.db.Query("SELECT name, value FROM counters")
	if err != nil || rows.Err() != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var value types.Counter

		err = rows.Scan(&name, &value)
		if err != nil {
			log.Fatal(err)
		}

		counters[name] = value
	}

	return gauges, counters
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
