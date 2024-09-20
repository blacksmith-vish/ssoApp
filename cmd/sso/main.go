package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"sso/internal/app"
	"sso/internal/lib/config"
	"sso/internal/lib/logger"
)

func main() {

	// Инициализация конфига
	conf := config.MustLoad()

	log := logger.SetupLogger(conf.Env)

	log.Info("start app")

	// Инициализация приложения
	application := app.NewApp(log, conf)

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
