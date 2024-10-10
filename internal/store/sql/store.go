package sqlstore

import (
	"database/sql"
	"sso/internal/store/sql/authentication"
)

type Store struct {
	db                  *sql.DB
	authenticationStore *authentication.AuthenticationStore
}

type StoreProvider interface {
	DB() *sql.DB
}

func NewStore(
	store StoreProvider,
) *Store {

	db := store.DB()

	return &Store{
		db:                  db,
		authenticationStore: authentication.NewAuthenticationStore(db),
	}
}

func (store *Store) Stop() error {
	return store.db.Close()
}

func (store *Store) AuthenticationStore() *authentication.AuthenticationStore {
	return store.authenticationStore
}
