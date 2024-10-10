package config

import (
	"sso/internal/lib/config"
	"time"
)

type Server struct {
	Port    uint16        `yaml:"port"`
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

func (srv Server) server() config.Server {
	return config.Server{
		Port:    srv.Port,
		Timeout: srv.Timeout,
	}
}

func (srvs Servers) getRestConfig() config.RESTConfig {
	return config.RESTConfig{
		Server: srvs.REST.server(),
	}
}

func (srvs Servers) getGrpcConfig() config.GRPCConfig {
	return config.GRPCConfig{
		Server: srvs.GRPC.server(),
	}
}
