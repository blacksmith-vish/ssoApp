package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"sso/internal/app"
	"sso/internal/lib/config"
	"sso/internal/lib/logger"
	config_yaml "sso/internal/store/filesystem/config/yaml"
)

func main() {

	// Инициализация конфига
	yaml := config_yaml.MustLoad()

	conf := config.NewConfig(yaml)

	log := logger.SetupLogger(conf.Env)

	log.Info("start app")

	// Инициализация приложения
	application := app.NewApp(log, conf)

	// Инициализация gRPC-сервер
	go application.GRPCServer.MustRun()

	go application.RESTServer.MustRun()

	// Graceful shut down

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sig := <-stop

	log.Info("app stopping", slog.String("signal", sig.String()))

	ctx, cancel := context.WithTimeout(context.Background(), conf.RestConfig.Timeout)
	defer func() {
		// extra handling here
		cancel()
	}()

	application.GRPCServer.Stop()

	application.RESTServer.Stop(ctx)

	log.Info("app stopped")
}
