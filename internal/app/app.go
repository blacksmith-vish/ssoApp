package app

import (
	"log/slog"
	grpcApp "sso/internal/app/grpc"
	"sso/internal/services/auth"
	"sso/internal/storage/sqlite"
	"time"
)

type App struct {
	GRPCServer *grpcApp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {

	// Инициализация хранилища
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	// Инициализация auth сервиса
	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcapp := grpcApp.New(log, authService, grpcPort)

	return &App{
		GRPCServer: grpcapp,
	}

}
