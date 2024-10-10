package authentication

import (
	"context"
	"log/slog"
	"sso/internal/services/authentication/models"
	auth_store "sso/internal/store/sql/authentication"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"

	"github.com/blacksmith-vish/sso/gen/go/sso"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srv *server) Register(
	ctx context.Context,
	request *sso.RegisterRequest,
) (*sso.RegisterResponse, error) {

	log := srv.log.With(
		slog.String("op", sso.Authentication_Register_FullMethodName),
	)

	serviceRequest := models.RegisterRequest{
		Email:    request.GetEmail(),
		Password: request.GetPassword(),
	}

	if err := validator.New().Struct(serviceRequest); err != nil {
		log.Error("validation failed", "err", err.Error())
		return nil, status.Error(codes.InvalidArgument, "login failed")
	}

	serviceResponse, err := srv.auth.RegisterNewUser(
		ctx,
		serviceRequest,
	)
	if err != nil {
		if errors.Is(err, auth_store.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "login failed")
		}
		return nil, status.Error(codes.Internal, "login failed")
	}

	response := &sso.RegisterResponse{
		UserId: serviceResponse.UserID,
	}

	return response, nil
}
