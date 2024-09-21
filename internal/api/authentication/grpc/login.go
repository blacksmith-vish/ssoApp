package authentication

import (
	"context"
	"log/slog"
	"sso/internal/services/authentication"
	"sso/internal/services/authentication/models"

	"github.com/pkg/errors"

	"github.com/blacksmith-vish/sso/gen/go/sso"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srv *server) Login(
	ctx context.Context,
	request *sso.LoginRequest,
) (*sso.LoginResponse, error) {

	log := srv.log.With(
		slog.String("op", sso.Authentication_Login_FullMethodName),
	)

	serviceRequest := models.LoginRequest{
		Email:    request.GetEmail(),
		Password: request.GetPassword(),
		AppID:    request.GetAppId(),
	}

	if err := validator.New().Struct(serviceRequest); err != nil {
		log.Error("validation failed", "err", err.Error())
		return nil, status.Error(codes.InvalidArgument, "login failed")
	}

	serviceResponse, err := srv.auth.Login(
		ctx,
		serviceRequest,
	)
	if err != nil {

		if errors.Is(err, authentication.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "login failed")
		}

		return nil, status.Error(codes.Internal, "login failed")
	}

	response := &sso.LoginResponse{
		Token: serviceResponse.Token,
	}

	return response, nil
}
