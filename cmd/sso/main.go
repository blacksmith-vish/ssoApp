package main

import (
	"log/slog"
	"os"
	"os/signal"
	"sso/internal/config"
	"syscall"

	"sso/internal/app"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	// инициализация конфига
	conf := config.MustLoad()

	// инициализация логирования
	log := setupLogger(conf.Env)

	log.Info("start app")

	// TODO: инициализация приложения (app)
	application := app.New(
		log,
		conf.GRPC.Port,
		conf.StoragePath,
		conf.TokenTTL,
	)

	go application.GRPCServer.MustRun()

	// TODO: инициализация gRPC-сервер

	// TODO: Graceful shut down

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sig := <-stop

	log.Info("application stopping", slog.String("signal", sig.String()))

	application.GRPCServer.Stop()

	log.Info("application stopped")

}

func setupLogger(env string) *slog.Logger {

	var log *slog.Logger

	switch env {

	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)

	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)

	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
