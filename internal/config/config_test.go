package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("read config", func(t *testing.T) {
		cfg, err := NewAppConfig("./config_test.yaml")
		require.NoError(t, err)
		require.Equal(t, "warn", cfg.Logger.Level)
		require.Equal(t, "127.0.0.1", cfg.Api.Host)
	})

	t.Run("default config", func(t *testing.T) {
		cfg, err := NewAppConfig("./config_test_err.yaml")
		require.NoError(t, err)
		require.Equal(t, "info", cfg.Logger.Level)
	})

	t.Run("reading config error", func(t *testing.T) {
		_, err := NewAppConfig("./config_test_not_found.yaml")
		require.ErrorIs(t, err, ErrUnreadableConfig)
	})
}
