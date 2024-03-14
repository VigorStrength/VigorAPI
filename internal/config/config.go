package config

import (
	"errors"

	"github.com/spf13/viper"
)

type Config struct {
	MongoDBURI   string
	DatabaseName string
	JWTSecretKey string
}

func LoadConfig() (*Config, error){
	viper.AutomaticEnv()

	viper.SetDefault("VIGOR_DB_URI", "mongodb://localhost:27017")
	viper.SetDefault("VIGOR_DB_NAME", "Vigor_Production")
	viper.SetDefault("JWT_SECRET_KEY", "VigorS3cr3tk3y#71124")

	mongoDBURI := viper.GetString("VIGOR_DB_URI")
	if mongoDBURI == "" {
		return nil, errors.New("missing VIGOR_DB_URI")
	}

	databaseName := viper.GetString("VIGOR_DB_NAME")
	if databaseName == "" {
		return nil, errors.New("missing VIGOR_DB_NAME")
	}

	jwtSecretKey := viper.GetString("JWT_SECRET_KEY")
	if jwtSecretKey == "" {
		return nil, errors.New("missing JWT_SECRET_KEY")
	}

	config := &Config{
		MongoDBURI: mongoDBURI,
		DatabaseName: databaseName,
		JWTSecretKey: jwtSecretKey,
	}

	return config, nil
}

