package config

type Config struct {
	Database    *string
	StoragePath string
	Verbose     bool
}

type ConfigFunc func(config *Config)

func New(database *string, configs ...ConfigFunc) (config Config) {
	config = Config{
		Database: database,
	}

	for _, configFunc := range configs {
		configFunc(&config)
	}

	return
}
