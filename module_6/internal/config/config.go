package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
)

type Config struct {
	Env               string `env:"APP_ENV" envDefault:"development"`
	MongoCn           string `env:"MONGO_CN"`
	MongoDb           string `env:"MONGO_DB"`
	LoggerLevel       string `env:"LOGGER_LEVEL" envDefault:"info"`
	AuthSecret        string `env:"AUTH_SECRET"`
	TokenTTLMinutes   string `env:"TOKEN_TTL_MINUTES" envDefault:"300"`
	PyroscopeServer   string `env:"PYROSCOPE_SERVER"`
	EnableCPUProfiler string `env:"ENABLE_CPU_PROFILER" envDefault:"false"`
	InstanceName      string `env:"INSTANCE_NAME" envDefault:"dev"`
}

func NewConfigFromEnv() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return cfg, nil
}
