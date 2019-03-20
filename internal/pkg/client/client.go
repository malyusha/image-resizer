package client

import (
	"errors"
	"fmt"

	"github.com/malyusha/image-resizer/internal/pkg/logger"
)

const (
	fsClient = "fs"
	httpClient = "http"
)

func New(config *Config) (Client, error) {
	return newClient(config)
}

func NewWithLogger(config *Config, logger logger.Logger) (Client, error) {
	c, err := newClient(config)

	if err != nil {
		return c, err
	}

	return WithLogger(c, logger), nil
}

func newClient(config *Config) (Client, error) {
	if config == nil {
		return new(NilClient), nil
	}

	switch config.Type {
	case fsClient:
		if config.FS == nil {
			return nil, errors.New("no configuration provided for filesystem client")
		}

		return NewFileSystemClient(config.FS)

	case httpClient:
		return NewHTTPClient(config.HTTP)
	}

	return nil, fmt.Errorf("client %s does not exist", config.Type)
}