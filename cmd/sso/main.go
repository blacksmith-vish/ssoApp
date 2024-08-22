package main

import (
	"log/slog"
	"os"
	"os/signal"
	"sso/internal/domain"
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
