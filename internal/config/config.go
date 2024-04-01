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

func LoadConfig(useDefaults bool) (*Config, error) {
	viper.AutomaticEnv()

	environment := viper.GetString("VIGOR_ENV")
	isTestEnv := environment == "test"

	if useDefaults || isTestEnv {
		viper.SetDefault("VIGOR_DB_URI", "mongodb://localhost:27017")
		if isTestEnv {
			viper.SetDefault("VIGOR_DB_NAME", "Vigor_Test")
		} else {
			viper.SetDefault("VIGOR_DB_NAME", "Vigor_Production")
		}
	}
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
