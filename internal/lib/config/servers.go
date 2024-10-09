package config

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

var (
	ErrServerPortAlreadyInUse = errors.New("port is already in use")
)

type Server struct {
	Port    uint16 `validate:"gte=8000,lte=65535"`
	Timeout time.Duration
}

type GRPCConfig struct {
	Server
}

type RESTConfig struct {
	Server
}

func comparePorts(servers ...Server) error {
	ports := make(map[uint16]*struct{})
	for _, server := range servers {
		if _, ok := ports[server.Port]; ok {
			return fmt.Errorf("%v: %w", server.Port, ErrServerPortAlreadyInUse)
		} else {
			ports[server.Port] = nil
		}
	}
	return nil
}
