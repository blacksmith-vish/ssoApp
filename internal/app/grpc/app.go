package grpcApp

import (
	"fmt"
	"log/slog"
	"net"
	authenticationGRPC "sso/internal/api/authentication"
	"sso/internal/domain"

	"github.com/blacksmith-vish/sso/gen/go/sso"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type App struct {
	ctx        *domain.Context
	gRPCServer *grpc.Server
	port       int
}

func New(
	ctx *domain.Context,
	authService authenticationGRPC.Authentication,
) *App {

	gRPCServer := grpc.NewServer()

	sso.RegisterAuthenticationServer(
		gRPCServer,
		authenticationGRPC.NewAuthenticationServer(
			ctx,
			authService,
		),
	)

	//	authenticationGRPC.Register(gRPCServer, authService)

	return &App{
		ctx:        ctx,
		gRPCServer: gRPCServer,
		port:       ctx.Config().GRPC.Port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {

	const op = "grpcApp.Run"

	log := a.ctx.Log().With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	log.Info("starting gRPC server")

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return errors.Wrap(err, op)
	}

	log.Info("gRPC server is running", slog.String("addr", listener.Addr().String()))

	if err := a.gRPCServer.Serve(listener); err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}

func (a *App) Stop() {

	const op = "grpcApp.Stop"

	a.ctx.Log().With(slog.String("op", op)).
		Info("stopping gRPC server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()

}
