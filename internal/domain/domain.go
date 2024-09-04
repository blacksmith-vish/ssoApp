package domain

import (
	"io"
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
		log:  setupLogger(conf.Env),
		conf: conf,
	}
}

// NewContext Возвращает новую структуру Context
// Принимает функции для задания опций в любом количестве
func NewContextWithOpts(
	opts ...optsFunc,
) *Context {
	o := defaultOpts()
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func NewContext(
	Log *slog.Logger,
	Config *config.Config,
) *Context {

	return &Context{
		log:  Log,
		conf: Config,
	}

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

	case config.EnvTest:
		return slog.New(
			slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)

	}

	return nil
}
