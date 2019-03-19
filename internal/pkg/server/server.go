package server

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"

	"github.com/malyusha/image-resizer/internal/app"
)

type Instance struct {
	app          app.Application
	httpServer   *http.Server
	router       *mux.Router
	shuttingDown uint32
}

// App returns application instance of server
func (s *Instance) App() app.Application {
	return s.app
}

// Start starts server and returns error channel. If any error occurred channel will receive
// error that can be processed.
func (s *Instance) Start() chan error {
	errChan := make(chan error, 1)

	go s.StartHTTP(errChan)

	return errChan
}

// StartHTTP starts HTTP listener on configured address and port.
func (s *Instance) StartHTTP(errChan chan<- error) {
	log := s.app.Logger()
	log.Infof("Starting HTTP listener on %s", s.app.Config().Server.Address())
	s.httpServer = &http.Server{
		Addr:         s.app.Config().Server.Address(),
		Handler:      s.router,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	if err := s.httpServer.ListenAndServe(); err != nil {
		errChan <- fmt.Errorf("http server startup error: %s", err)
	}
}

// Shutdowns running servers
func (s *Instance) Shutdown() {
	log := s.app.Logger()

	if atomic.LoadUint32(&s.shuttingDown) == 1 {
		log.Warnf("Shutdown already in progress")
		return
	} else {
		atomic.AddUint32(&s.shuttingDown, 1)
	}

	if s.httpServer != nil {
		wait := s.app.Config().Server.GetGracefulTimeout()
		ctx, cancel := context.WithTimeout(context.Background(), wait)
		defer cancel()

		log.Info("Stopping http server...")
		if err := s.httpServer.Shutdown(ctx); err != nil {
			log.Errorf("failed to shutdown http server: %s", err.Error())
		}
	}

	atomic.AddUint32(&s.shuttingDown, 0)
	log.Info("Shutdown successfully completed")
}

// NewInstance returns new instance of server
func NewInstance(app app.Application) *Instance {
	server := &Instance{app: app, router: mux.NewRouter()}
	server.registerRoutes()

	return server
}
