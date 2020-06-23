package config

import (
	"flag"
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"location-service/internal/types"
)

// ParseEnvFlag reads environment command line. Return an error if an
// invalid flag is passed in. (flag must be "prod", "dev", or "test").
func ParseEnvFlag() (string, error) {
	var env string
	flag.StringVar(&env, "GO_ENV", "", "Runtime environment")
	flag.Parse()

	if !validEnvFlag(env) {
		return "", errors.Errorf("Invalid environment flag: %s", env)
	}
	return env, nil
}

func validEnvFlag(env string) bool {
	return env == "prod" || env == "dev" || env == "test"
}

// NewConfig returns a new configuration struct containing values
// needed for the websocket and rest server. The specific yaml config
// file parsed into the struct depends on the runtime environment flag
// passed in the command line. The file path is relative to the caller.
func NewConfig(filePath, env string) (*types.Config, error) {
	yamlData, err := loadFile(filePath, fileName(env))
	if err != nil {
		return nil, err
	}

	cfg := &types.Config{}

	cfg, err = unmarshalYAML(yamlData)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func fileName(env string) string {
	return "config." + env + ".yaml"
}

func loadFile(filePath, fileName string) ([]byte, error) {
	yamlData, err := ioutil.ReadFile(filePath + fileName)
	if err != nil {
		return nil, errors.Errorf("Error opening file '%s': %v", fileName, err)
	}
	return yamlData, nil
}

func unmarshalYAML(yamlData []byte) (*types.Config, error) {
	cfg := &types.Config{}
	if err := yaml.Unmarshal(yamlData, &cfg); err != nil {
		return nil, errors.Errorf("Error unmarshalling yaml: %v", err)
	}
	return cfg, nil
}
