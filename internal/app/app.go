package app

import (
	"log/slog"
	grpcApp "sso/internal/app/grpc"
	"sso/internal/lib/config"
	"sso/internal/services/authentication"
	"sso/internal/store/sqlite"
)

type App struct {
	GRPCServer *grpcApp.App
}

func NewApp(
	log *slog.Logger,
	conf *config.Config,
) *App {

	// Инициализация хранилища
	storage, err := sqlite.New(conf.StorePath)
	if err != nil {
		panic(err)
	}

	// Инициализация auth сервиса
	authService := authentication.NewService(
		log,
		storage,
		storage,
		storage,
		conf.Services.Authentication.TokenTTL,
	)

	grpcapp := grpcApp.NewGrpcApp(log, conf.Servers.GRPC, authService)

	return &App{
		GRPCServer: grpcapp,
	}

}
