package config

import (
	"time"

	"github.com/caarlos0/env/v6"
	flag "github.com/spf13/pflag"
)

type ServerConfig struct {
	Address       string        `env:"ADDRESS"`        // Адрес запуска HTTP-сервера
	StoreInterval time.Duration `env:"STORE_INTERVAL"` // Интервал времени в секундах, по истечении которого текущие показания сервера сбрасываются на диск
	StoreFile     string        `env:"STORE_FILE"`     // Имя файла, где хранятся значения
	Restore       bool          `env:"RESTORE"`        // Возможность восстановления данных с диска при запуске
	Key           string        `env:"KEY"`            // Ключ шифрования
	DatabaseDsn   string        `env:"DATABASE_DSN"`   // Адрес подключения к БД
}

// CreateServerConfig возвращает структуру конфига ServerConfig со значениями для работы сервера.
// Значения для конфига задаются через флаги или переменные окружения
// Приоритет значений у переменных окружения
func CreateServerConfig() (*ServerConfig, error) {
	var cfg ServerConfig

	flag.StringVarP(&cfg.Address, "address", "a", "127.0.0.1:8080", "Address. Format: ip:port (for example: 127.0.0.1:8080")
	flag.DurationVarP(&cfg.StoreInterval, "store", "i", 300*time.Second, "Store interval. Format: any input valid for time.ParseDuration (for example: 1s)")
	flag.StringVarP(&cfg.StoreFile, "file", "f", "/tmp/devops-metrics-db.json", "Store file. Format: local path (for example: /tmp/devops-metrics-db.json)")
	flag.BoolVarP(&cfg.Restore, "restore", "r", true, "Restore. Format: bool (for example: true")
	flag.StringVarP(&cfg.Key, "key", "k", "", "Key. Format: string (for example: ?)")
	flag.StringVarP(&cfg.DatabaseDsn, "database-dsn", "d", "", "Database dsn. Format: string (for example: postgres://username:password@localhost:5432/database_name)")

	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
