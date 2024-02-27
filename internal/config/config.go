package config

import (
	"errors"
	"os"
)

type Config struct {
	MongoDBURI string
}

func LoadConfig() (*Config, error) {
	mongoDBURI := os.Getenv("MONGODB_URI")
	if mongoDBURI == "" {
		return nil, errors.New("MONGODB_URI environment variable is not set")
	}

	config := &Config{
		MongoDBURI: mongoDBURI,
	}

	return config, nil
}
