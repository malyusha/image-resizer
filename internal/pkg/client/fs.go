package client

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/malyusha/image-resizer/internal/pkg/util"
)

type FileSystemClient struct {
	Directory string
}

func (c *FileSystemClient) GetImageContent(filepath string) ([]byte, error) {
	return util.GetFileContent(c.Path(filepath))
}

func (c *FileSystemClient) Path(filepath string) string {
	return path.Join(c.Directory, filepath)
}

// Returns new FileSystemClient
func NewFileSystemClient(config *FSClientConfig) (*FileSystemClient, error) {
	location := config.Location
	if location == "" {
		return nil, errors.New("file system client requires directory location of files to serve from")
	}

	if _, err := os.Stat(location); os.IsNotExist(err) {
		return nil, fmt.Errorf(`directory "%s", passed to FileSystemClient doesn't exist`, location)
	}

	return &FileSystemClient{location}, nil
}
