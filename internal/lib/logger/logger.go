package logger

import (
	"io"
	"log/slog"
	"os"
	"sso/internal/lib/config"
	"sso/internal/lib/logger/handlers/dev"
)

func SetupLogger(env string) *slog.Logger {

	switch env {

	case config.EnvDev:
		return slog.New(
			dev.NewHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)

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
