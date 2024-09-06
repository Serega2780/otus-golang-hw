package repository

import (
	"context"
	"errors"

	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/storage/model"
)

var (
	ErrDateBusy = errors.New("specified event time range is already occupied error")
	ErrIDExists = errors.New("insert entity with non empty id error")
	ErrBeginTx  = errors.New("transaction begin error")
	ErrCommitTx = errors.New("transaction commit error")
)

type EventInterface interface {
	Remove(ctx context.Context, id string) (err error)
	Find(ctx context.Context, id string) (*model.Event, error)
	FindAll(ctx context.Context) ([]model.Event, error)
	Update(ctx context.Context, event *model.Event) (updatedEvent *model.Event, err error)
	Insert(ctx context.Context, event *model.Event) (newEvent *model.Event, err error)
}
