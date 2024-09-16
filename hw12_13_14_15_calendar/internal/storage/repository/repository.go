package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/model"
)

var (
	ErrDateBusy = errors.New("specified event time range is already occupied error")
	ErrIDExists = errors.New("insert entity with non empty id error")
	ErrBeginTx  = errors.New("transaction begin error")
	ErrCommitTx = errors.New("transaction commit error")
)

type EventInterface interface {
	Remove(ctx context.Context, id string) (err error)
	Find(ctx context.Context, id string) (*model.DBEvent, error)
	FindAll(ctx context.Context) ([]*model.DBEvent, error)
	FindAllByDay(ctx context.Context, date time.Time) ([]*model.DBEvent, error)
	FindAllByWeek(ctx context.Context, date time.Time) ([]*model.DBEvent, error)
	FindAllByMonth(ctx context.Context, date time.Time) ([]*model.DBEvent, error)
	Update(ctx context.Context, event *model.DBEvent) (updatedEvent *model.DBEvent, err error)
	Insert(ctx context.Context, event *model.DBEvent) (newEvent *model.DBEvent, err error)
}
