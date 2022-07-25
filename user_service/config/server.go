package config

type server struct {
	Port        string `mapstructure:"server_port"`
	Host        string `mapstructure:"server_host"`
	GRPCAddress string `mapstructure:"server_grpc_address"`
}
