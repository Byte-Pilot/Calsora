package config

import (
	"github.com/caarlos0/env/v10"
	"log"
)

type Config struct {
	PostgresURL string `env:"POSTGRES_DB_STR"`
	GinMode     string `env:"GIN_MODE" envDefault:"debug"`
	Port        string `env:"PORT" envDefault:"8080"`
}

func GetConfig() Config {
	envConfig := Config{}
	if err := env.Parse(&envConfig); err != nil {
		log.Panicf("apperrors parse env config: %s", err)
		return envConfig
	}

	return envConfig
}
