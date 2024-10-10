package authentication

import (
	"context"

	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

var (
	ErrUserExists = errors.New("user exists already")
)

// SaveUser saves user to db.
func (store *Store) SaveUser(ctx context.Context, nickname, email string, passHash []byte) (string, error) {
	const op = "Store.sqlite.SaveUser"

	stmt, err := store.db.Prepare("INSERT INTO users(id, nickname, email, pass_hash) VALUES(?, ?, ?, ?)")
	if err != nil {
		return "", errors.Wrap(err, op)
	}

	ID := uuid.New().String()

	_, err = stmt.ExecContext(ctx, ID, nickname, email, passHash)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return "", errors.Wrap(ErrUserExists, op)
		}

		return "", errors.Wrap(err, op)
	}

	return ID, nil
}
