package internalhttp

import (
	"net/http"
	"time"

	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
)

type status string

const (
	STATUS status = "status"
)

func loggingMiddleware(next http.Handler, logger *logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		logger.Infof("%s %s %s %s %s %d %v %s",
			r.RemoteAddr, start.Format(time.RFC822), r.Method, r.URL.Path, r.Proto, r.Context().Value(STATUS),
			time.Since(start), r.UserAgent())
	})
}
