package config

type Config struct {
	database    *string
	storagePath string
	verbose     bool
}

type ConfigFunc func(config *Config)

func New(database *string, configs ...ConfigFunc) (config *Config) {
	config = &Config{
		database: database,
	}

	for _, configFunc := range configs {
		configFunc(config)
	}

	return
}

func (config *Config) GetDatabase() *string {
	return config.database
}
