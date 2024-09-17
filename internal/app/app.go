package app

import (
	grpcApp "sso/internal/app/grpc"
	"sso/internal/domain"
	"sso/internal/services/authentication"
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

	// Инициализация auth сервиса
	authService := authentication.NewService(
		ctx,
		storage,
		storage,
		storage,
	)

	grpcapp := grpcApp.New(ctx, authService)

	return &App{
		GRPCServer: grpcapp,
	}

}
