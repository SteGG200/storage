package config

func SetStoragePath(path string) ConfigFunc {
	return func(config *Config) {
		config.storagePath = path
	}
}

func (config *Config) GetStoragePath() string {
	return config.storagePath
}
