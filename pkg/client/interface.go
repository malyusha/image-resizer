package client

import (
	log "github.com/sirupsen/logrus"
)

// Client is the client retrieving content of source image file.
type Client interface {
	// GetImageContent returns data of given file. If file doesn't exist, or other error occurred
	// error must be returned.
	GetImageContent(path string) ([]byte, error)

	// FullPath path returns string, representing normalized path for client.
	// (e.g. it can be full path with host prefix for http client)
	FullPath(path string) string
}

// LoggingClient is decorator for Client. It logs image retrieving, using given logger
type LoggingClient struct {
	Client Client
	Logger *log.Logger
}

func (c *LoggingClient) GetImageContent(path string) ([]byte, error) {
	c.Logger.Infof(`Retrieving image content for path "%s"`, path)

	b, err := c.Client.GetImageContent(path)

	if err != nil {
		c.Logger.WithError(err).Error("Client error")
	}

	return b, err
}

func (c *LoggingClient) FullPath(path string) string {
	return c.Client.FullPath(path)
}

// WithLogger returns decorated Client, that will log it's activity.
func WithLogger(c Client, logger *log.Logger) Client {
	return &LoggingClient{c, logger}
}
