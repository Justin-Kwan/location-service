package testutil

import (
	"log"

	"location-service/internal/config"
	"location-service/internal/types"
)

const (
	ConfigFilePath = "../../../"
	Env            = "test"
)

func GetConfig() *types.Config {
	cfg, err := config.NewConfig(ConfigFilePath, Env)
	if err != nil {
		log.Fatalf(err.Error())
	}

	return cfg
}
