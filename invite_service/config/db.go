package config

type db struct {
	Driver   string `mapstructure:"db_driver"`
	Port     int    `mapstructure:"db_port"`
	Host     string `mapstructure:"db_host"`
	Password string `mapstructure:"db_password"`
	Username string `mapstructure:"db_username"`
	Database string `mapstructure:"db_database"`
	SSLMode  string `mapstructure:"db_ssl_mode"`

	RedisDSN string `mapstructure:"db_redis_dsn"`
}
