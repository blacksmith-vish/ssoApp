package suite

import (
	"context"
	"net"
	"sso/internal/lib/config"
	"strconv"
	"testing"

	config_yaml "sso/internal/store/filesystem/config/yaml"

	"github.com/blacksmith-vish/sso/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	grpcHost = "localhost"
)

type Suite struct {
	*testing.T
	Conf       *config.Config
	AuthClient sso.AuthenticationClient
}

func newConfig() *config.Config {
	yaml := config_yaml.MustLoadByPath("../../config/local.yaml")
	return config.NewConfig(yaml)
}

func New(t *testing.T) (context.Context, *Suite) {

	t.Helper()
	t.Parallel()

	conf := newConfig()

	ctx, cancelCtx := context.WithTimeout(context.Background(), conf.AuthenticationService.TokenTTL)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	cc, err := grpc.NewClient(
		grpcAddress(conf),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		Conf:       conf,
		AuthClient: sso.NewAuthenticationClient(cc),
	}

}

func grpcAddress(conf *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(int(conf.GrpcConfig.Port)))
}
