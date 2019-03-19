package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Load(t *testing.T) {
	cfg := loadTestConfig()

	assert.Equal(t, "127.0.0.1", cfg.Server.HTTPAddr)
	assert.Equal(t, "127.0.0.1:8080", cfg.Server.Address())
}

func TestConfig_Get(t *testing.T) {
	cfg := loadTestConfig()
	expectedMap := map[string]interface{}{
		"addr": "https://example.com/store",
		"username": "test",
	}
	assert.Equal(t, expectedMap, cfg.Get("storage.external").Map().Raw())
	assert.Equal(t, "test", cfg.Get("storage.external.username").String())
}

func TestMustLoad(t *testing.T) {
	assert.NotPanics(t, func() {
		MustLoad("testdata/config.testing.yaml")
	}, "Loading config should not panic")
}

func loadTestConfig() *Config {
	return MustLoad("testdata/config.testing.yaml")
}