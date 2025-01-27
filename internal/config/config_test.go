package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/pingvincible/kvdatabase/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigNegative(t *testing.T) {
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
			wantError:       "",
		},
		{
			name:            "config path is set but config file does not exist",
			configPathIsSet: true,
			configFileExist: false,
			wantError:       "",
		},
		{
			name:            "config file is empty",
			configPathIsSet: true,
			configFileExist: true,
			configFileBody:  ``,
			wantError:       "",
		},
		{
			name:            "config file has invalid data",
			configPathIsSet: true,
			configFileExist: true,
			configFileBody:  `not a yaml file`,
			wantError:       "",
		},
	}

	t.Parallel()

	for _, tc := range testCases {

		testCase := tc

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			var configPath string

			configFile, err := os.CreateTemp("", "kvdatabase*.yaml")
			if err != nil {
				t.Fatal("failed to create temporary config file, %w", err)
			}
			defer func() { _ = os.Remove(configFile.Name()) }()

			if testCase.configPathIsSet {
				configPath = configFile.Name()
			}

			if !testCase.configFileExist {
				err = os.Remove(configFile.Name())
				if err != nil {
					t.Fatal("failed to remove config file, %w", err)
				}
			}

			if _, err = configFile.Write([]byte(testCase.configFileBody)); err != nil {
				t.Fatal("failed to write to temporary config file:", err)
			}

			cfg, err := config.Load(configPath)

			require.NoError(t, err)

			assert.Equal(t, "in-memory", cfg.Engine.Type)
			assert.Equal(t, "127.0.0.1:3223", cfg.Network.Address)
			assert.Equal(t, 100, cfg.Network.MaxConnections)
			assert.Equal(t, "4KB", cfg.Network.MaxMessageSize)
			assert.Equal(t, 5*time.Minute, cfg.Network.IdleTimeout)
			assert.Equal(t, "info", cfg.Logging.Level)
			assert.Equal(t, "./kvdatabase.log", cfg.Logging.Output)
		})
	}
}

func prepareConfigFile(configFileBody string) (*os.File, error) {
	configFile, err := os.CreateTemp("", "kvdatabase*.yaml")
	if err != nil {
		return nil, err
	}

	if _, err = configFile.Write([]byte(configFileBody)); err != nil {
		return nil, err
	}

	return configFile, err
}

func TestConfigFillAllValuesFromFile(t *testing.T) {
	t.Parallel()

	configFileBody := `
engine:
  type: "other_type"
network:
  address: "127.0.0.1:3225"
  max_connections: 1000
  max_message_size: "10KB"
  idle_timeout: 50m
logging:
  level: "debug"
  output: "./debug.log"`

	configFile, err := prepareConfigFile(configFileBody)
	if err != nil {
		t.Fatal("failed to prepare config file, %w", err)
	}

	defer func() { _ = os.Remove(configFile.Name()) }()

	cfg, err := config.Load(configFile.Name())

	require.NoError(t, err)

	assert.Equal(t, "other_type", cfg.Engine.Type)
	assert.Equal(t, "127.0.0.1:3225", cfg.Network.Address)
	assert.Equal(t, 1000, cfg.Network.MaxConnections)
	assert.Equal(t, "10KB", cfg.Network.MaxMessageSize)
	assert.Equal(t, 50*time.Minute, cfg.Network.IdleTimeout)
	assert.Equal(t, "debug", cfg.Logging.Level)
	assert.Equal(t, "./debug.log", cfg.Logging.Output)

}

func TestConfigFillAbsentSectionWithDefaultValues(t *testing.T) {
	t.Parallel()

	configFileBody := `
engine:
  type: "other_type"
logging:
  level: "debug"
  output: "./debug.log"`

	configFile, err := prepareConfigFile(configFileBody)
	if err != nil {
		t.Fatal("failed to prepare config file, %w", err)
	}

	defer func() { _ = os.Remove(configFile.Name()) }()

	cfg, err := config.Load(configFile.Name())

	require.NoError(t, err)

	assert.Equal(t, "other_type", cfg.Engine.Type)
	assert.Equal(t, "127.0.0.1:3223", cfg.Network.Address)
	assert.Equal(t, 100, cfg.Network.MaxConnections)
	assert.Equal(t, "4KB", cfg.Network.MaxMessageSize)
	assert.Equal(t, 5*time.Minute, cfg.Network.IdleTimeout)
	assert.Equal(t, "debug", cfg.Logging.Level)
	assert.Equal(t, "./debug.log", cfg.Logging.Output)
}

func TestConfigFillAbsentKeysWithDefaultValues(t *testing.T) {
	t.Parallel()

	configFileBody := `
engine:
network:
  address: "127.0.0.1:3225"
  idle_timeout: 50m
logging:
  level: "debug"`

	configFile, err := prepareConfigFile(configFileBody)
	if err != nil {
		t.Fatal("failed to prepare config file, %w", err)
	}

	defer func() { _ = os.Remove(configFile.Name()) }()

	cfg, err := config.Load(configFile.Name())

	require.NoError(t, err)

	assert.Equal(t, "in-memory", cfg.Engine.Type)
	assert.Equal(t, "127.0.0.1:3225", cfg.Network.Address)
	assert.Equal(t, 100, cfg.Network.MaxConnections)
	assert.Equal(t, "4KB", cfg.Network.MaxMessageSize)
	assert.Equal(t, 50*time.Minute, cfg.Network.IdleTimeout)
	assert.Equal(t, "debug", cfg.Logging.Level)
	assert.Equal(t, "./kvdatabase.log", cfg.Logging.Output)
}
