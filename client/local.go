package client

import (
	"github.com/pkg/errors"
	"path"
	"github.com/malyusha/image-resizer/app/util"
)

var (
	ErrorDirectoryNotExists = errors.New("Given directory does not exist")
)

type LocalStorageClient struct {
	Directory string
}

func (c *LocalStorageClient) GetImageContent(path string) ([]byte, error) {
	if c.Directory == "" {
		return nil, ErrorDirectoryNotExists
	}

	return util.GetFileContent(c.absPath(path))
}

func (c *LocalStorageClient) absPath(filename string) (string) {
	return path.Join(c.Directory, filename)
}

func NewLocalStorageClient(dir string) *LocalStorageClient {
	if dir == "" {
		dir = "/"
	}

	return &LocalStorageClient{dir}
}

