package internalhttp

import (
	"net/http"
	"time"
)

type status string

const (
	STATUS status = "status"
)

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		s.log.Infof("%s %s %s %s %s %d %v %s",
			r.RemoteAddr, start.Format(time.RFC822), r.Method, r.URL.Path, r.Proto, r.Context().Value(STATUS),
			time.Since(start), r.UserAgent())
	})
}
