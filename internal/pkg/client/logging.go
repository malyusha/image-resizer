package client

import (
	"github.com/malyusha/image-resizer/internal/pkg/logger"
)

// LoggingClient is decorator for Client. It logs image retrieving, using given logger
type LoggingClient struct {
	Client Client
	Logger logger.Logger
}

func (c *LoggingClient) GetImageContent(path string) ([]byte, error) {
	c.Logger.Infof(`Retrieving image content for path "%s"`, path)

	b, err := c.Client.GetImageContent(path)

	if err != nil {
		c.Logger.Errorf("Client error: %s", err)
	}

	return b, err
}

func (c *LoggingClient) Path(path string) string {
	return c.Client.Path(path)
}

// WithLogger returns decorated Client, that will log it's activity.
func WithLogger(c Client, logger logger.Logger) Client {
	return &LoggingClient{c, logger}
}
