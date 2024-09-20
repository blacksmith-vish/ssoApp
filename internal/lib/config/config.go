package config

import (
	"flag"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
)

const (
	EnvDev  = "dev"
	EnvProd = "prod"
	EnvTest = "test"
)

type Config struct {
	Env       string   `yaml:"env" validate:"oneof=dev prod"`
	StorePath string   `yaml:"store_path" validate:"required"`
	Services  Services `yaml:"services"`
	Servers   Servers  `yaml:"servers"`
}

func MustLoad() *Config {

	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	return MustLoadByPath(path)
}

func MustLoadByPath(path string) *Config {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}

	conf := new(Config)

	if err := cleanenv.ReadConfig(path, conf); err != nil {
		panic("failed to parse config file: " + err.Error())
	}

	mustValidate(
		conf,
		&conf.Servers,
	)

	return conf
}

// flag > env > default
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")

	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}

func (conf *Config) validate() error {
	if err := validator.New().Struct(conf); err != nil {
		return errors.Wrap(err, "config")
	}
	return nil
}
