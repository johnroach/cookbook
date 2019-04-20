package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	assert.NoError(t, Init("test"))
}

func TestGetConfig(t *testing.T) {
	os.Setenv("COOKBOOK_TESTING_ENV", "testing_env_setup")
	Init("test")
	c := GetConfig()
	assert.Equal(t, c.Get("TESTING_ENV"), "testing_env_setup")
	assert.Equal(t, c.Get("TESTING_YAML"), "testing_yaml_setup")
}

func TestGetWithDefault(t *testing.T) {
	c := GetConfig()
	assert.Equal(t, GetWithDefault("NON_EXISTENT_CONF", "DEFAULT_VALUE"), c.Get("NON_EXISTENT_CONF"))
}