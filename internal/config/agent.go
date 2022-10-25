// Package config служит для получения параметров запуска агента и сервера
package config

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/caarlos0/env/v6"

	flag "github.com/spf13/pflag"
)

type fileConfig struct {
	JSONConfig string `env:"CONFIG"` // Путь до файла json с конфигом
}

type jsonAgentConfig struct {
	Address        string `json:"address"`
	ReportInterval string `json:"report_interval"`
	PollInterval   string `json:"poll_interval"`
	CryptoKey      string `json:"crypto_key"`
}

type AgentConfig struct {
	Address        string        `env:"ADDRESS"`         // Адрес для отправки значений
	ReportInterval time.Duration `env:"REPORT_INTERVAL"` // Периодичность отправки значений на сервер
	PollInterval   time.Duration `env:"POLL_INTERVAL"`   // Периодичность получения значений
	Key            string        `env:"KEY"`             // Ключ шифрования сообщений
	CryptoKey      string        `env:"CRYPTO_KEY"`      // Путь до файла с публичным ключом
}

// CreateAgentConfig возвращает структуру конфига AgentConfig со значениями для работы агента.
// Значения для конфига задаются через флаги или переменные окружения
// Приоритет значений у переменных окружения
func CreateAgentConfig() (*AgentConfig, error) {
	var config AgentConfig
	var jsonFileConfig fileConfig
	var jsonConfig = jsonAgentConfig{
		Address:        "127.0.0.1:8080",
		ReportInterval: "10s",
		PollInterval:   "2s",
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

	reportInterval, err := time.ParseDuration(jsonConfig.ReportInterval)
	if err != nil {
		return nil, err
	}

	pollInterval, err := time.ParseDuration(jsonConfig.PollInterval)
	if err != nil {
		return nil, err
	}

	flag.StringVarP(&config.Address, "address", "a", jsonConfig.Address, "Address. Format: ip:port (for example: 127.0.0.1:8080")
	flag.DurationVarP(&config.ReportInterval, "report", "r", reportInterval, "Report interval. Format: any input valid for time.ParseDuration (for example: 1s)")
	flag.DurationVarP(&config.PollInterval, "poll", "p", pollInterval, "Poll interval. Format: any input valid for time.ParseDuration (for example: 1s)")
	flag.StringVarP(&config.Key, "key", "k", "", "Key. Format: string (for example: ?)")
	flag.StringVarP(&config.CryptoKey, "crypto-key", "y", jsonConfig.CryptoKey, "Path for public key")

	flag.Parse()

	err = env.Parse(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// AddressWithHTTP получение адреса с http префиксом
func (c AgentConfig) AddressWithHTTP() string {
	return "http://" + c.Address
}

func (c AgentConfig) GetServiceIP() (ip string, err error) {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return ip, err
	}

	for _, a := range addresses {
		if ipNet, ok := a.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			ip = ipNet.IP.String()
		}
	}

	if ip == "" {
		return ip, fmt.Errorf("ip адрес не найден")
	}

	return ip, nil
}
