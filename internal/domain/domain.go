package domain

import (
	"log/slog"
	"os"
	"sso/internal/lib/config"
	"sso/internal/lib/log/handlers/dev"
)

type Context struct {
	log  *slog.Logger
	conf *config.Config
}

// Сигнатура функции для задания параметров
type optsFunc func(*Context)

// Задает опции по умолчанию
func defaultOpts() *Context {
	conf := config.MustLoad()
	return &Context{
		conf: conf,
		log:  setupLogger(conf.Env),
	}
}

// NewContext Возвращает новую структуру Context
// Принимает функции для задания опций в любом количестве
func NewContext(
	opts ...optsFunc,
) *Context {
	o := defaultOpts()
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func (ctx *Context) Log() *slog.Logger {
	return ctx.log
}

func (ctx *Context) Config() *config.Config {
	return ctx.conf
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
