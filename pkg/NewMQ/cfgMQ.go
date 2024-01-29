package cfgMQ

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type ConfigMQ struct {
	URI string
}

func NewCfgMQ() (*ConfigMQ, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err

	}
	var cfgMQ ConfigMQ
	if err := envconfig.Process("mq", &cfgMQ); err != nil {
		return nil, err
	}
	return &cfgMQ, nil
}
