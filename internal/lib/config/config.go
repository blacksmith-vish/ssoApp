package config

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

const (
	EnvDev  = "dev"
	EnvProd = "prod"
	EnvTest = "test"
)

type Config struct {
	Env string `yaml:"env" env-default:"prod" validate:"oneof=dev prod"`

	StorePath string `yaml:"store_path" env-required:"true"`

	TokenTTL time.Duration `yaml:"token_ttl" env-required:"true"`

	GRPC GRPCConfig `yaml:"grpc"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port" validate:"gte=1000,lte=99999"`
	Timeout time.Duration `yaml:"timeout"`
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

	if err := validator.New().Struct(conf); err != nil {
		panic("failed to validate config: " + err.Error())
	}

	return conf
}

// flag > env > default
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")

	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file")
	}

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
