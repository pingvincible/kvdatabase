package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Flags struct {
	EngineType     *string
	Address        *string
	MaxConnections *int
	MaxMessageSize *string
	IdleTimeout    *time.Duration
	LoggingLevel   *string
	LoggingOutput  *string
}

type Config struct {
	Engine  EngineConfig  `yaml:"engine" env-description:"database engine configuration"`
	Network NetworkConfig `yaml:"network" env-description:"network configuration"`
	Logging LogConfig     `yaml:"logging" env-description:"logging configuration"`
}

type EngineConfig struct {
	Type string `yaml:"type" env:"ENGINE_TYPE" env-default:"in-memory" env-description:"database engine type"`
}

type NetworkConfig struct {
	Address        string        `yaml:"address" env:"ADDRESS" env-default:"127.0.0.1:3223" env-description:"address to listen"`
	MaxConnections int           `yaml:"maxConnections" env:"MAX_CONNECTIONS" env-default:"100" env-description:"max client connections"`
	MaxMessageSize string        `yaml:"maxMessageSize" env:"MAX_MESSAGE_SIZE" env-default:"4KB" env-description:"max message size"`
	IdleTimeout    time.Duration `yaml:"idleTimeout" env:"IDLE_TIMEOUT" env-default:"5m" env-description:"idle timeout"`
}

type LogConfig struct {
	Level  string `yaml:"level" env:"LOG_LEVEL" env-default:"info" env-description:"log level"`
	Output string `yaml:"output" env:"LOG_OUTPUT" env-default:"./kvdatabase.log" env-description:"log output filename"`
}

func Load(configPath string) (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("failed to load config file: %w", err)
	}

	return &cfg, nil
}

func (c *Config) UpdateWithFlags(flags Flags) {
	c.Engine.Type = *flags.EngineType
	c.Network.Address = *flags.Address
	c.Network.MaxConnections = *flags.MaxConnections
	c.Network.MaxMessageSize = *flags.MaxMessageSize
	c.Network.IdleTimeout = *flags.IdleTimeout
	c.Logging.Level = *flags.LoggingLevel
	c.Logging.Output = *flags.LoggingOutput
}
