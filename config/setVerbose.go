package config

func SetVerbose() ConfigFunc {
	return func(config *Config) {
		config.Verbose = true
	}
}
