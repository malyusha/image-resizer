package main

import (
	"net/http"
	"os"
	"fmt"

	"github.com/malyusha/image-resizer/app"
	"github.com/malyusha/image-resizer/console"
	log "github.com/sirupsen/logrus"
)

var (
	application app.Application
)

func main() {
	defer application.Destroy()

	addr := fmt.Sprintf("%v:%v", console.Args.Address, console.Args.Port)
	log.Infof("Serving application on %s", addr)
	log.Fatal(http.ListenAndServe(addr, application.Router()))
}

func init() {
	presetsFile := console.Args.PresetsFile

	if presetsFile == "" {
		exit("You should provide presets JSON file path to run application\n")
	}

	if _, err := os.Stat(presetsFile); os.IsNotExist(err) {
		exit("File %s doesn't exist\n", presetsFile)
	}

	application = app.GetInstance()
}

func exit(mess string, args ...interface{}) {
	if len(args) == 0 {
		fmt.Printf(mess)
	} else {
		fmt.Printf(mess, args...)
	}

	os.Exit(1)
}