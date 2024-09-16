package grpc

import (
	"context"
	"net"
	"strings"

	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/api/pb"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/go-kit/log"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	ctx  context.Context
	ip   string
	port string
	app  *app.App
	log  *logger.Logger
	srv  *grpc.Server
}

func NewServer(ctx context.Context, app *app.App) *Server {
	return &Server{ctx: ctx, app: app, log: app.Log, ip: app.Cfg.GRPC.IP, port: app.Cfg.GRPC.Port}
}

func (s *Server) Start(ctx context.Context) {
	h := NewEventsHandler(s.log, s.app.EventsService)
	l := log.NewJSONLogger(s.log.GetWriter())
	rpcLogger := log.With(l, "service", "gRPC/server")
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(interceptorLogger(rpcLogger),
				logging.WithFieldsFromContext(logging.ExtractFields)),
		),
	)
	pb.RegisterEventServiceServer(grpcServer, h)
	grpcListener, err := net.Listen("tcp", strings.Join([]string{s.ip, s.port}, ":"))
	if err != nil {
		s.log.Errorf("error listen on port %s %v", s.port, err)
	}
	s.srv = grpcServer
	reflection.Register(grpcServer)
	s.ctx = ctx
	s.log.Infof("grpc server start on port %s", s.port)
	go func() {
		_ = grpcServer.Serve(grpcListener)
	}()

	<-ctx.Done()
	s.Stop(ctx)
}

func (s *Server) Stop(_ context.Context) {
	s.srv.Stop()
	s.log.Info("grpc server stopped")
}
