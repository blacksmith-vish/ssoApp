package auth

import (
	"context"
	"sso/internal/store/models"
)

type AuthStoreProvider struct {
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
	) (userID int64, err error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.45.0 --name=UserProvider
type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.45.0 --name=AppProvider
type AppProvider interface {
	App(ctx context.Context, appID int32) (models.App, error)
}

// Сигнатура функции для задания параметров
type optsFunc func(*AuthStoreProvider)

// Задает опции по умолчанию
func defaultOpts() *AuthStoreProvider {
	return &AuthStoreProvider{
		//
	}
}

// NewStoreProvider Возвращает новую структуру AuthStoreProvider
// Принимает функции для задания опций в любом количестве
func NewStoreProvider(
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	opts ...optsFunc,
) *AuthStoreProvider {
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
