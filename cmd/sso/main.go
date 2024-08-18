package main

import (
	"log/slog"
	"os"
	"os/signal"
	"sso/internal/domain"
	"sso/internal/lib/config"
	"sso/internal/lib/log/handlers/dev"
	"syscall"

	"sso/internal/app"
)

func main() {

	// Инициализация контекста приложения
	ctx := domain.NewContext()

	log := ctx.Log()

	log.Info("start app")

	// Инициализация приложения
	application := app.New(ctx)

	// Инициализация gRPC-сервер
	go application.GRPCServer.MustRun()

	// Graceful shut down

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sig := <-stop

	log.Info("application stopping", slog.String("signal", sig.String()))

	application.GRPCServer.Stop()

	log.Info("application stopped")
}

func setupLogger(env string) *slog.Logger {

	switch env {

	case config.EnvDev:
		return slog.New(
			dev.NewHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelDebug},
			),
		)

	case config.EnvProd:
		return slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return nil
}
