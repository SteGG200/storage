package config

func SetStoragePath(path string) ConfigFunc {
	return func(config *Config) {
		config.StoragePath = path
	}
}
