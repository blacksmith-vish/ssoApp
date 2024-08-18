package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	def "sso/internal/api/auth"
	"sso/internal/domain"
	errs "sso/internal/domain/errors"
	"sso/internal/lib/jwt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var _ def.Auth = (*Auth)(nil)

type Auth struct {
	ctx      *domain.Context
	store    AuthStoreProvider
	tokenTTL time.Duration
}

// New returns a new instance of Auth
func New(
	ctx *domain.Context,
	storeProvider AuthStoreProvider,
) *Auth {
	return &Auth{
		ctx:      ctx,
		store:    storeProvider,
		tokenTTL: ctx.Config().TokenTTL,
	}
}

// Login checks if user's credentials exists
func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string,
	appID int32,
) (string, error) {

	const op = "auth.Login"

	log := a.ctx.Log().With(
		slog.String("op", op),
		slog.String("email", email), // TODO email лучше не логировать
	)

	log.Info("attempting to login user")

	user, err := a.store.userProvider.User(ctx, email)
	if err != nil {

		if errors.Is(err, errs.ErrUserNotFound) {
			log.Warn("user not found", slog.String("", err.Error()))
			return "", fmt.Errorf("%s: %w", op, errs.ErrInvalidCredentials)
		}

		log.Error("failed to get user", slog.String("", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		log.Error("invalid credentials", slog.String("", err.Error()))
		return "", fmt.Errorf("%s: %w", op, errs.ErrInvalidCredentials)
	}

	app, err := a.store.appProvider.App(ctx, appID)
	if err != nil {

		if errors.Is(err, errs.ErrAppNotFound) {
			log.Warn("user not found", slog.String("", err.Error()))
			return "", fmt.Errorf("%s: %w", op, errs.ErrInvalidAppID)
		}

		log.Error("failed getting app", slog.String("", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logged in succesfully")

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		log.Error("failed to get token", slog.String("", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (a *Auth) RegisterNewUser(
	ctx context.Context,
	email string,
	password string,
) (userID int64, err error) {

	const op = "auth.RegisterNewUser"

	log := a.ctx.Log().With(
		slog.String("op", op),
		slog.String("email", email), // TODO email лучше не логировать
	)

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate pass hash", slog.String("", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.store.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {

		if errors.Is(err, errs.ErrUserExists) {
			log.Warn("user exists", slog.String("", err.Error()))
			return 0, fmt.Errorf("%s: %w", op, errs.ErrUserExists)
		}

		log.Error("failed to save user", slog.String("", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user registered")

	return id, nil

}

func (a *Auth) IsAdmin(
	ctx context.Context,
	userID int64,
) (bool, error) {
	const op = "auth.IsAdmin"

	log := a.ctx.Log().With(
		slog.String("op", op),
		slog.Int64("userID", userID),
	)

	log.Info("checking if user is admin")

	isAdmin, err := a.store.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		log.Error("error occured", slog.String("", err.Error()))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("checked if user is admin", slog.Bool("is_admin", isAdmin))

	return isAdmin, nil
}
