package auth

import (
	"context"
	"errors"

	errs "sso/internal/domain/errors"

	ssov1 "github.com/blacksmith-vish/sso/protos/gen/go/sso"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	_        ssov1.AuthServer = (*serverAPI)(nil)
	validate *validator.Validate
)

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		appID int32,
	) (token string, err error)

	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (userID int64, err error)

	IsAdmin(
		ctx context.Context,
		userID int64,
	) (bool, error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

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

func (s *serverAPI) Register(
	ctx context.Context,
	request *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {

	if validate.Var(request.GetEmail(), "required,email") != nil {
		return nil, status.Error(codes.InvalidArgument, "email required")
	}

	if validate.Var(request.GetPassword(), "required") != nil {
		return nil, status.Error(codes.InvalidArgument, "password required")
	}

	userID, err := s.auth.RegisterNewUser(
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
