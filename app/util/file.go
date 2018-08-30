package util

import (
	"github.com/pkg/errors"
	"os"
	"fmt"
	"io/ioutil"
)

// GetFileContent returns data from given file if it exists. Otherwise returns error
func GetFileContent(file string) ([]byte, error) {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return nil, errors.New(fmt.Sprintf("File %s does not exist", file))
	}

	// Using read file because we need to read whole file
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return data, nil
}