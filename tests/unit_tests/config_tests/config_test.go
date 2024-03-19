/* trunk-ignore-all(golangci-lint/typecheck) */
package config_test

import (
	"os"
	"testing"

	"github.com/GhostDrew11/vigor-api/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfigSuccess(t *testing.T) {
	// Set environment variables for the test
	os.Setenv("VIGOR_DB_URI", "mongodb://localhost:27017")
	os.Setenv("VIGOR_DB_NAME", "Vigor_Production")
	os.Setenv("JWT_SECRET_KEY", "VigorSuperSecretKey")

	// Unset environment variable after the test
	defer os.Unsetenv("VIGOR_DB_URI")
	defer os.Unsetenv("VIGOR_DB_NAME")
	defer os.Unsetenv("JWT_SECRET_KEY")

	want := &config.Config{
		MongoDBURI:   "mongodb://localhost:27017",
		DatabaseName: "Vigor_Production",
		JWTSecretKey: "VigorSuperSecretKey",
	}

	got, err := config.LoadConfig()
	assert.NoError(t, err, "LoadConfig() should not error")
	assert.Equal(t, want, got, "LoadConfig() should return the expected configuration.")
}
