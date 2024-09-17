package auth

import (
	"context"
	errs "sso/internal/domain/errors"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"

	"github.com/blacksmith-vish/sso/gen/go/sso"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srv *authenticationServerAPI) IsAdmin(
	ctx context.Context,
	request *sso.IsAdminRequest,
) (*sso.IsAdminResponse, error) {

	validate := validator.New()

	err := validate.Var(request.GetUserId(), "gte=0")
	if err != nil {
		return nil, err
	}

	isAdmin, err := srv.auth.IsAdmin(
		ctx,
		request.GetUserId(),
	)

	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return nil, status.Error(codes.AlreadyExists, "login failed")
		}
		return nil, status.Error(codes.Internal, "login failed")
	}

	response := &sso.IsAdminResponse{
		IsAdmin: isAdmin,
	}

	return response, nil
}
