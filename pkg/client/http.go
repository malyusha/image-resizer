package client

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

var (
	// Single http.Client instance
	client *http.Client
	once sync.Once
)

// HTTPImageClient represents client, with remote retrieving image content from URL.
// If Host property is set, it will be used as base host for image retrieving.
type HTTPImageClient struct {
	Host string
	client *http.Client
}

// GetImageContent returns image content from remote server by given path.
func (c *HTTPImageClient) GetImageContent(path string) ([]byte, error) {
	urlString := c.url(path)
	parsed, err := url.Parse(urlString)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to parse valid url: %s", err)
	}

	// Requesting image content from URL
	resp, err := c.client.Get(fmt.Sprintf("%s/%s", parsed.Host, parsed.Path))

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

func (c *HTTPImageClient) FullPath(path string) string {
	return c.url(path)
}

// url returns URL for given path with host
func (c *HTTPImageClient) url(path string) string {
	if strings.Contains(path, "http://") || strings.Contains(path, "https://") || strings.Index(path, "//") == 0 {
		return path
	}

	return fmt.Sprintf("%s/%s", c.Host, path)
}

// NewHTTPImageClient returns new http client for image content retrieve
func NewHTTPImageClient(host string) *HTTPImageClient {
	return &HTTPImageClient{
		Host: host,
		client: getClient(),
	}
}

// getClient creates new http client with timeout parameters only once
func getClient() *http.Client {
	once.Do(func() {
		transport := &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).DialContext,
		}

		client = &http.Client{
			Timeout: 5 * time.Second,
			Transport: transport,
		}
	})

	return client
}