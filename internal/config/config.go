package config

import (
	"errors"
	"os"
)

type Config struct {
	MongoDBURI string
	DatabaseName string
}

func LoadConfig() (*Config, error) {
	mongoDBURI := os.Getenv("VIGOR_DB_URI")
	databaseName := os.Getenv("VIGOR_DB_NAME")

	if mongoDBURI == "" || databaseName == "" {
		return nil, errors.New("database configuration not set")
	}

	config := &Config{
		MongoDBURI: mongoDBURI,
		DatabaseName: databaseName,
	}

	return config, nil
}
