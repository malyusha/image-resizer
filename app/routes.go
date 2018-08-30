package app

func (a *app) routes() {
	a.router.HandleFunc("/{preset}/{path:.+}", a.HandleImagesRequest())
	a.router.HandleFunc("/health", a.HandleHealthCheck())

	a.router.NotFoundHandler = a.HandleNotFound()
}