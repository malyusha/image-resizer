package app

import (
	"errors"
	"fmt"

	"github.com/ulule/gokvstores"
	"github.com/ulule/gostorages"

	"github.com/malyusha/image-resizer/internal/pkg/client"
	"github.com/malyusha/image-resizer/internal/pkg/config"
	"github.com/malyusha/image-resizer/internal/pkg/kvstore"
	"github.com/malyusha/image-resizer/internal/pkg/logger"
	"github.com/malyusha/image-resizer/internal/pkg/storage"
)

var (
	application          *app
	ErrAppNotInitialized = errors.New("Application hasn't been initialized yet")
)

// App is a singleton struct of Application interface
type app struct {
	config      *config.Config
	storage     gostorages.Storage
	kvStorage   gokvstores.KVStore
	client      client.Client
	logger      logger.Logger
}

func (a *app) KVStorage() gokvstores.KVStore {
	return a.kvStorage
}

// Config is a getter for app configuration instance
func (a *app) Config() *config.Config {
	return a.config
}

func (a *app) Logger() logger.Logger {
	return a.logger
}

// Storage returns storage for app
func (a *app) Storage() gostorages.Storage {
	return a.storage
}

// ImageClient returns image client of app
func (a *app) ImageClient() client.Client {
	return a.client
}

// Initializes application
func (a *app) init(config *config.Config) error {
	var (
		err       error
		kvStorage gokvstores.KVStore
		s         gostorages.Storage
		c         client.Client
		l         logger.Logger
	)

	l, err = logger.New(config.Log)
	kvStorage, err = kvstore.New(config.KVStorage)
	s, err = storage.New(config.Storage)
	c, err = client.New(config.Client)

	if config.Client != nil && config.Client.Log {
		c, err = client.NewWithLogger(config.Client, l)
	} else {
		c, err = client.New(config.Client)
	}

	if err != nil {
		return err
	}

	a.config = config
	a.logger = l
	a.kvStorage = kvStorage
	a.storage = s
	a.client = c

	return nil
}

// GetInstance returns application instance only if it's been initialized calling createInstance
// function. If no application initialized panic will be called.
func GetInstance() (Application, error) {
	if application == nil {
		return nil, ErrAppNotInitialized
	}

	return application, nil
}

// Creates an instance of application
func NewInstance(c *config.Config) (*app, error) {
	application = &app{}

	if err := application.init(c); err != nil {
		return nil, fmt.Errorf("Application init error: %s", err)
	}

	return application, nil
}
