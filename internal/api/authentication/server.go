package auth

import (
	"context"
	"sso/internal/domain"

	"github.com/blacksmith-vish/sso/gen/go/sso"
)

type Authentication interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		appID string,
	) (token string, err error)

	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (userID string, err error)

	IsAdmin(
		ctx context.Context,
		userID string,
	) (isAdmin bool, err error)
}

type authenticationServerAPI struct {
	sso.UnimplementedAuthenticationServer
	ctx  *domain.Context
	auth Authentication
}

func NewAuthenticationServer(
	ctx *domain.Context,
	auth Authentication,
) *authenticationServerAPI {

	return &authenticationServerAPI{
		ctx:  ctx,
		auth: auth,
	}

}

// func Register(gRPC *grpc.Server, auth Auth) {
// 	sso.RegisterAuthenticationServer(gRPC, &serverAPI{auth: auth})
// }
