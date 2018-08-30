package storage

import (
	"sync"
	"os"
	log "github.com/sirupsen/logrus"
	"fmt"
	"errors"
	"path"
	"io/ioutil"
	"github.com/malyusha/image-resizer/app/util"
)

type LocalStorageConfig struct {
	Dir string
}

type localStorage struct {
	mux sync.Mutex
	config *LocalStorageConfig
}

var (
	once    sync.Once
	storage *localStorage
)

const (
	DirectoryPerm = 0755
)

// New creates new LocalStorage and initializes it with given config
func NewLocalStorage(config *LocalStorageConfig) (*localStorage, error) {
	var (
		syncError error
	)
	once.Do(func() {
		if err := createDirectoryIfNotExists(config.Dir); err != nil {
			syncError = err
			return
		}

		log.Info(fmt.Sprintf("Using directory '%s' for local storage driver", config.Dir))

		storage = &localStorage{config: config}
	})

	return storage, syncError
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

// creates given directory recursively if it doesn't exist
func createDirectoryIfNotExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Info(fmt.Sprintf("Directory '%s' doesn't exist. Trying to create...\n", dir))
		if err := os.MkdirAll(dir, DirectoryPerm); err != nil {
			return err
		}
	}

	return nil
}
