package app

import (
	grpcApp "sso/internal/app/grpc"
	"sso/internal/domain"
	"sso/internal/services/auth"
	"sso/internal/store/sqlite"
)

type App struct {
	GRPCServer *grpcApp.App
}

func New(
	ctx *domain.Context,
) *App {

	// Инициализация хранилища
	storage, err := sqlite.New(ctx.Config().StorePath)
	if err != nil {
		panic(err)
	}

	authStoreProvider := auth.NewAuthenticationStoreProvider(
		storage,
		storage,
		storage,
	)

	// Инициализация auth сервиса
	authService := auth.New(ctx, authStoreProvider)

	grpcapp := grpcApp.New(ctx, authService)

	return &App{
		GRPCServer: grpcapp,
	}

}
