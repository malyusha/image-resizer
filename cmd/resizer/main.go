package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"

	"github.com/malyusha/image-resizer/pkg/app"
	"github.com/malyusha/image-resizer/pkg/config"
	"github.com/malyusha/image-resizer/pkg/server"
)

var (
	flagSet    = flag.NewFlagSet("resizer", flag.ExitOnError)
	configFile string
)

func init() {
	flagSet.StringVar(&configFile, "c", os.Getenv("RESIZER_CONFIG_FILE"), "Image resizer configuration file path")

	flagSet.Usage = usage
}

func main() {
	_ = flagSet.Parse(os.Args[1:])
	if configFile == "" || flagSet.Arg(0) == "help" {
		flagSet.Usage()
		return
	}

	application := app.CreateInstance(config.MustLoad(configFile))
	srv := server.NewServer(application)

	// Start srv
	errCh := srv.Start()

	// Wait for os notification in channel
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, os.Kill)

	// Graceful handle fatal errors in running app
	logrus.RegisterExitHandler(func() {
		stop <- os.Interrupt
	})

	select {
	case err := <-errCh:
		// Server startup error occurred
		application.Logger().Errorf("Error running srv: %s", err.Error())
		srv.Shutdown()
	case sig := <-stop:
		// Received system signal
		application.Logger().Warnf("Received %s signal. Stopping...", sig.String())
		srv.Shutdown()
		fmt.Println("done")
	}
}

// Usage prints usage for seeder in CLI
func usage() {
	_, _ = fmt.Fprintf(flagSet.Output(), "Usage: %s [options]\n", flagSet.Name())
	flagSet.PrintDefaults()
}

func exitLoggerFunc(stop chan os.Signal) func(code int) {
	return func(code int) {
		stop <- os.Kill
	}
}