package client

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

// HTTPClient represents client, with remote retrieving image content from URL.
// If BaseURL property is set, it will be used as base host for image retrieving.
type HTTPClient struct {
	BaseURL string
	client  *http.Client
}

// GetImageContent returns image content from remote server by given path.
func (c *HTTPClient) GetImageContent(filepath string) ([]byte, error) {
	urlString := c.Path(filepath)
	parsed, err := url.Parse(urlString)
	var b []byte

	if parsed.Host == "" {
		return b, errors.New("no host provided in HTTP Image client")
	}

	if err != nil {
		return b, fmt.Errorf("failed to parse valid url: %s", err)
	}

	// Requesting image content from URL
	resp, err := c.client.Get(urlString)

	if err != nil {
		return nil, fmt.Errorf("failed to load image content: %s", err)
	}

	// Maybe handle error?
	defer resp.Body.Close()

	// TODO: enhance status code error handling
	if resp.StatusCode >= 400 || resp.StatusCode < 200 {
		return []byte{}, fmt.Errorf("server responded with unsuccessfull code: %d", resp.StatusCode)
	}

	// Returning reader result
	return ioutil.ReadAll(resp.Body)
}

func (c *HTTPClient) Path(filepath string) string {
	return c.BaseURL + filepath
}

// NewHTTPClient returns new http client for image content retrieve
func NewHTTPClient(config *HTTPClientConfig) (*HTTPClient, error) {
	if config.BaseURL == "" {
		return nil, errors.New("base url is required for http client")
	}

	return &HTTPClient{
		BaseURL: config.BaseURL,
		client:  createClient(),
	}, nil
}

// getClient creates new http client with timeout parameters only once
func createClient() *http.Client {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).DialContext,
	}

	return &http.Client{
		Timeout:   5 * time.Second,
		Transport: transport,
	}
}
