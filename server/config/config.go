package config

import "github.com/SteGG200/storage/db"

type Config struct {
	database    *db.DB
	storagePath string
	verbose     bool
}

type ConfigFunc func(config *Config)

func New(database *db.DB, configs ...ConfigFunc) (config *Config) {
	config = &Config{
		database: database,
	}

	for _, configFunc := range configs {
		configFunc(config)
	}

	return
}

func (config *Config) GetDatabase() *db.DB {
	return config.database
}
