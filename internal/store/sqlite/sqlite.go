package sqlite

import (
	"context"
	"database/sql"
	errs "sso/internal/domain/errors"
	def "sso/internal/services/auth"
	"sso/internal/store/models"

	"github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

var (
	_ def.UserSaver    = (*Store)(nil)
	_ def.UserProvider = (*Store)(nil)
	_ def.AppProvider  = (*Store)(nil)
)

type Store struct {
	db *sql.DB
}

func New(StorePath string) (*Store, error) {

	const op = "Store.sqlite.New"

	db, err := sql.Open("sqlite3", StorePath)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return &Store{db: db}, nil
}

func (s *Store) Stop() error {	
	return s.db.Close()
}

// SaveUser saves user to db.
func (s *Store) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "Store.sqlite.SaveUser"

	stmt, err := s.db.Prepare("INSERT INTO users(email, pass_hash) VALUES(?, ?)")
	if err != nil {
		return 0, errors.Wrap(err, op)
	}

	res, err := stmt.ExecContext(ctx, email, passHash)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, errors.Wrap(errs.ErrUserExists, op)
		}

		return 0, errors.Wrap(err, op)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, errors.Wrap(err, op)
	}

	return id, nil
}

// User returns user by email.
func (s *Store) User(ctx context.Context, email string) (models.User, error) {
	const op = "Store.sqlite.User"

	stmt, err := s.db.Prepare("SELECT id, email, pass_hash FROM users WHERE email = ?")
	if err != nil {
		return models.User{}, errors.Wrap(err, op)
	}

	row := stmt.QueryRowContext(ctx, email)

	var user models.User
	err = row.Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, errors.Wrap(errs.ErrUserNotFound, op)
		}

		return models.User{}, errors.Wrap(err, op)
	}

	return user, nil
}

//func (s *Store) SavePermission(ctx context.Context, userID int64, permission models.Permission, appID string) error {
//	const op = "Store.sqlite.SavePermission"
//
//	stmt, err := s.db.Prepare("INSERT INTO permissions(user_id, permission, app_id) VALUES(?, ?, ?)")
//	if err != nil {
//		return errors.Wrap(err, op)
//	}
//
//	_, err = stmt.ExecContext(ctx, userID, permission, appID)
//	if err != nil {
//		return errors.Wrap(err, op)
//	}
//
//	return nil
//}

// App returns app by id.
func (s *Store) App(ctx context.Context, id int32) (models.App, error) {
	const op = "Store.sqlite.App"

	stmt, err := s.db.Prepare("SELECT id, name, secret FROM apps WHERE id = ?")
	if err != nil {
		return models.App{}, errors.Wrap(err, op)
	}

	row := stmt.QueryRowContext(ctx, id)

	var app models.App
	err = row.Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, errors.Wrap(errs.ErrAppNotFound, op)
		}

		return models.App{}, errors.Wrap(err, op)
	}

	return app, nil
}

func (s *Store) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "Store.sqlite.IsAdmin"

	stmt, err := s.db.Prepare("SELECT is_admin FROM users WHERE id = ?")
	if err != nil {
		return false, errors.Wrap(err, op)
	}

	row := stmt.QueryRowContext(ctx, userID)

	var isAdmin bool

	err = row.Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, errors.Wrap(errs.ErrUserNotFound, op)
		}

		return false, errors.Wrap(err, op)
	}

	return isAdmin, nil
}
