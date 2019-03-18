package app

import (
	"errors"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/malyusha/image-resizer/pkg/client"
	"github.com/malyusha/image-resizer/pkg/config"
	"github.com/malyusha/image-resizer/pkg/storage"
)

var (
	once                 sync.Once
	application          *app
	ErrAppNotInitialized = errors.New("Application hasn't been initialized yet")
)

type Application interface {
	// Returns storage for image resizing
	Storage() storage.Storage
	// Returns instance of image client
	ImageClient() client.Client
	// Returns configuration for application
	Config() *config.Config
	// Returns logger for application
	Logger() *log.Logger
}

// App is a singleton struct of Application interface
type app struct {
	config      *config.Config
	storage     storage.Storage
	imageClient client.Client
	initialized bool
	connected   bool
	logger      *log.Logger
}

// Config is a getter for app configuration instance
func (a *app) Config() *config.Config {
	return a.config
}

func (a *app) Logger() *log.Logger {
	return a.logger
}

// GetInstance returns application instance only if it's been initialized calling createInstance
// function. If no application initialized panic will be called.
func GetInstance() (Application, error) {
	if application == nil {
		return nil, ErrAppNotInitialized
	}

	return application, nil
}

// IsProduction checks whether app running in production mode
func (a *app) IsProduction() bool {
	env := a.config.ENV

	return env == "prod" || env == "production"
}

// Creates an instance of application
func CreateInstance(c *config.Config) *app {
	once.Do(func() {
		application = &app{config: c}
		if err := application.config.Check(); err != nil {
			panic(fmt.Sprintf("ERROR CHECKING CONFIG: %s", err))
		}

		if err := application.init(); err != nil {
			panic(fmt.Sprintf("ERROR INITIALIZING APPLICATION: %s", err))
		}
	})

	return application
}

// Storage returns storage for app
func (a *app) Storage() storage.Storage {
	return a.storage
}

// ImageClient returns image client of app
func (a *app) ImageClient() client.Client {
	return a.imageClient
}

// Initializes application
func (a *app) init() error {
	a.createLogger()
	a.resolveImageClient()
	a.resolveStorage()

	return nil
}

// createLogger initializes new logger on app instance
func (a *app) createLogger() {
	a.logger = log.New()

	level, err := log.ParseLevel(a.config.LogLevel)
	if err != nil {
		// Cannot set level to given
		a.logger.Warnf("Failed to set log level to given `%s`", a.config.LogLevel)
		return
	}

	a.logger.SetLevel(level)
}

// resolveStorage resolves storage for application from arguments of CLI run
func (a *app) resolveStorage() {
	var (
		s           storage.Storage
		err         error
		storageName = a.config.Storage
	)

	switch storageName {
	case "local":
		c := &storage.LocalStorageConfig{
			Dir: a.config.Get("local_storage_dir").String(),
		}

		s, err = storage.NewLocalStorage(c)
	}

	if err != nil {
		a.logger.Fatalf("failed to resolve client: %s", err)
	}

	if s == nil {
		a.logger.Fatalf("No resolver for images client %s found", storageName)
	}

	a.logger.Infof(`Using "%s" storage`, storageName)
	a.storage = s
}

// resolveImageClient resolves image client for application from arguments of CLI run
func (a *app) resolveImageClient() {
	var (
		c          client.Client
		err        error
		clientName = a.config.ImageClient
	)
	switch clientName {
	case "local":
		directory := a.config.Get("local_client_dir").String()
		c, err = client.NewLocalStorageClient(directory)
	}

	if err != nil {
		a.logger.Fatalf("failed to resolve client: %s", err)
	}

	if c == nil {
		a.logger.Fatalf("No resolver for images client %s found", clientName)
	}

	a.logger.Infof(`Using "%s" image client`, clientName)
	a.imageClient = c
}