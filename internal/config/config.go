package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Engine  EngineConfig  `yaml:"engine"`
	Network NetworkConfig `yaml:"network"`
	Logging LogConfig     `yaml:"logging"`
}

type EngineConfig struct {
	Type string `yaml:"type" env-default:"in-memory"`
}

type NetworkConfig struct {
	Address        string        `yaml:"address" env-default:"127.0.0.1:3223"`
	MaxConnections int           `yaml:"maxConnections" env-default:"100"`
	MaxMessageSize string        `yaml:"maxMessageSize" env-default:"4KB"`
	IdleTimeout    time.Duration `yaml:"idleTimeout" env-default:"5m"`
}

type LogConfig struct {
	Level  string `yaml:"level" env-default:"info"`
	Output string `yaml:"output" env-default:"./kvdatabase.log"`
}

func Load(configPath string) (*Config, error) {
	var cfg Config

	if configPath == "" {
		err := cleanenv.ReadEnv(&cfg)
		if err != nil {
			return &cfg, fmt.Errorf("config path is not set, failed to load from env: %w", err)
		}

		return &cfg, nil
	}

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		errEnv := cleanenv.ReadEnv(&cfg)
		if errEnv != nil {
			return &cfg, fmt.Errorf("failed to read config file: %w, failed to load from env: %w", err, errEnv)
		}

		return &cfg, nil
	}

	return &cfg, nil
}
