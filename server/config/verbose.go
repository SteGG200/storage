package config

func SetVerbose() ConfigFunc {
	return func(config *Config) {
		config.verbose = true
	}
}

func (config *Config) GetVerbose() bool {
	return config.verbose
}
