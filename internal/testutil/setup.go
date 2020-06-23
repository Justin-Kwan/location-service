package testutil

import (
	"log"

	"location-service/internal/types"
  "location-service/internal/config"
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
