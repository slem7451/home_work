package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	configFile := "../../configs/config.toml"
	config := NewConfig(configFile)

	require.Equal(t, "sql", config.Storage)

	require.Equal(t, "INFO", config.Logger.Level)

	require.Equal(t, "localhost", config.Db.Host)
	require.Equal(t, 5432, config.Db.Port)
	require.Equal(t, "root", config.Db.User)
	require.Equal(t, "root", config.Db.Password)
	require.Equal(t, "app_db", config.Db.Name)

	require.Equal(t, "localhost", config.Http.Host)
	require.Equal(t, 8080, config.Http.Port)
}