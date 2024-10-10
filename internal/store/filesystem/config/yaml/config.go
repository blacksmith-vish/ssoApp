package config

import (
	"flag"
	"os"
	"sso/internal/lib/config"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env       string   `yaml:"env"`
	StorePath string   `yaml:"store_path"`
	Services  Services `yaml:"services"`
	Servers   Servers  `yaml:"servers"`
}

func (conf Config) Convert() *config.Config {
	return &config.Config{
		Env:                   conf.Env,
		StorePath:             conf.StorePath,
		AuthenticationService: conf.Services.getAuthenticationService(),
		GrpcConfig:            conf.Servers.getGrpcConfig(),
		RestConfig:            conf.Servers.getRestConfig(),
	}
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
