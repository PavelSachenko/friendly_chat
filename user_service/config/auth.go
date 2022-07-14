package config

import "time"

type auth struct {
	AuthAccessTokenSalt    string        `mapstructure:"auth_access_token_salt"`
	AuthRefreshTokenSalt   string        `mapstructure:"auth_refresh_token_salt"`
	AuthRefreshTokenExpire time.Duration `mapstructure:"auth_refresh_token_expire"`
	AuthAccessTokenExpire  time.Duration `mapstructure:"auth_access_token_expire"`
}
