package config

var (
	// config : this will hold all the application configuration
	config appConfig
)

type appConfig struct {
	Application application `toml:"application"`
	Database    Database    `toml:"database"`
}

func GetConfig() appConfig {
	return config
}
