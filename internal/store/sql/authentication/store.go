package authentication

import (
	"database/sql"
)

type AuthenticationStore struct {
	db *sql.DB
}

type Store = AuthenticationStore

func NewAuthenticationStore(db *sql.DB) *AuthenticationStore {
	return &AuthenticationStore{
		db: db,
	}
}
