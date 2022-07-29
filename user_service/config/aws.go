package config

type aws struct {
	Region     string `mapstructure:"aws_region"`
	AccessKey  string `mapstructure:"aws_access_key_id"`
	SecretKey  string `mapstructure:"aws_secret_access_key"`
	BucketName string `mapstructure:"aws_bucket_name"`
}
