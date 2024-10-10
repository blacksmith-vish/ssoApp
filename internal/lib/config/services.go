package config

import "time"

type AuthenticationService struct {
	TokenTTL time.Duration
}

type Services struct {
	AuthenticationService AuthenticationService
}
