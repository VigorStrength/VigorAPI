package config

import (
	"errors"

	"github.com/spf13/viper"
)

var (
	ErrMissingDBURI     = errors.New("missing VIGOR_DB_URI")
	ErrMissingDBNAME    = errors.New("missing VIGOR_DB_NAME")
	ErrMissingSecretKey = errors.New("missing JWT_SECRET_KEY")
)

type Config struct {
	MongoDBURI   string
	DatabaseName string
	JWTSecretKey string
}

func LoadConfig() (*Config, error) {
	viper.AutomaticEnv()

	// Set default values for development and test environments
	viper.SetDefault("VIGOR_DB_URI", "")
	viper.SetDefault("VIGOR_DB_NAME", "")
	viper.SetDefault("JWT_SECRET_KEY", "")

	// Check for test environment
	environment := viper.GetString("VIGOR_ENV")
	if environment == "test" {
		viper.Set("VIGOR_DB_URI", "mongodb://localhost:27017")
		viper.Set("VIGOR_DB_NAME", "Vigor_Test")
	} 
	if environment == "dev" {
		viper.Set("VIGOR_DB_URI", "mongodb://localhost:27017")
		viper.Set("VIGOR_DB_NAME", "Vigor_Dev")
		viper.SetDefault("JWT_SECRET_KEY", "your_default_secret")
	}

	// Retrieve the actual values considering environment variables
	mongoDBURI := viper.GetString("VIGOR_DB_URI")
	if mongoDBURI == "" {
		return nil, ErrMissingDBURI
	}

	databaseName := viper.GetString("VIGOR_DB_NAME")
	if databaseName == "" {
		return nil, ErrMissingDBNAME
	}

	jwtSecretKey := viper.GetString("JWT_SECRET_KEY")
	if jwtSecretKey == "" {
		return nil, ErrMissingSecretKey
	}

	config := &Config{
		MongoDBURI:   mongoDBURI,
		DatabaseName: databaseName,
		JWTSecretKey: jwtSecretKey,
	}

	return config, nil
}
