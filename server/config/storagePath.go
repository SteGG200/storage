package config

/*
SetStoragePath sets the directory where to store the user's data

Params:

	path string // The directory path where to store the user's data.
*/
func SetStoragePath(path string) ConfigFunc {
	return func(config *Config) {
		config.storagePath = path
	}
}

func (config *Config) GetStoragePath() string {
	return config.storagePath
}
