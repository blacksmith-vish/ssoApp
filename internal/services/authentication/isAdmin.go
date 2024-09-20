package authentication

import (
	"context"
	"log/slog"
	"sso/internal/services/authentication/models"

	"github.com/pkg/errors"
)

func (a *Authentication) IsAdmin(
	ctx context.Context,
	request models.IsAdminRequest,
) (models.IsAdminResponse, error) {
	const op = "auth.IsAdmin"

	log := a.log.With(
		slog.String("op", op),
		slog.String("userID", request.UserID),
	)

	log.Info("checking if user is admin")

	isAdmin, err := a.userProvider.IsAdmin(ctx, request.UserID)
	if err != nil {
		log.Error("error occured", slog.String("", err.Error()))
		return models.IsAdminResponse{
			IsAdmin: false,
		}, errors.Wrap(err, op)
	}

	log.Info("checked if user is admin", slog.Bool("is_admin", isAdmin))

	return models.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}
