package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Load(t *testing.T) {
	var cfg *Config
	assert.NotPanics(t, func() {
		cfg = MustLoad("testdata/config.testing.yaml")
	}, "Loading config should not panic")

	assert.Equal(t, "127.0.0.1", cfg.Server.HTTPAddr)
	assert.Equal(t, "127.0.0.1:8080", cfg.Server.Address())
}
