package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidEnvFlag(t *testing.T) {
	assert.Equal(t, true, validEnvFlag("prod"), "should be a valid environment flag")
	assert.Equal(t, true, validEnvFlag("dev"), "should be a valid environment flag")
	assert.Equal(t, true, validEnvFlag("test"), "should be a valid environment flag")
	assert.Equal(t, false, validEnvFlag("prod "), "should be a invalid environment flag")
	assert.Equal(t, false, validEnvFlag(" test"), "should be a invalid environment flag")
	assert.Equal(t, false, validEnvFlag(" dev"), "should be a invalid environment flag")
	assert.Equal(t, false, validEnvFlag("invalid"), "should be a invalid environment flag")
}

func TestFileName(t *testing.T) {
	assert.Equal(t, "config.prod.yaml", fileName("prod"), "should return production configuration file name")
	assert.Equal(t, "config.dev.yaml", fileName("dev"), "should return development configuration file name")
	assert.Equal(t, "config.test.yaml", fileName("test"), "should return test configuration file name")
	assert.Equal(t, "config. .yaml", fileName(" "), "should return configuration file with ' ' in name")
	assert.Equal(t, "config.*.yaml", fileName("*"), "should return configuration file with '*' in name")
}
