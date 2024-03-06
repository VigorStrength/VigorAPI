package config

import (
	"errors"

	"github.com/spf13/viper"
)

type Config struct {
	MongoDBURI   string
	DatabaseName string
}

func LoadConfig() (*Config, error){
	viper.AutomaticEnv()

	viper.SetDefault("VIGOR_DB_URI", "mongodb://localhost:27017")
	viper.SetDefault("VIGOR_DB_NAME", "Vigor_Production")

	mongoDBURI := viper.GetString("VIGOR_DB_URI")
	if mongoDBURI == "" {
		return nil, errors.New("missing VIGOR_DB_URI")
	}

	databaseName := viper.GetString("VIGOR_DB_NAME")
	if databaseName == "" {
		return nil, errors.New("missing VIGOR_DB_NAME")
	}

	config := &Config{
		MongoDBURI: mongoDBURI,
		DatabaseName: databaseName,
	}

	return config, nil
}

