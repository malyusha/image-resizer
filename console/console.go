package console

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/malyusha/image-resizer/config"
	"github.com/joho/godotenv"
)

// Defines console arguments struct
type argv struct {
	Env         string
	Port        string
	Address     string
	PresetsFile string
	Storage     string
	StorageDir  string
	ImageClient string
	SourceDir   string
	LogLevel    string
	help        bool
}

var Args *argv

func init() {
	loadEnvFiles()

	Args = &argv{}

	flag.Usage = func() {
		name := filepath.Base(os.Args[0])
		fmt.Fprintf(os.Stdout, "Usage: %s [options]\n", name)
		flag.PrintDefaults()
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "fatal"
	}

	flag.BoolVar(&Args.help, "h", false, "Show this help")
	flag.StringVar(&Args.Address, "addr", "127.0.0.1", "Server address")
	flag.StringVar(&Args.Port, "port", "9898", "Server port")
	flag.StringVar(&Args.Env, "e", os.Getenv("ENV"), "Environment")
	flag.StringVar(&Args.PresetsFile, "presets", os.Getenv("RESIZE_PRESETS_PATH"), "JSON file containing presets for images storage.")
	flag.StringVar(&Args.Storage, "s", os.Getenv("STORAGE_NAME"), "Storage name for application")
	flag.StringVar(&Args.StorageDir, "d", os.Getenv("STORAGE_DIR"), "Storage directory for application. Used only when storage type is local.")
	flag.StringVar(&Args.ImageClient, "client", os.Getenv("IMAGES_CLIENT"), "Images client for application")
	flag.StringVar(&Args.SourceDir, "source", os.Getenv("SOURCE_DIR"), "Source images directory for application. Used only when images client type is local.")
	flag.StringVar(&Args.LogLevel, "log", logLevel, "Log level for application. Allowed types: panic, fatal, error, warn, warning, info, debug.")
	flag.Parse()

	if Args.help {
		flag.Usage()
		os.Exit(0)
	}
}

func loadEnvFiles() {
	var filteredFiles []string
	envFiles := []string{".env"}

	if config.IsTesting() {
		envFiles = append(envFiles, ".testing.env")
	}

	env := strings.Split(strings.Trim(os.Getenv("ENV"), " "), ",")

	if len(env) != 0 && env[0] != "" {
		envFiles = append(envFiles, env...)
	}

	for _, file := range envFiles {
		path := config.RootPath(fmt.Sprintf("env/%s", file))
		file, err := os.Open(path)

		if err != nil {
			continue
		}

		file.Close()
		filteredFiles = append(filteredFiles, path)
	}

	if len(filteredFiles) == 0 {
		return
	}

	if err := godotenv.Overload(filteredFiles...); err != nil {
		panic(err)
	}
}
