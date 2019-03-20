package config

import (
	"fmt"
	"time"

	"github.com/malyusha/image-resizer/internal/pkg/client"
	"github.com/malyusha/image-resizer/internal/pkg/kvstore"
	"github.com/malyusha/image-resizer/internal/pkg/logger"
	"github.com/malyusha/image-resizer/internal/pkg/storage"
)

type Server struct {
	// Application HTTP server address
	HTTPAddr string `mapstructure:"address"`
	// Application HTTP server port
	HTTPPort string `mapstructure:"port"`
	// Duration in seconds for which the server will wait existing connections to finish
	GracefulTimeout int `mapstructure:"graceful_timeout"`
}

// Config is the main structure for application configuring
type Config struct {
	// Log configuration
	Log *logger.Config
	// key/value storage configuration
	KVStorage *kvstore.Config
	// Storage type
	Storage *storage.Config
	// Image client configuration
	Client *client.Config
	// Server configuration
	Server *Server
}

// AddressString returns HTTP address with port
func (s Server) Address() string {
	return fmt.Sprintf("%s:%s", s.HTTPAddr, s.HTTPPort)
}

// GetGracefulTimeout returns number of seconds to wait until server shuts down
func (s *Server) GetGracefulTimeout() time.Duration {
	timeout := time.Duration(5)
	if s.GracefulTimeout != 0 {
		timeout = time.Duration(s.GracefulTimeout)
	}

	return time.Second * timeout
}

// Load returns new config struct from config file path
func Load(path string) (*Config, error) {
	return load(path)
}

func load(path string) (*Config, error) {
	return nil, nil
}
