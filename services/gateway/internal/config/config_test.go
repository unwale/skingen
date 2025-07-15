package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		os.Setenv("PORT", "8080")                    //nolint:errcheck
		os.Setenv("TASK_SERVICE_URL", "localhost:1") //nolint:errcheck

		cfg, err := LoadConfig()

		assert.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Equal(t, "8080", cfg.Port)
		assert.Equal(t, "localhost:1", cfg.TaskServiceUrl)
	})

	t.Run("missing env vars", func(t *testing.T) {
		os.Unsetenv("PORT")             //nolint:errcheck
		os.Unsetenv("TASK_SERVICE_URL") //nolint:errcheck

		cfg, err := LoadConfig()

		assert.Error(t, err)
		assert.Nil(t, cfg)
	})
}
