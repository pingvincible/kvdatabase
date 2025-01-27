package config_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/pingvincible/kvdatabase/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigInvalidConfigFile(t *testing.T) {
	testCases := []struct {
		name            string
		configPathIsSet bool
		configFileExist bool
		configFileBody  string
		wantError       string
	}{
		{
			name:            "config path is not set",
			configPathIsSet: false,
			configFileExist: false,
			wantError:       "failed to load config file",
		},
		{
			name:            "config path is set but config file does not exist",
			configPathIsSet: true,
			configFileExist: false,
			wantError:       "failed to load config file",
		},
		{
			name:            "config file is empty",
			configPathIsSet: true,
			configFileExist: true,
			configFileBody:  ``,
			wantError:       "failed to load config file",
		},
		{
			name:            "config file has invalid data",
			configPathIsSet: true,
			configFileExist: true,
			configFileBody:  `not a yaml file`,
			wantError:       "failed to load config file",
		},
	}

	t.Parallel()

	for _, tc := range testCases {
		testCase := tc

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			var configPath string

			configFile, err := os.CreateTemp(t.TempDir(), "kvdatabase*.yaml")
			if err != nil {
				t.Fatal("failed to create temporary config file, %w", err)
			}

			if _, err = configFile.WriteString(testCase.configFileBody); err != nil {
				t.Fatal("failed to write to temporary config file:", err)
			}

			defer func() { _ = os.Remove(configFile.Name()) }()

			if !testCase.configFileExist {
				err = configFile.Close()
				if err != nil {
					t.Fatal("failed to close config file, %w", err)
				}
				err = os.Remove(configFile.Name())
				if err != nil {
					t.Fatal("failed to remove config file, %w", err)
				}
			}

			if testCase.configPathIsSet {
				configPath = configFile.Name()
			}

			_, err = config.Load(configPath)

			assert.Contains(t, err.Error(), testCase.wantError)
		})
	}
}

func prepareConfigFile(configFileBody string) (*os.File, error) {
	configFile, err := os.CreateTemp("", "kvdatabase*.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp config file: %w", err)
	}

	if _, err = configFile.WriteString(configFileBody); err != nil {
		return nil, fmt.Errorf("failed to write data to temp config file: %w", err)
	}

	return configFile, nil
}

func TestConfigFillValues(t *testing.T) {
	testCases := []struct {
		name               string
		configFileBody     string
		wantType           string
		wantAddress        string
		wantMaxConnections int
		wantMaxMessageSize string
		wantIdleTimeout    time.Duration
		wantLevel          string
		wantOutput         string
	}{
		{
			name: "fill all values from file",
			configFileBody: `
engine:
  type: "other-type"
network:
  address: "127.0.0.1:3225"
  maxConnections: 1000
  maxMessageSize: "10KB"
  idleTimeout: 50m
logging:
  level: "debug"
  output: "./debug.log"`,
			wantType:           "other-type",
			wantAddress:        "127.0.0.1:3225",
			wantMaxConnections: 1000,
			wantMaxMessageSize: "10KB",
			wantIdleTimeout:    50 * time.Minute,
			wantLevel:          "debug",
			wantOutput:         "./debug.log",
		},
		{
			name: "fill absent section with default values",
			configFileBody: `
engine:
  type: "other-type"
logging:
  level: "debug"
  output: "./debug.log"`,
			wantType:           "other-type",
			wantAddress:        "127.0.0.1:3223",
			wantMaxConnections: 100,
			wantMaxMessageSize: "4KB",
			wantIdleTimeout:    5 * time.Minute,
			wantLevel:          "debug",
			wantOutput:         "./debug.log",
		},
		{
			name: "fill absent keys with default values",
			configFileBody: `
engine:
network:
  address: "127.0.0.1:3225"
  idleTimeout: 50m
logging:
  level: "debug"`,
			wantType:           "in-memory",
			wantAddress:        "127.0.0.1:3225",
			wantMaxConnections: 100,
			wantMaxMessageSize: "4KB",
			wantIdleTimeout:    50 * time.Minute,
			wantLevel:          "debug",
			wantOutput:         "./kvdatabase.log",
		},
	}
	t.Parallel()

	for _, tc := range testCases {
		testCase := tc

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			configFile, err := prepareConfigFile(testCase.configFileBody)

			if err != nil {
				t.Fatal("failed to prepare config file, %w", err)
			}

			defer func() { _ = os.Remove(configFile.Name()) }()

			cfg, err := config.Load(configFile.Name())

			require.NoError(t, err)

			assert.Equal(t, testCase.wantType, cfg.Engine.Type)
			assert.Equal(t, testCase.wantAddress, cfg.Network.Address)
			assert.Equal(t, testCase.wantMaxConnections, cfg.Network.MaxConnections)
			assert.Equal(t, testCase.wantMaxMessageSize, cfg.Network.MaxMessageSize)
			assert.Equal(t, testCase.wantIdleTimeout, cfg.Network.IdleTimeout)
			assert.Equal(t, testCase.wantLevel, cfg.Logging.Level)
			assert.Equal(t, testCase.wantOutput, cfg.Logging.Output)
		})
	}
}

func TestConfigUpdateWithFlags(t *testing.T) {
	flags := createFlags()
	configFileBody := `
engine:
  type: "other-type"
network:
  address: "127.0.0.1:3225"
  maxConnections: 1000
  maxMessageSize: "10KB"
  idleTimeout: 50m
logging:
  level: "debug"
  output: "./debug.log"`
	wantType := "flag-type"
	wantAddress := "flag-address"
	wantMaxConnections := 10000
	wantMaxMessageSize := "100KB"
	wantIdleTimeout := 500 * time.Minute
	wantLevel := "error"
	wantOutput := "./flag.log"

	t.Parallel()

	configFile, err := prepareConfigFile(configFileBody)

	if err != nil {
		t.Fatal("failed to prepare config file, %w", err)
	}

	defer func() { _ = os.Remove(configFile.Name()) }()

	cfg, err := config.Load(configFile.Name())

	require.NoError(t, err)

	cfg.UpdateWithFlags(*flags)

	assert.Equal(t, wantType, cfg.Engine.Type)
	assert.Equal(t, wantAddress, cfg.Network.Address)
	assert.Equal(t, wantMaxConnections, cfg.Network.MaxConnections)
	assert.Equal(t, wantMaxMessageSize, cfg.Network.MaxMessageSize)
	assert.Equal(t, wantIdleTimeout, cfg.Network.IdleTimeout)
	assert.Equal(t, wantLevel, cfg.Logging.Level)
	assert.Equal(t, wantOutput, cfg.Logging.Output)
}

func createFlags() *config.Flags {
	engineType := "flag-type"
	address := "flag-address"
	maxConnections := 10000
	maxMessageSize := "100KB"
	idleTimeout := 500 * time.Minute
	loggingLevel := "error"
	loggingOutput := "./flag.log"

	return &config.Flags{
		EngineType:     &engineType,
		Address:        &address,
		MaxConnections: &maxConnections,
		MaxMessageSize: &maxMessageSize,
		IdleTimeout:    &idleTimeout,
		LoggingLevel:   &loggingLevel,
		LoggingOutput:  &loggingOutput,
	}
}
