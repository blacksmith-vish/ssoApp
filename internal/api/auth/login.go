package auth

import (
	"context"
	"errors"
	errs "sso/internal/domain/errors"
	apiValidator "sso/internal/lib/validators"

	ssov1 "github.com/blacksmith-vish/sso/protos/gen/go/sso"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serverAPI) Login(
	ctx context.Context,
	request *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {

	if err := apiValidator.Validate(request); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(
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

	response := &ssov1.LoginResponse{
		Token: token,
	}

	if err := apiValidator.Validate(response); err != nil {
		return nil, err
	}

	return response, nil
}
