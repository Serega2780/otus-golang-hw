package internalhttp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
)

type Service struct {
	log *logger.Logger
}

func NewService(logger *logger.Logger) *Service {
	return &Service{log: logger}
}

func (s *Service) HelloWorld(w http.ResponseWriter, r *http.Request) {
	resp := "Hello, World!"
	if r.Method != http.MethodGet {
		resp = fmt.Sprintf("method %s not not supported on uri %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusMethodNotAllowed)
		*r = *(r.WithContext(context.WithValue(r.Context(), STATUS, http.StatusMethodNotAllowed)))
		s.writeResponse(w, resp)
		return
	}
	w.WriteHeader(http.StatusOK)
	*r = *(r.WithContext(context.WithValue(r.Context(), STATUS, http.StatusOK)))
	s.writeResponse(w, resp)
}

func (s *Service) writeResponse(w http.ResponseWriter, resp string) {
	_, err := w.Write([]byte(resp))
	if err != nil {
		s.log.Errorf("response marshal error: %s", err)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}
