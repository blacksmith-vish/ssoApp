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

func (srv *server) Register(
	ctx context.Context,
	request *sso.RegisterRequest,
) (*sso.RegisterResponse, error) {

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

	response := &sso.RegisterResponse{
		UserId: userID,
	}

	if err := apiValidator.Validate(response); err != nil {
		return nil, err
	}

	return response, nil
}
