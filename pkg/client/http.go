package client

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/malyusha/image-resizer/pkg/dot"
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
func (c *HTTPImageClient) GetImageContent(p string) ([]byte, error) {
	urlString := c.url(p)
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

func (c *HTTPImageClient) FullPath(path string) string {
	return c.url(path)
}

// url returns URL for given path with host
func (c *HTTPImageClient) url(p string) string {
	if strings.Contains(p, "http://") || strings.Contains(p, "https://") || strings.Index(p, "//:") == 0 {
		return p
	}
	if p[0] != '/' {
		p = "/" + p
	}

	p = path.Clean(p)

	if c.Host == "" {
		return p
	}


	return c.Host + p
}

// NewHTTPImageClient returns new http client for image content retrieve
func NewHTTPImageClient(config *dot.Map) (*HTTPImageClient, error) {
	return &HTTPImageClient{
		Host: config.Get("host").String(),
		client: getClient(),
	}, nil
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