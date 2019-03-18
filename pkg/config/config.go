package config

import (
	"fmt"
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

const (
	defaultAddr           = "0.0.0.0"
	defaultPort           = "8080"
	defaultResizeStrategy = "presets"
	defaultImageClient    = "local"
	defaultStorage        = "local"
	defaultENV            = "dev"
)

type Server struct {
	// Application HTTP server address
	HTTPAddr string `yaml:"address"`
	// Application HTTP server port
	HTTPPort string `yaml:"port"`
	// Duration in seconds for which the server will wait existing connections to finish
	GracefulTimeout int `yaml:"graceful_timeout"`
}

// Config is the main structure for application configuring
type Config struct {
	// ENV type
	ENV string `yaml:"env"`
	// Log level
	LogLevel string `yaml:"log_level"`
	// Resizing strategy for images. Available strategies can be found in `strategy` package
	ResizeStrategy string `yaml:"resize_strategy"`
	// Storage type
	Storage string `yaml:"storage"`
	// Image client type
	ImageClient string `yaml:"image_client"`
	// Server settings
	Server Server
	// Additional dynamically typed configuration
	Additional map[string]interface{} `yaml:"additional"`
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

// Get returns *info with value from Additional property. If it's not there default value or nil will be
// returned
func (c *Config) Get(key string, defValue ...interface{}) *info {
	v, ok := c.Additional[key]

	if !ok && len(defValue) == 0 {
		return &info{}
	}

	if !ok {
		v = defValue[0]
	}

	return &info{Key: key, value: v}
}

// Hash checks whether config has a value with given key in Additional property
func (c *Config) Has(key string) bool {
	_, ok := c.Additional[key]
	return ok
}

// Check checks validity of config
func (c *Config) Check() error {
	if c.ENV == "" {
		c.ENV = defaultENV
	}
	if c.Server.HTTPAddr == "" {
		c.Server.HTTPAddr = defaultAddr
	}
	if c.Server.HTTPPort == "" {
		c.Server.HTTPPort = defaultPort
	}
	if c.Storage == "" {
		c.Storage = defaultStorage
	}

	if c.ResizeStrategy == "" {
		c.ResizeStrategy = defaultResizeStrategy
	}

	if c.ImageClient == "" {
		c.ImageClient = defaultImageClient
	}

	return nil
}

// Load loads configuration from file
func (c *Config) Load(path string) error {
	var (
		err error
		b   []byte
	)

	b, err = ioutil.ReadFile(path)

	if err != nil {
		return fmt.Errorf("failed to open config %s\n%s", path, err.Error())
	}

	if err := yaml.Unmarshal(b, &c); err != nil {
		return fmt.Errorf("failed to unmarshal config file %s\n%s", path, err)
	}

	if err := c.Check(); err != nil {
		return err
	}

	return nil
}

// MustLoad initializes configuration from given config path
func MustLoad(path string) *Config {
	var c Config

	if err := c.Load(path); err != nil {
		panic(err)
	}

	return &c
}
