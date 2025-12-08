package config

import (
	"github.com/caarlos0/env/v10"
	"log"
)

type Config struct {
	PostgresURL string `env:"POSTGRES_DB_STR"`
}

func GetConfig() Config {
	envConfig := Config{}
	if err := env.Parse(&envConfig); err != nil {
		log.Panicf("Error parse env config: %s", err)
		return envConfig
	}

	return envConfig
}
