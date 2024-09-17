package authentication

import (
	"sso/internal/domain"
	"time"
)

type Authentication struct {
	ctx      *domain.Context
	store    *AuthenticationStoreProvider
	tokenTTL time.Duration
}

// New returns a new instance of Auth
func New(
	ctx *domain.Context,
	storeProvider *AuthenticationStoreProvider,
) *Authentication {
	return &Authentication{
		ctx:      ctx,
		store:    storeProvider,
		tokenTTL: ctx.Config().TokenTTL,
	}
}
