package senderconfig

import (
	"testing"

	"github.com/stretchr/testify/require" //nolint:depguard
)

func TestConfig(t *testing.T) {
	configFile := "../../../configs/sender_config.toml"
	config := NewConfig(configFile)

	require.Equal(t, "sql", config.Storage)

	require.Equal(t, "INFO", config.Logger.Level)

	require.Equal(t, "postgres", config.DB.Host)
	require.Equal(t, 5432, config.DB.Port)
	require.Equal(t, "root", config.DB.User)
	require.Equal(t, "root", config.DB.Password)
	require.Equal(t, "app_db", config.DB.Name)

	require.Equal(t, "rabbitmq", config.Rabbit.Host)
	require.Equal(t, 5672, config.Rabbit.Port)
	require.Equal(t, "guest", config.Rabbit.User)
	require.Equal(t, "guest", config.Rabbit.Password)
	require.Equal(t, "app", config.Rabbit.Exchange)
	require.Equal(t, "calendar", config.Rabbit.Queue)
}
