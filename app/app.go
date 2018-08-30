package app

import (
	"os"
	"fmt"
	"sync"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/malyusha/image-resizer/storage"
	"github.com/malyusha/image-resizer/console"
	"github.com/malyusha/image-resizer/client"
)

type Application interface {
	Router() *mux.Router
	Destroy()
	Storage() storage.Storage
	ImageClient() client.Client
}

// App is a singleton struct
type app struct {
	initialized bool
	router      *mux.Router
	server      *http.Server
	storage     storage.Storage
	imageClient client.Client
}

const (
	defaultImageClientName = "local"
	defaultStorageName     = "local"
)

var (
	once        sync.Once
	application *app
)

// Destroy destroys application
func (a *app) Destroy() {
	if !a.initialized {
		return
	}
}

// Router is the getter for mux.Router instance
func (a *app) Router() *mux.Router {
	return a.router
}

func (a *app) Storage() storage.Storage {
	return a.storage
}

func (a *app) ImageClient() client.Client {
	return a.imageClient
}

// GetInstance returns new instance of application
func GetInstance() *app {
	once.Do(func() {
		application = &app{
			router:      mux.NewRouter(),
			storage:     resolveStorage(),
			imageClient: resolveImageClient(),
		}
		application.init()
	})

	return application
}

// Initializes application
func (a *app) init() error {
	// Initialize router
	a.initialized = true
	a.routes()

	return nil
}

// resolveStorage resolves storage for application from arguments of CLI run
func resolveStorage() storage.Storage {
	storageName := defaultStorageName
	if console.Args.Storage != "" {
		storageName = console.Args.Storage
	}

	switch storageName {
	case "local":
		if console.Args.StorageDir == "" {
			fmt.Println("When using `local` type of storage you must provide storage " +
				"directory for storage")
			os.Exit(1)
		}
		s, err := storage.NewLocalStorage(&storage.LocalStorageConfig{
			Dir: console.Args.StorageDir,
		})
		if err != nil {
			panic(err)
		}

		return s
	}

	panic(fmt.Sprintf("No resolver for storage %s found", storageName))
}

// resolveImageClient resolves image client for application from arguments of CLI run
func resolveImageClient() client.Client {
	clientName := defaultImageClientName

	if console.Args.ImageClient != "" {
		clientName = console.Args.ImageClient
	}

	switch clientName {
	case "local":
		if console.Args.SourceDir == "" {
			fmt.Println("When using `local` type of images client you must provide source " +
				"directory for static files")
			os.Exit(1)
		}

		return client.NewLocalStorageClient(console.Args.SourceDir)
	}

	panic(fmt.Sprintf("No resolver for images client %s found", clientName))
}
