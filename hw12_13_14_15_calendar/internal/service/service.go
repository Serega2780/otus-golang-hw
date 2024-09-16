package service

import (
	"context"
	"time"

	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/model"
)

type EventService interface {
	RemoveEvent(ctx context.Context, id string) (err error)
	FindEvent(ctx context.Context, id string) (*model.Event, error)
	FindAllEvents(ctx context.Context) ([]*model.Event, error)
	FindAllEventsByDay(ctx context.Context, date time.Time) ([]*model.Event, error)
	FindAllEventsByWeek(ctx context.Context, date time.Time) ([]*model.Event, error)
	FindAllEventsByMonth(ctx context.Context, date time.Time) ([]*model.Event, error)
	UpdateEvent(ctx context.Context, event *model.Event) (updatedEvent *model.Event, err error)
	InsertEvent(ctx context.Context, event *model.Event) (newEvent *model.Event, err error)
}
