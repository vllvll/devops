package config

import (
	"encoding/json"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
	flag "github.com/spf13/pflag"
)

type jsonServerConfig struct {
	Address       string `json:"address"`
	Restore       bool   `json:"restore"`
	StoreInterval string `json:"store_interval"`
	StoreFile     string `json:"store_file"`
	DatabaseDsn   string `json:"database_dsn"`
	CryptoKey     string `json:"crypto_key"`
}

type ServerConfig struct {
	Address       string        `env:"ADDRESS"`        // Адрес запуска HTTP-сервера
	StoreInterval time.Duration `env:"STORE_INTERVAL"` // Интервал времени в секундах, по истечении которого текущие показания сервера сбрасываются на диск
	StoreFile     string        `env:"STORE_FILE"`     // Имя файла, где хранятся значения
	Restore       bool          `env:"RESTORE"`        // Возможность восстановления данных с диска при запуске
	Key           string        `env:"KEY"`            // Ключ шифрования
	DatabaseDsn   string        `env:"DATABASE_DSN"`   // Адрес подключения к БД
	CryptoKey     string        `env:"CRYPTO_KEY"`     // Путь до файла с приватным ключом
}

// CreateServerConfig возвращает структуру конфига ServerConfig со значениями для работы сервера.
// Значения для конфига задаются через флаги или переменные окружения
// Приоритет значений у переменных окружения
func CreateServerConfig() (*ServerConfig, error) {
	var config ServerConfig
	var jsonFileConfig fileConfig
	var jsonConfig = jsonServerConfig{
		Address:       "127.0.0.1:8080",
		Restore:       true,
		StoreInterval: "300s",
		StoreFile:     "/tmp/devops-metrics-db.json",
	}

	jsonConfigFlag := flag.NewFlagSet("file", flag.ContinueOnError)
	jsonConfigFlag.StringVarP(&jsonFileConfig.JSONConfig, "config", "c", "", "JSON Config file")
	err := jsonConfigFlag.Parse([]string{"c"})
	if err != nil {
		return nil, err
	}

	err = env.Parse(&jsonFileConfig)
	if err != nil {
		return nil, err
	}

	if jsonFileConfig.JSONConfig != "" {
		content, err := os.ReadFile(jsonFileConfig.JSONConfig)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(content, &jsonConfig); err != nil {
			return nil, err
		}
	}

	storeInterval, err := time.ParseDuration(jsonConfig.StoreInterval)
	if err != nil {
		return nil, err
	}

	flag.StringVarP(&config.Address, "address", "a", jsonConfig.Address, "Address. Format: ip:port (for example: 127.0.0.1:8080")
	flag.DurationVarP(&config.StoreInterval, "store", "i", storeInterval, "Store interval. Format: any input valid for time.ParseDuration (for example: 1s)")
	flag.StringVarP(&config.StoreFile, "file", "f", jsonConfig.StoreFile, "Store file. Format: local path (for example: /tmp/devops-metrics-db.json)")
	flag.BoolVarP(&config.Restore, "restore", "r", jsonConfig.Restore, "Restore. Format: bool (for example: true")
	flag.StringVarP(&config.Key, "key", "k", "", "Key. Format: string (for example: ?)")
	flag.StringVarP(&config.DatabaseDsn, "database-dsn", "d", jsonConfig.DatabaseDsn, "Database dsn. Format: string (for example: postgres://username:password@localhost:5432/database_name)")
	flag.StringVarP(&config.CryptoKey, "crypto-key", "y", jsonConfig.CryptoKey, "Path for private key")

	flag.Parse()

	err = env.Parse(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
