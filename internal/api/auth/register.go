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

func (srv *serverAPI) Register(
	ctx context.Context,
	request *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {

	if err := apiValidator.Validate(request); err != nil {
		return nil, err
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

	response := &ssov1.RegisterResponse{
		UserId: userID,
	}

	if err := apiValidator.Validate(response); err != nil {
		return nil, err
	}

	return response, nil
}
