package migrate

import (
	"database/sql"
	embed "sso"

	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
)

type Storage interface {
	DB() *sql.DB
	Dialect() string
	MigrationsPath() string
}

func MustMigrate(store Storage) {
	if err := Migrate(store); err != nil {
		panic(err)
	}
}

func Migrate(store Storage) error {

	goose.SetBaseFS(embed.SQLiteMigrations)

	dialect := store.Dialect()

	if err := goose.SetDialect(dialect); err != nil {
		return errors.Wrap(err, "could not set dialect "+dialect)
	}

	if err := goose.Up(store.DB(), store.MigrationsPath()); err != nil {
		return errors.Wrap(err, "could not run migrations up")
	}
	return nil
}
