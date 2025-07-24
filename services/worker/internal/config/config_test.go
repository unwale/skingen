package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		os.Setenv("PORT", "8080")                                 //nolint:errcheck
		os.Setenv("RABBITMQ_URL", "url")                          //nolint:errcheck
		os.Setenv("MODEL_SERVER_URL", "model_server_url")         //nolint:errcheck
		os.Setenv("IMAGE_BUCKET", "image_bucket")                 //nolint:errcheck
		os.Setenv("S3_ENDPOINT", "s3_endpoint")                   //nolint:errcheck
		os.Setenv("S3_ACCESS_KEY", "s3_access_key")               //nolint:errcheck
		os.Setenv("S3_SECRET_KEY", "s3_secret_key")               //nolint:errcheck
		os.Setenv("S3_BUCKET", "s3_bucket")                       //nolint:errcheck
		os.Setenv("GENERATE_IMAGE_QUEUE", "generate_image_queue") //nolint:errcheck
		os.Setenv("TASK_RESULT_QUEUE", "task_result_queue")       //nolint:errcheck
		cfg, err := LoadConfig()

		assert.NoError(t, err)
		assert.NotNil(t, cfg)
	})

	t.Run("missing env vars", func(t *testing.T) {
		os.Unsetenv("PORT") //nolint:errcheck

		cfg, err := LoadConfig()

		assert.Error(t, err)
		assert.Nil(t, cfg)
	})
}
