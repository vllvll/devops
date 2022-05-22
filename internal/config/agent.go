package config

import (
	"github.com/caarlos0/env/v6"
	flag "github.com/spf13/pflag"
	"time"
)

type AgentConfig struct {
	Address        string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	Key            string        `env:"KEY"`
}

func CreateAgentConfig() (*AgentConfig, error) {
	var cfg AgentConfig

	flag.StringVarP(&cfg.Address, "address", "a", "127.0.0.1:8080", "Address. Format: ip:port (for example: 127.0.0.1:8080")
	flag.DurationVarP(&cfg.ReportInterval, "report", "r", 10*time.Second, "Report interval. Format: any input valid for time.ParseDuration (for example: 1s)")
	flag.DurationVarP(&cfg.PollInterval, "poll", "p", 2*time.Second, "Poll interval. Format: any input valid for time.ParseDuration (for example: 1s)")
	flag.StringVarP(&cfg.Key, "key", "k", "", "Key. Format: string (for example: ?)")

	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c AgentConfig) AddressWithHTTP() string {
	return "http://" + c.Address
}
