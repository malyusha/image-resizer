package app

import (
	"github.com/ulule/gokvstores"
	"github.com/ulule/gostorages"

	"github.com/malyusha/image-resizer/internal/pkg/client"
	"github.com/malyusha/image-resizer/internal/pkg/config"
	"github.com/malyusha/image-resizer/internal/pkg/logger"
)

type Application interface {
	// Returns storage for image resizing
	Storage() gostorages.Storage
	// KVStorage returns KV storage
	KVStorage() gokvstores.KVStore
	// Returns instance of image client
	ImageClient() client.Client
	// Returns configuration for application
	Config() *config.Config
	// Returns logger for application
	Logger() logger.Logger
}
