package mapper

import (
	"time"

	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/api/pb"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Event(dbEvent *model.DBEvent) *model.Event {
	if dbEvent == nil {
		return nil
	}
	return &model.Event{
		ID: dbEvent.ID, Title: dbEvent.Title, StartTime: dbEvent.StartTime, EndTime: dbEvent.EndTime,
		Description: dbEvent.Description, UserID: dbEvent.UserID, NotifyBeforeEvent: dbEvent.NotifyBeforeEvent,
	}
}

func EventGRPC(event *model.Event) *pb.EventResponse {
	if event == nil {
		return nil
	}
	return &pb.EventResponse{Event: &pb.Event{
		Id: event.ID, Title: event.Title,
		StartTime: timestamppb.New(event.StartTime), EndTime: timestamppb.New(event.EndTime),
		Description: event.Description, UserId: event.UserID,
		NotifyBeforeEvent: event.NotifyBeforeEvent.Milliseconds(),
	}}
}

func GRPCEventCreate(ge *pb.CreateEventRequest) *model.Event {
	return &model.Event{
		Title: ge.Event.Title, StartTime: ge.Event.StartTime.AsTime(),
		EndTime: ge.Event.EndTime.AsTime(), Description: ge.Event.Description, UserID: ge.Event.UserId,
		NotifyBeforeEvent: time.Duration(ge.Event.NotifyBeforeEvent),
	}
}

func GRPCEventUpdate(ge *pb.UpdateEventRequest) *model.Event {
	return &model.Event{
		ID: ge.Event.Id, Title: ge.Event.Title, StartTime: ge.Event.StartTime.AsTime(),
		EndTime: ge.Event.EndTime.AsTime(), Description: ge.Event.Description, UserID: ge.Event.UserId,
		NotifyBeforeEvent: time.Duration(ge.Event.NotifyBeforeEvent),
	}
}

func Events(dbEvents []*model.DBEvent) []*model.Event {
	events := make([]*model.Event, len(dbEvents))
	if len(dbEvents) == 0 {
		return make([]*model.Event, 0)
	}
	for i, dbe := range dbEvents {
		e := Event(dbe)
		events[i] = e
	}
	return events
}

func EventsGRPC(events []*model.Event) *pb.EventsResponse {
	response := &pb.EventsResponse{}
	if len(events) == 0 {
		return response
	}
	for _, e := range events {
		response.Events = append(response.Events, &pb.Event{
			Id: e.ID, Title: e.Title,
			StartTime: timestamppb.New(e.StartTime), EndTime: timestamppb.New(e.EndTime),
			Description: e.Description, UserId: e.UserID,
			NotifyBeforeEvent: e.NotifyBeforeEvent.Milliseconds(),
		})
	}
	return response
}

func DBEvent(event *model.Event) *model.DBEvent {
	return &model.DBEvent{
		ID: event.ID, Title: event.Title, StartTime: event.StartTime, EndTime: event.EndTime,
		Description: event.Description, UserID: event.UserID, NotifyBeforeEvent: event.NotifyBeforeEvent,
	}
}

func DBEvents(events []*model.Event) []*model.DBEvent {
	dbEvents := make([]*model.DBEvent, len(events))
	for i, e := range events {
		dbe := DBEvent(e)
		dbEvents[i] = dbe
	}
	return dbEvents
}
