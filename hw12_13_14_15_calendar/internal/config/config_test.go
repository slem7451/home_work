package config

import (
	"testing"

	"github.com/stretchr/testify/require" //nolint:depguard
)

func TestConfig(t *testing.T) {
	configFile := "../../configs/config.toml"
	config := NewConfig(configFile)

	require.Equal(t, "sql", config.Storage)

	require.Equal(t, "INFO", config.Logger.Level)

	require.Equal(t, "localhost", config.DB.Host)
	require.Equal(t, 5432, config.DB.Port)
	require.Equal(t, "root", config.DB.User)
	require.Equal(t, "root", config.DB.Password)
	require.Equal(t, "app_db", config.DB.Name)

	require.Equal(t, "localhost", config.HTTP.Host)
	require.Equal(t, 8080, config.HTTP.Port)
}
