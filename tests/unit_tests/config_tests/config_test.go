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
	os.Setenv("VIGOR_DB_NAME", "Vigor_Test")
	os.Setenv("JWT_SECRET_KEY", "VigorSuperSecretKey")

	// Unset environment variable after the test
	defer os.Unsetenv("VIGOR_DB_URI")
	defer os.Unsetenv("VIGOR_DB_NAME")
	defer os.Unsetenv("JWT_SECRET_KEY")

	want := &config.Config{
		MongoDBURI:   "mongodb://localhost:27017",
		DatabaseName: "Vigor_Test",
		JWTSecretKey: "VigorSuperSecretKey",
	}

	got, err := config.LoadConfig()
	assert.NoError(t, err, "LoadConfig() should not error")
	assert.Equal(t, want, got, "LoadConfig() should return the expected configuration.")
}

// func TestLoadConfigFailureMissingDBURI(t *testing.T) {
// 	os.Setenv("VIGOR_DB_NAME", "Vigor_Test")
// 	os.Setenv("JWT_SECRET_KEY", "VigorSuperSecretKey")
// 	os.Unsetenv("VIGOR_DB_URI")

// 	defer func() {
// 		os.Unsetenv("VIGOR_DB_NAME")
// 		os.Unsetenv("JWT_SECRET_KEY")
// 	}()

// 	_, err := config.LoadConfig(false)

// 	assert.Error(t, err, "LoadConfig() should error due to missing Database URI")
// 	// assert.Error(t, err, "LoadConfig() should error due to missing Database URI")
// 	// assert.ErrorIs(t, config.ErrMissingDBURI, err)
// }

// func TestLoadConfigFailureMissingDBName(t *testing.T) {
// 	os.Setenv("VIGOR_DB_URI", "mongodb://localhost:27017")
// 	os.Setenv("JWT_SECRET_KEY", "VigorSuperSecretKey")
// 	os.Unsetenv("VIGOR_DB_NAME")

// 	defer func() {
// 		os.Unsetenv("VIGOR_DB_URI")
// 		os.Unsetenv("JWT_SECRET_KEY")
// 	}()

// 	_, err := config.LoadConfig(false)

// 	assert.Error(t, err, "LoadConfig() should error due to missing Database URI")
// 	assert.Equal(t, config.ErrMissingDBNAME, err)
// }

func TestLoadConfigFailureMissingJWTSecretKey(t *testing.T) {
	os.Setenv("VIGOR_DB_URI", "mongodb://localhost:27017")
	os.Setenv("VIGOR_DB_NAME", "Vigor_Test")
	os.Unsetenv("JWT_SECRET_KEY")

	defer func() {
		os.Unsetenv("VIGOR_DB_URI")
		os.Unsetenv("VIGOR_DB_NAME")
	}()

	_, err := config.LoadConfig()

	assert.Error(t, err, "LoadConfig() should error due to missing secret key")
	assert.Equal(t, config.ErrMissingSecretKey, err)
}