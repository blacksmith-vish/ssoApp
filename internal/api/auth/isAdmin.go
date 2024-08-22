package auth

import (
	"context"
	errs "sso/internal/domain/errors"

	"github.com/pkg/errors"

	apiValidator "sso/internal/lib/validators"

	ssov1 "github.com/blacksmith-vish/sso/protos/gen/go/sso"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serverAPI) IsAdmin(
	ctx context.Context,
	request *ssov1.IsAdminRequest,
) (*ssov1.IsAdminResponse, error) {

	if err := apiValidator.Validate(request); err != nil {
		return nil, err
	}

	isAdmin, err := s.auth.IsAdmin(
		ctx,
		request.GetUserId(),
	)

	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return nil, status.Error(codes.AlreadyExists, "login failed")
		}
		return nil, status.Error(codes.Internal, "login failed")
	}

	response := &ssov1.IsAdminResponse{
		IsAdmin: isAdmin,
	}

	if err := apiValidator.Validate(response); err != nil {
		return nil, err
	}

	return response, nil
}
