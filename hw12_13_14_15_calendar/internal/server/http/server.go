package http

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/gorilla/mux"
)

type Server struct {
	ctx  context.Context
	ip   string
	port string
	app  *app.App
	log  *logger.Logger
	srv  *http.Server
}

func NewServer(ctx context.Context, app *app.App) *Server {
	return &Server{ctx: ctx, app: app, log: app.Log, ip: app.Cfg.HTTP.IP, port: app.Cfg.HTTP.Port}
}

func (s *Server) Start(ctx context.Context) {
	h := NewEventsHandler(ctx, s.log, s.app.EventsService)

	r := mux.NewRouter()
	r.HandleFunc("/", h.HelloWorld).Methods(http.MethodGet)
	r.HandleFunc("/hello", h.HelloWorld).Methods(http.MethodGet)
	r.HandleFunc("/v1/events", h.InsertEvent).Methods(http.MethodPost)
	r.HandleFunc("/v1/events", h.UpdateEvent).Methods(http.MethodPut)
	r.HandleFunc("/v1/events/{id}", h.DeleteEvent).Methods(http.MethodDelete)
	r.HandleFunc("/v1/events/{id}", h.FindEventByID).Methods(http.MethodGet)
	r.HandleFunc("/v1/events", h.FindAllEvent).Methods(http.MethodGet)
	r.HandleFunc("/v1/events/day/{date}", h.FindAllEventByDay).Methods(http.MethodGet)
	r.HandleFunc("/v1/events/week/{date}", h.FindAllEventByWeek).Methods(http.MethodGet)
	r.HandleFunc("/v1/events/month/{date}", h.FindAllEventByMonth).Methods(http.MethodGet)
	r.Use(s.loggingMiddleware)
	server := &http.Server{
		Addr:              strings.Join([]string{s.ip, s.port}, ":"),
		Handler:           r,
		ReadHeaderTimeout: 2 * time.Second,
	}
	s.srv = server
	s.ctx = ctx

	s.log.Infof("http server start on port %s", s.port)
	go func() {
		_ = server.ListenAndServe()
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	s.Stop(ctx)
}

func (s *Server) Stop(ctx context.Context) {
	err := s.srv.Shutdown(ctx)
	if err != nil {
		s.log.Errorf("could not shutdown http server %v", err)
		return
	}
	s.log.Info("http server stopped")
}
