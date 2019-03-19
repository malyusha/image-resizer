package server

func (s *Instance) registerRoutes() {
	var (
		storage = s.app.Storage()
		client = s.app.ImageClient()
	)

	s.router.HandleFunc("/{preset}/{path:.+}", s.HandleImagesRequest(storage, client))
	s.router.HandleFunc("/health", s.HandleHealthCheck())

	s.router.NotFoundHandler = s.HandleNotFound()

	s.router.Use(s.LogsRequests)
}

