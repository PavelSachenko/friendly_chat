package config

type server struct {
	GatewayAddress string `mapstructure:"gateway_server"`
	CommandAddress string `mapstructure:"command_server"`
	QueryAddress   string `mapstructure:"query_server"`
}
