package config

type application struct {
	Name       string `toml:"app_name"`
	ListenPort int    `toml:"listen_port"`
	ListenIP   string `toml:"listen_ip"`
	Email      string `toml:"winathon_mail"`
	Password   string `toml:"winathon_password"`
}
