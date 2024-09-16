package grpc

import (
	"context"

	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/api/pb"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/mapper"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/service"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EventsHandler struct {
	log     *logger.Logger
	service service.EventService
	pb.UnimplementedEventServiceServer
}

func NewEventsHandler(logger *logger.Logger, s service.EventService) *EventsHandler {
	return &EventsHandler{log: logger, service: s}
}

func (s *EventsHandler) CreateEvent(ctx context.Context, req *pb.CreateEventRequest) (*pb.EventResponse, error) {
	event := mapper.GRPCEventCreate(req)
	newEvent, err := s.service.InsertEvent(ctx, event)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return mapper.EventGRPC(newEvent), nil
}

func (s *EventsHandler) UpdateEvent(ctx context.Context, req *pb.UpdateEventRequest) (*pb.EventResponse, error) {
	event := mapper.GRPCEventUpdate(req)
	newEvent, err := s.service.UpdateEvent(ctx, event)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return mapper.EventGRPC(newEvent), nil
}

func (s *EventsHandler) RemoveEvent(ctx context.Context, req *pb.RemoveEventRequest) (*pb.RemoveEventResponse, error) {
	err := s.service.RemoveEvent(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return &pb.RemoveEventResponse{Id: req.Id}, nil
}

func (s *EventsHandler) FindEvent(ctx context.Context, req *pb.FindEventRequest) (*pb.EventResponse, error) {
	event, err := s.service.FindEvent(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return mapper.EventGRPC(event), nil
}

func (s *EventsHandler) FindEventsByDay(ctx context.Context, req *pb.GetEventByDate) (*pb.EventsResponse, error) {
	events, err := s.service.FindAllEventsByDay(ctx, req.Date.AsTime())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return mapper.EventsGRPC(events), nil
}

func (s *EventsHandler) FindEventsByWeek(ctx context.Context, req *pb.GetEventByDate) (*pb.EventsResponse, error) {
	events, err := s.service.FindAllEventsByWeek(ctx, req.Date.AsTime())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return mapper.EventsGRPC(events), nil
}

func (s *EventsHandler) FindEventsByMonth(ctx context.Context, req *pb.GetEventByDate) (*pb.EventsResponse, error) {
	events, err := s.service.FindAllEventsByMonth(ctx, req.Date.AsTime())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return mapper.EventsGRPC(events), nil
}

func (s *EventsHandler) FindEvents(ctx context.Context, _ *empty.Empty) (*pb.EventsResponse, error) {
	events, err := s.service.FindAllEvents(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return mapper.EventsGRPC(events), nil
}
