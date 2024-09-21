package authentication

import (
	"context"
	"log/slog"
	"sso/internal/services/authentication/models"
	"sso/internal/store/sqlite"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"

	"github.com/blacksmith-vish/sso/gen/go/sso"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srv *server) IsAdmin(
	ctx context.Context,
	request *sso.IsAdminRequest,
) (*sso.IsAdminResponse, error) {

	log := srv.log.With(
		slog.String("op", sso.Authentication_IsAdmin_FullMethodName),
		slog.String("userID", request.GetUserId()),
	)

	serviceRequest := models.IsAdminRequest{
		UserID: request.GetUserId(),
	}

	if err := validator.New().Struct(serviceRequest); err != nil {
		log.Error("validation failed", "err", err.Error())
		return nil, status.Error(codes.InvalidArgument, "login failed")
	}

	serviceResponse, err := srv.auth.IsAdmin(
		ctx,
		serviceRequest,
	)

	if err != nil {
		if errors.Is(err, sqlite.ErrUserNotFound) {
			return nil, status.Error(codes.AlreadyExists, "login failed")
		}
		return nil, status.Error(codes.Internal, "login failed")
	}

	response := &sso.IsAdminResponse{
		IsAdmin: serviceResponse.IsAdmin,
	}

	return response, nil
}
