package logger

import (
	"io"
	"log/slog"
	"os"
	"sso/internal/lib/config"
	"sso/internal/lib/logger/handlers/dev"
)

type Config interface {
	GetEnv() string
}

func SetupLoggerWithConf(conf Config) *slog.Logger {
	return SetupLogger(conf.GetEnv())
}

func SetupLogger(env string) *slog.Logger {

	devHandler := dev.NewHandler(
		os.Stdout,
		&slog.HandlerOptions{Level: slog.LevelDebug},
	)

	switch env {

	case config.EnvDev:
		return slog.New(devHandler)

	case config.EnvProd:
		return slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)

	case config.EnvTest:
		return slog.New(
			slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)

	}

	return nil
}
