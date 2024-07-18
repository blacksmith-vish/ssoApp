package app

import (
	"log/slog"
	grpcApp "sso/internal/app/grpc"
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

	// TODO: инициализация хранилища

	// TODO: инициализация auth сервиса

	grpcapp := grpcApp.New(log, grpcPort)

	return &App{
		GRPCServer: grpcapp,
	}

}
