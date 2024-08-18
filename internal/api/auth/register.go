package auth

import (
	"context"
	"errors"
	errs "sso/internal/domain/errors"

	ssov1 "github.com/blacksmith-vish/sso/protos/gen/go/sso"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srv *serverAPI) Register(
	ctx context.Context,
	request *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {

	if validate.Var(request.GetEmail(), "required,email") != nil {
		return nil, status.Error(codes.InvalidArgument, "email required")
	}

	if validate.Var(request.GetPassword(), "required") != nil {
		return nil, status.Error(codes.InvalidArgument, "password required")
	}

	userID, err := srv.auth.RegisterNewUser(
		ctx,
		request.GetEmail(),
		request.GetPassword(),
	)
	if err != nil {
		if errors.Is(err, errs.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "login failed")
		}
		return nil, status.Error(codes.Internal, "login failed")
	}

	return &ssov1.RegisterResponse{
		UserId: userID,
	}, nil
}
