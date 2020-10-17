package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type DB struct {
	Address string `yaml:"address" env:"DB_ADDRESS" env-default:"mongodb://localhost:27017" env-description:"database address"`
	Name    string `yaml:"database" env:"DB_NAME" env-default:"atlant" env-description:"database name"`
	Timeout int64  `yaml:"timeout" env:"DB_TIMEOUT" env-default:"10" env-description:"datebase timeout in sec."`
}

type GRPC struct {
	Port string `yaml:"port" env:"GRPS_PORT" env-default:"50051" env-description:"gRPC port"`
}

type Config struct {
	GRPC GRPC `yaml:"grpc"`
	DB   DB   `yaml:"database"`
}

func Get(name string) (cfg *Config, err error) {
	cfg = new(Config)
	if name != "" {
		err = cleanenv.ReadConfig(name, cfg)
	} else {
		err = cleanenv.ReadEnv(cfg)
	}

	return cfg, err
}
