package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		os.Setenv("PORT", "8080")                  //nolint:errcheck
		os.Setenv("POSTGRES_HOST", "localhost")    //nolint:errcheck
		os.Setenv("POSTGRES_PORT", "5432")         //nolint:errcheck
		os.Setenv("POSTGRES_USER", "user")         //nolint:errcheck
		os.Setenv("POSTGRES_PASSWORD", "password") //nolint:errcheck
		os.Setenv("POSTGRES_DB", "skingen")        //nolint:errcheck

		cfg, err := LoadConfig()

		assert.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Equal(t, "8080", cfg.Port)
		assert.Equal(t, "localhost", cfg.PostgresHost)
		assert.Equal(t, "5432", cfg.PostgresPort)
		assert.Equal(t, "user", cfg.PostgresUser)
		assert.Equal(t, "password", cfg.PostgresPassword)
		assert.Equal(t, "skingen", cfg.PostgresDB)
	})

	t.Run("missing env vars", func(t *testing.T) {
		os.Unsetenv("PORT") //nolint:errcheck

		cfg, err := LoadConfig()

		assert.Error(t, err)
		assert.Nil(t, cfg)
	})
}
