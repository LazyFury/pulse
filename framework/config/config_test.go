package config_test

import (
	"testing"

	"github.com/lazyfury/pulse/framework/config"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	assert := assert.New(t)

	conf := config.NewConfig()
	assert.NotNil(conf)
}

func TestConfigFromFile(t *testing.T) {
	assert := assert.New(t)

	conf := config.NewConfigFromFile()
	assert.NotNil(conf)
}

func TestConfigFromBytes(t *testing.T) {
	assert := assert.New(t)
	configBytes := []byte(`
port: 8900
`)
	conf := config.NewConfigFromBytes(configBytes)
	assert.NotNil(conf)
}
