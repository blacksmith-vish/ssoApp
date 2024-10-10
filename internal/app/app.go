package app

import (
	"log/slog"
	authentication "sso/internal/api/authentication/rest"
	grpcApp "sso/internal/app/grpc"
	restApp "sso/internal/app/rest"
	"sso/internal/lib/config"
	"sso/internal/lib/migrate"
	authService "sso/internal/services/authentication"
	sqlstore "sso/internal/store/sql"
	"sso/internal/store/sql/sqlite"
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
	sqliteStore := sqlite.MustInitSqlite(conf.StorePath)
	migrate.MustMigrate(sqliteStore)

	store := sqlstore.NewStore(sqliteStore)

	// Инициализация auth сервиса
	authService := authService.NewService(
		log,
		conf.AuthenticationService,
		store.AuthenticationStore(),
		store.AuthenticationStore(),
		store.AuthenticationStore(),
	)

	grpcapp := grpcApp.NewGrpcApp(log, conf.GrpcConfig, authService)

	restapp := restApp.NewRestApp(log, conf.RestConfig, authentication.NewAuthenticationServer(log, authService))

	return &App{
		GRPCServer: grpcapp,
		RESTServer: restapp,
	}

}
