package app

import (
	"log/slog"
	authentication "sso/internal/api/authentication/rest"
	grpcApp "sso/internal/app/grpc"
	restApp "sso/internal/app/rest"
	"sso/internal/lib/config"
	"sso/internal/lib/migrate"
	authService "sso/internal/services/authentication"
	"sso/internal/store/sqlite"
)

type App struct {
	GRPCServer *grpcApp.App
	RESTServer *restApp.App
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

	if err := migrate.Migrate(storage); err != nil {
		panic(err)
	}

	// Инициализация auth сервиса
	authService := authService.NewService(
		log,
		storage,
		storage,
		storage,
		conf.Services.Authentication.TokenTTL,
	)

	grpcapp := grpcApp.NewGrpcApp(log, conf.Servers.GRPC, authService)

	restapp := restApp.NewRestApp(log, conf.Servers.REST, authentication.NewAuthenticationServer(log, authService))

	return &App{
		GRPCServer: grpcapp,
		RESTServer: restapp,
	}

}
