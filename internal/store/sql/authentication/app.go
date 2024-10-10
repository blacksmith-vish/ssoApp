package authentication

import (
	"context"
	"database/sql"
	"sso/internal/store/models"

	"github.com/pkg/errors"
)

var (
	ErrAppNotFound = errors.New("app not found")
)

// App returns app by id.
func (store *Store) App(ctx context.Context, id string) (models.App, error) {
	const op = "Store.sqlite.App"

	stmt, err := store.db.Prepare("SELECT id, name, secret FROM apps WHERE id = ?")
	if err != nil {
		return models.App{}, errors.Wrap(err, op)
	}

	row := stmt.QueryRowContext(ctx, id)

	var app models.App
	err = row.Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, errors.Wrap(ErrAppNotFound, op)
		}

		return models.App{}, errors.Wrap(err, op)
	}

	return app, nil
}
