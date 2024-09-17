package auth

import (
	"context"
	"sso/internal/store/models"
)

type AuthenticationStoreProvider struct {
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
}

//go:generate go run github.com/vektra/mockery/v2@v2.45.0 --name=UserSaver
type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passwordHash []byte,
	) (userID string, err error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.45.0 --name=UserProvider
type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID string) (bool, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.45.0 --name=AppProvider
type AppProvider interface {
	App(ctx context.Context, appID string) (models.App, error)
}

// Сигнатура функции для задания параметров
type optsFunc func(*AuthenticationStoreProvider)

// Задает опции по умолчанию
func defaultOpts() *AuthenticationStoreProvider {
	return &AuthenticationStoreProvider{
		//
	}
}

// NewStoreProvider Возвращает новую структуру AuthenticationStoreProvider
// Принимает функции для задания опций в любом количестве
func NewAuthenticationStoreProvider(
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	opts ...optsFunc,
) *AuthenticationStoreProvider {
	o := defaultOpts()

	o.userSaver = userSaver
	o.userProvider = userProvider
	o.appProvider = appProvider

	for _, opt := range opts {
		opt(o)
	}
	return o
}

// WithUuid Опция задающая uuid
// func WithUuid(uuid string) optsFunc {
// 	return func(chart *AuthStoreProvider) {
// 		if len(uuid) != 0 {
// 			chart.uuid = uuid
// 		}
// 	}
// }
