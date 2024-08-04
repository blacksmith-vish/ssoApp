package auth

import (
	"context"
	"errors"
	errs "sso/internal/domain/errors"

	ssov1 "github.com/blacksmith-vish/sso/protos/gen/go/sso"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serverAPI) IsAdmin(
	ctx context.Context,
	request *ssov1.IsAdminRequest,
) (*ssov1.IsAdminResponse, error) {

	if validate.Var(request.GetUserId(), "gte=0") != nil {
		return nil, status.Error(codes.InvalidArgument, "app_id required")
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

	return &ssov1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}
