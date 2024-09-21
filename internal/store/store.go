package store

import (
	"database/sql"
	"embed"

	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type Storage interface {
	DB() *sql.DB
}

func Migrate(store Storage) error {

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("sqlite"); err != nil {
		return errors.Wrap(err, "could not set dialect")
	}

	if err := goose.Up(store.DB(), "migrations"); err != nil {
		return errors.Wrap(err, "could not run up")
	}
	return nil
}
