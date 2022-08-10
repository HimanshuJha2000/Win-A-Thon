package config

// Database : struct to hold live / test Database config
type Database struct {
	Dialect            string `toml:"dialect"`
	Protocol           string `toml:"protocol"`
	Host               string `toml:"host"`
	Port               int    `toml:"port"`
	Username           string `env:"username"`
	Password           string `env:"password"`
	Name               string `toml:"name"`
	MaxOpenConnections int    `toml:"max_open_connections"`
	MaxIdleConnections int    `toml:"max_idle_connections"`
}
