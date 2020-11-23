package main

import (
	"log"

	"location-service/internal/config"
	"location-service/internal/wire"
)

const (
	ConfigFilePath = ""
)

func main() {
	env, err := config.ParseEnvFlag()
	if err != nil {
		log.Fatalf(err.Error())
	}

	cfg, err := config.NewConfig(ConfigFilePath, env)
	if err != nil {
		log.Fatalf(err.Error())
	}

	log.Printf("App config: %+v \n", *cfg)

	provider := wire.NewProvider(*cfg)

	// sh := provider.ProvideSocketHandler()
	gh := provider.ProvideGrpcHandler()

	// go sh.Serve()
	if err := gh.Serve(); err != nil {
		panic(err)
	}
}
