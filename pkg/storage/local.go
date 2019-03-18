package storage

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"sync"

	"github.com/malyusha/image-resizer/pkg/util"
)

type LocalStorageConfig struct {
	Dir string
}

type localStorage struct {
	mux    sync.Mutex
	config *LocalStorageConfig
}

const (
	DirectoryPerm = 0755
)

// New creates new LocalStorage and initializes it with given config
func NewLocalStorage(config *LocalStorageConfig) (*localStorage, error) {
	var (
		syncError error
	)

	if config.Dir == "" {
		return nil, errors.New(`when using "local" type of storage you must provide storage directory for files`)
	}

	if err := createDirectoryIfNotExists(config.Dir); err != nil {
		return nil, err
	}

	return &localStorage{config: config}, syncError
}

// Save saves file into file system
func (s *localStorage) Save(filename string, content []byte) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if err := checkFile(filename); err != nil {
		return err
	}

	// Creating full path for file
	fullPath := s.fullPath(filename)

	// Next, we'll need to create directory for file if it doesn't exist
	if err := createDirectoryIfNotExists(path.Dir(fullPath)); err != nil {
		return err
	}

	// Writing content to file
	if err := ioutil.WriteFile(fullPath, content, 0644); err != nil {
		return err
	}

	return nil
}

// Get returns data from file existing in file system. If given file doesn't exist it will
// return an error
func (s *localStorage) Get(filename string) ([]byte, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	return util.GetFileContent(s.fullPath(filename))
}

// Delete removes file with given name from filesystem
func (s *localStorage) Delete(filename string) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if err := os.Remove(s.fullPath(filename)); err != nil {
		return err
	}

	return nil
}

// Dir returns config base directory of storage. Just a getter for config parameter.
func (s *localStorage) Dir() string {
	return s.config.Dir
}

// fullPath returns full path of given filename including base config directory for current storage
func (s *localStorage) fullPath(filename string) string {
	return path.Join(s.Dir(), filename)
}

// checkFile checks if file name is empty string
func checkFile(filename string) error {
	if filename == "" {
		return errors.New("filename is empty")
	}

	return nil
}

// createDirectoryIfNotExists creates given directory recursively if it doesn't exist
func createDirectoryIfNotExists(dir string) error {
	if _, err := os.Stat(dir); os.IsExist(err) {
		return nil
	}

	if err := os.MkdirAll(dir, DirectoryPerm); err != nil {
		return err
	}

	return nil
}
