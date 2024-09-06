package internalhttp

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
)

type Server struct {
	ip   string
	port string
	app  *app.App
	log  *logger.Logger
	srv  *http.Server
}

func NewServer(conf *config.HTTPServerConfig, logger *logger.Logger, app *app.App) *Server {
	return &Server{app: app, log: logger, ip: conf.IP, port: conf.Port}
}

func (s *Server) Start(_ context.Context) error {
	h := NewService(s.log)

	mux := http.NewServeMux()
	mux.HandleFunc("/", h.HelloWorld)
	mux.HandleFunc("/hello", h.HelloWorld)

	server := &http.Server{
		Addr:              strings.Join([]string{s.ip, s.port}, ":"),
		Handler:           loggingMiddleware(mux, s.log),
		ReadHeaderTimeout: 2 * time.Second,
	}
	s.srv = server

	s.log.Infof("server start on port %s", s.port)
	log.Fatal(server.ListenAndServe())

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.srv.Shutdown(ctx)
	if err != nil {
		s.log.Errorf("could not shutdown http server %v", err)
	}
	return nil
}
