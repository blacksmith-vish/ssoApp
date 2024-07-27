package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string `yaml:"env" env-default:"local"`

	StoragePath string `yaml:"storage_path" env-required:"true"`

	TokenTTL time.Duration `yaml:"token_ttl" env-required:"true"`

	GRPC GRPCConfig `yaml:"grpc"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {

	path := fetchConfigPath()
	if len(path) == 0 {
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

	if len(res) == 0 {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
