package config

import "time"

type Authentication struct {
	TokenTTL time.Duration `yaml:"token_ttl"`
}

type Services struct {
	Authentication Authentication `yaml:"authentication"`
}
