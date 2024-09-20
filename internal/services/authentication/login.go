package authentication

import (
	"context"
	"log/slog"
	"sso/internal/lib/jwt"
	"sso/internal/services/authentication/models"
	"sso/internal/store/sqlite"

	"github.com/pkg/errors"

	"golang.org/x/crypto/bcrypt"
)

// Login checks if user's credentials exists
func (a *Authentication) Login(
	ctx context.Context,
	request models.LoginRequest,
) (models.LoginResponse, error) {

	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", request.Email), // TODO email лучше не логировать
	)

	log.Info("attempting to login user")

	user, err := a.userProvider.User(ctx, request.Email)
	if err != nil {

		if errors.Is(err, sqlite.ErrUserNotFound) {
			log.Warn("user not found", slog.String("", err.Error()))

			return models.LoginResponse{}, errors.Wrap(ErrInvalidCredentials, op)
		}

		log.Error("failed to get user", slog.String("", err.Error()))
		return models.LoginResponse{}, errors.Wrap(err, op)
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(request.Password)); err != nil {
		log.Error("invalid credentials", slog.String("", err.Error()))
		return models.LoginResponse{}, errors.Wrap(ErrInvalidCredentials, op)
	}

	app, err := a.appProvider.App(ctx, request.AppID)
	if err != nil {

		if errors.Is(err, sqlite.ErrAppNotFound) {
			log.Warn("user not found", slog.String("", err.Error()))
			return models.LoginResponse{}, errors.Wrap(ErrInvalidAppID, op)
		}

		log.Error("failed getting app", slog.String("", err.Error()))
		return models.LoginResponse{}, errors.Wrap(err, op)
	}

	log.Info("user logged in succesfully")

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		log.Error("failed to get token", slog.String("", err.Error()))
		return models.LoginResponse{}, errors.Wrap(err, op)
	}

	return models.LoginResponse{
		Token: token,
	}, nil
}
