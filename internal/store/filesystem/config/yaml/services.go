package config

import (
	"sso/internal/lib/config"
	"time"
)

type Authentication struct {
	TokenTTL time.Duration `yaml:"token_ttl"`
}

type Services struct {
	Authentication Authentication `yaml:"authentication"`
}

func (services Services) getAuthenticationService() config.AuthenticationService {
	return config.AuthenticationService{
		TokenTTL: services.Authentication.TokenTTL,
	}
}
