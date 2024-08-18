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
	// log *slog.Logger,
	// grpcPort int,
	// storagePath string,
	// tokenTTL time.Duration,
) *App {

	// Инициализация хранилища
	storage, err := sqlite.New(ctx.Config().StorePath)
	if err != nil {
		panic(err)
	}

	authStoreProvider := auth.NewStoreProvider(
		storage,
		storage,
		storage,
	)

	// Инициализация auth сервиса
	authService := auth.New(ctx, *authStoreProvider)

	grpcapp := grpcApp.New(ctx.Log(), authService, ctx.Config().GRPC.Port)

	return &App{
		GRPCServer: grpcapp,
	}

}
