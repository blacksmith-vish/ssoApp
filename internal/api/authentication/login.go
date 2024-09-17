package auth

import (
	"context"
	"errors"
	errs "sso/internal/domain/errors"
	apiValidator "sso/internal/lib/validators"

	"github.com/blacksmith-vish/sso/gen/go/sso"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srv *authenticationServerAPI) Login(
	ctx context.Context,
	request *sso.LoginRequest,
) (*sso.LoginResponse, error) {

	if err := apiValidator.Validate(request); err != nil {
		return nil, err
	}

	token, err := srv.auth.Login(
		ctx,
		request.GetEmail(),
		request.GetPassword(),
		request.GetAppId(),
	)
	if err != nil {

		if errors.Is(err, errs.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "login failed")
		}

		return nil, status.Error(codes.Internal, "login failed")
	}

	response := &sso.LoginResponse{
		Token: token,
	}

	if err := apiValidator.Validate(response); err != nil {
		return nil, err
	}

	return response, nil
}
