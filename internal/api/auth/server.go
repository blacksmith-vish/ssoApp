package auth

import (
	"context"

	ssov1 "github.com/blacksmith-vish/sso/protos/gen/go/sso"
	"google.golang.org/grpc"
)

var (
	_ ssov1.AuthServer = (*serverAPI)(nil)
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name=Auth --structname=AuthApi
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
