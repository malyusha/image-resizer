package client

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/malyusha/image-resizer/pkg/dot"
	"github.com/malyusha/image-resizer/pkg/util"
)

type LocalStorageClient struct {
	Directory string
}

func (c *LocalStorageClient) GetImageContent(path string) ([]byte, error) {
	return util.GetFileContent(c.absPath(path))
}

func (c *LocalStorageClient) FullPath(path string) string {
	return c.absPath(path)
}

// absPath returns absolute path with storage directory for given file
func (c *LocalStorageClient) absPath(filename string) string {
	return path.Join(c.Directory, filename)
}

// Returns new LocalStorageClient
func NewLocalStorageClient(config *dot.Map) (*LocalStorageClient, error) {
	dir := config.Get("dir").String()
	if dir == "" {
		return nil, errors.New(`when using "local" type of images client you must provide source directory for static files`)
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, fmt.Errorf(`directory "%s", passed to LocalStorageClient doesn't exist`, dir)
	}

	return &LocalStorageClient{dir}, nil
}
