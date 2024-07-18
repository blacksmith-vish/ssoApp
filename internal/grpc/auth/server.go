package auth

import (
	"context"

	ssov1 "github.com/blacksmith-vish/sso/protos/gen/go/sso"
	"google.golang.org/grpc"
)

var _ ssov1.AuthServer = (*serverAPI)(nil)

type serverAPI struct {
	ssov1.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	ssov1.RegisterAuthServer(gRPC, new(serverAPI))
}

func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	return &ssov1.LoginResponse{
		Token: req.GetEmail(),
	}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	panic("implement me")
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	panic("implement me")
}
