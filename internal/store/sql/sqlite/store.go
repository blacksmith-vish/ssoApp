package sqlite

import (
	"database/sql"
	"path"

	"github.com/pkg/errors"
)

const (
	dialect string = "sqlite"
)

type store struct {
	db *sql.DB
}

func MustInitSqlite(StorePath string) *store {
	store, err := NewSqliteStore(StorePath)
	if err != nil {
		panic(err)
	}
	return store
}

func NewSqliteStore(StorePath string) (*store, error) {

	const op = "Store.sqlite.New"

	db, err := sql.Open("sqlite3", StorePath)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return &store{
		db: db,
	}, nil
}

func (store *store) DB() *sql.DB {
	return store.db
}

func (store *store) Dialect() string {
	return dialect
}

func (store *store) MigrationsPath() string {
	return path.Join("migrations", dialect)
}

func (store *store) Stop() error {
	return store.db.Close()
}
