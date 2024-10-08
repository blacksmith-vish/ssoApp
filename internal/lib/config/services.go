package config

import "time"

type Authentication struct {
	TokenTTL time.Duration `yaml:"token_ttl"`
}

func (auth Authentication) GetTokenTTL() time.Duration {
	return auth.TokenTTL
}

type Services struct {
	Authentication Authentication `yaml:"authentication"`
}
