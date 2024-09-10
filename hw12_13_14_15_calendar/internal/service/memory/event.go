package memory

import (
	"context"
	"time"

	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/mapper"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/model"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/storage/memory"
)

type EventMemoryService struct {
	eventsRepo *memory.Storage
}

func NewEventMemoryService(repo *memory.Storage) *EventMemoryService {
	return &EventMemoryService{eventsRepo: repo}
}

func (ess *EventMemoryService) RemoveEvent(ctx context.Context, id string) (err error) {
	return ess.eventsRepo.Remove(ctx, id)
}

func (ess *EventMemoryService) FindEvent(ctx context.Context, id string) (*model.Event, error) {
	dbe, err := ess.eventsRepo.Find(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapper.Event(dbe), nil
}

func (ess *EventMemoryService) FindAllEvents(ctx context.Context) ([]*model.Event, error) {
	dbes, err := ess.eventsRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return mapper.Events(dbes), nil
}

func (ess *EventMemoryService) FindAllEventsByDay(ctx context.Context, date time.Time) ([]*model.Event, error) {
	dbes, err := ess.eventsRepo.FindAllByDay(ctx, date)
	if err != nil {
		return nil, err
	}
	return mapper.Events(dbes), nil
}

func (ess *EventMemoryService) FindAllEventsByWeek(ctx context.Context, date time.Time) ([]*model.Event, error) {
	dbes, err := ess.eventsRepo.FindAllByWeek(ctx, date)
	if err != nil {
		return nil, err
	}
	return mapper.Events(dbes), nil
}

func (ess *EventMemoryService) FindAllEventsByMonth(ctx context.Context, date time.Time) ([]*model.Event, error) {
	dbes, err := ess.eventsRepo.FindAllByMonth(ctx, date)
	if err != nil {
		return nil, err
	}
	return mapper.Events(dbes), nil
}

func (ess *EventMemoryService) UpdateEvent(ctx context.Context, event *model.Event) (updatedEvent *model.Event,
	err error,
) {
	dbe, err := ess.eventsRepo.Update(ctx, mapper.DBEvent(event))
	if err != nil {
		return nil, err
	}
	return mapper.Event(dbe), nil
}

func (ess *EventMemoryService) InsertEvent(ctx context.Context, event *model.Event) (newEvent *model.Event, err error) {
	dbe, err := ess.eventsRepo.Insert(ctx, mapper.DBEvent(event))
	if err != nil {
		return nil, err
	}
	return mapper.Event(dbe), nil
}
