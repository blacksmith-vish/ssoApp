package authentication

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

func (store *Store) IsAdmin(ctx context.Context, userID string) (bool, error) {
	const op = "Store.sqlite.IsAdmin"

	stmt, err := store.db.Prepare("SELECT is_admin FROM users WHERE id = ?")
	if err != nil {
		return false, errors.Wrap(err, op)
	}

	row := stmt.QueryRowContext(ctx, userID)

	var isAdmin bool

	err = row.Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, errors.Wrap(ErrUserNotFound, op)
		}

		return false, errors.Wrap(err, op)
	}

	return isAdmin, nil
}
