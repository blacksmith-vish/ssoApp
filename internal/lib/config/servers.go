package config

import (
	"time"

	"github.com/pkg/errors"
)

var (
	ErrServerPortsAreEqual = errors.New("REST & gRPC ports are equal")
)

type Servers struct {
	GRPC GRPCConfig `yaml:"grpc"`
	REST RESTConfig `yaml:"rest"`
}

type GRPCConfig struct {
	Port    uint16        `yaml:"port" validate:"gte=8000,lte=65535"`
	Timeout time.Duration `yaml:"timeout"`
}

type RESTConfig struct {
	Port    uint16        `yaml:"port" validate:"gte=8000,lte=65535"`
	Timeout time.Duration `yaml:"timeout"`
}

func (srv *Servers) validate() error {
	if srv.GRPC.Port == srv.REST.Port {
		return errors.Wrap(ErrServerPortsAreEqual, "servers")
	}
	return nil
}
