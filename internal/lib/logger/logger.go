package logger

import (
	"io"
	"log/slog"
	"os"
	"sso/internal/lib/env"
	"sso/internal/lib/logger/handlers/dev"
)

type Config interface {
	GetEnv() string
}

func SetupLoggerWithConf(conf Config) *slog.Logger {
	return SetupLogger(conf.GetEnv())
}

func SetupLogger(Env string) *slog.Logger {

	devHandler := dev.NewHandler(
		os.Stdout,
		&slog.HandlerOptions{Level: slog.LevelDebug},
	)

	switch Env {

	case env.EnvDev:
		return slog.New(devHandler)

	case env.EnvProd:
		return slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)

	case env.EnvTest:
		return slog.New(
			slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)

	}

	return nil
}
