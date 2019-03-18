package server

import (
	"net/http"
)

// Represents basic request logging
func (s *Server) LogsRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.app.Logger().Infof("Incoming request: Addr - %s, Method - %s, URI - %s", r.RemoteAddr, r.Method, r.RequestURI)

		next.ServeHTTP(w, r)
	})
}
