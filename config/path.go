package config

import (
	"os"
	"path/filepath"
	"runtime"
)

var (
	root       string
)

const (
	envAppRootPathName = "APP_ROOT_PATH"
)

func RootPath(additional ...string) string {
	path := root

	if envRootPath := os.Getenv(envAppRootPathName); envRootPath != "" {
		root = envRootPath
		path = root
	}

	if path == "" {
		_, filename, _, ok := runtime.Caller(0)
		if !ok {
			panic("No caller for root path")
		}

		root = filepath.Dir(filepath.Dir(filename))
		path = root
	}

	if len(additional) > 0 {
		path = filepath.Join(append([]string{root}, additional...)...)
	}

	return path
}