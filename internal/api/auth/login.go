package auth

import (
	"context"
	"errors"
	errs "sso/internal/domain/errors"

	ssov1 "github.com/blacksmith-vish/sso/protos/gen/go/sso"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serverAPI) Login(
	ctx context.Context,
	request *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {

	if validate.Var(request.GetEmail(), "required,email") != nil {
		return nil, status.Error(codes.InvalidArgument, "email required")
	}

	if validate.Var(request.GetPassword(), "required") != nil {
		return nil, status.Error(codes.InvalidArgument, "password required")
	}

	if validate.Var(request.GetAppId(), "gte=0") != nil {
		return nil, status.Error(codes.InvalidArgument, "app_id required")
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

	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}
