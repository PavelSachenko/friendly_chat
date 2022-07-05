package config

type server struct {
	Port string `mapstructure:"server_port"`
	Host string `mapstructure:"server_host"`
}
