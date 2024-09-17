package auth

import (
	"context"
	"log/slog"
	errs "sso/internal/domain/errors"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"

	"github.com/blacksmith-vish/sso/gen/go/sso"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srv *server) IsAdmin(
	ctx context.Context,
	request *sso.IsAdminRequest,
) (*sso.IsAdminResponse, error) {

	log := srv.ctx.Log().With(
		slog.String("op", sso.Authentication_IsAdmin_FullMethodName),
		slog.String("userID", request.GetUserId()),
	)

	validate := validator.New()

	err := validate.Var(request.GetUserId(), "required,uuid4")
	if err != nil {
		log.Error("validation failed", "err", err.Error())
		return nil, status.Error(codes.InvalidArgument, "login failed")
	}

	log.Debug("validation failed", "err",
		[]string{
			"shit1",
			"shit2",
			"shit3",
		},
	)

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
