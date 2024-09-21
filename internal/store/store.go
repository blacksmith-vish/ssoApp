package store

import (
	"database/sql"
	embed "sso"

	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
)

type Storage interface {
	DB() *sql.DB
}

func Migrate(store Storage) error {

	goose.SetBaseFS(embed.SQLiteMigrations)

	if err := goose.SetDialect("sqlite"); err != nil {
		return errors.Wrap(err, "could not set dialect")
	}

	if err := goose.Up(store.DB(), "migrations/sqlite"); err != nil {
		return errors.Wrap(err, "could not run up")
	}
	return nil
}
