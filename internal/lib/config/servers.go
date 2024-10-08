package config

import (
	"time"

	"github.com/pkg/errors"
)

var (
	ErrServerPortsAreEqual = errors.New("REST & gRPC ports are equal")
)

type Server struct {
	Port    uint16        `yaml:"port" validate:"gte=8000,lte=65535"`
	Timeout time.Duration `yaml:"timeout"`
}

type GRPCConfig struct {
	Server `yaml:"server"`
}

type RESTConfig struct {
	Server `yaml:"server"`
}

type Servers struct {
	GRPC GRPCConfig `yaml:"grpc"`
	REST RESTConfig `yaml:"rest"`
}

func (srvs *Servers) validate() error {
	if srvs.GRPC.Server.Port == srvs.REST.Server.Port {
		return errors.Wrap(ErrServerPortsAreEqual, "servers")
	}
	return nil
}

func (srv Server) GetTimeout() time.Duration {
	return srv.Timeout
}

func (srv Server) GetPort() uint16 {
	return srv.Port
}
