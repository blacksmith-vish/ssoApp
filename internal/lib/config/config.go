package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

const (
	EnvDev  = "dev"
	EnvProd = "prod"
	EnvTest = "test"
)

type Config struct {
	Env                   string `validate:"oneof=dev prod test"`
	StorePath             string `validate:"required"`
	AuthenticationService AuthenticationService
	GrpcConfig            GRPCConfig
	RestConfig            RESTConfig
}

type ConfigProvider interface {
	Convert() *Config
}

func NewConfig(
	config ConfigProvider,
) *Config {

	conf := config.Convert()

	if err := conf.validate(); err != nil {
		panic(err)
	}

	return conf
}

func (conf *Config) validate() error {

	const op = "config.validate"

	if err := validator.New().Struct(conf); err != nil {
		return errors.Wrap(err, op)
	}

	if err := comparePorts(conf.GrpcConfig.Server, conf.RestConfig.Server); err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}
