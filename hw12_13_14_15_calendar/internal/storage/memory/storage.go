package memorystorage

import (
	"context"
	"sync"

	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/storage/model"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/storage/repository"
	"github.com/google/uuid"
)

type Storage struct {
	db  map[string]*model.Event
	dmu sync.RWMutex
}

func New() *Storage {
	return &Storage{db: make(map[string]*model.Event)}
}

func (s *Storage) Remove(_ context.Context, id string) (err error) {
	s.dmu.Lock()
	defer s.dmu.Unlock()
	delete(s.db, id)
	return nil
}

func (s *Storage) Find(_ context.Context, id string) (*model.Event, error) {
	s.dmu.RLock()
	defer s.dmu.RUnlock()
	return s.db[id], nil
}

func (s *Storage) FindAll(_ context.Context) ([]model.Event, error) {
	s.dmu.RLock()
	defer s.dmu.RUnlock()
	v := make([]model.Event, 0, len(s.db))

	for _, value := range s.db {
		v = append(v, *value)
	}
	return v, nil
}

func (s *Storage) Update(_ context.Context, event *model.Event) (updatedEvent *model.Event, err error) {
	s.dmu.Lock()
	defer s.dmu.Unlock()
	if err := s.IsOccupied(event); err != nil {
		return nil, err
	}
	s.db[event.ID] = event
	return s.db[event.ID], nil
}

func (s *Storage) Insert(_ context.Context, event *model.Event) (newEvent *model.Event, err error) {
	s.dmu.Lock()
	defer s.dmu.Unlock()
	if err := s.IsOccupied(event); err != nil {
		return nil, err
	}
	if len(event.ID) != 0 {
		return nil, repository.ErrIDExists
	}
	event.ID = uuid.New().String()
	s.db[event.ID] = event
	return s.db[event.ID], nil
}

func (s *Storage) IsOccupied(event *model.Event) error {
	for id, ev := range s.db {
		if id == event.ID {
			continue
		}
		startEvent := event.StartTime
		endEvent := event.EndTime
		if (startEvent.After(ev.StartTime) || startEvent.Equal(ev.StartTime)) &&
			(startEvent.Before(ev.EndTime) || startEvent.Equal(ev.EndTime)) {
			return repository.ErrDateBusy
		}
		if (endEvent.After(ev.StartTime) || endEvent.Equal(ev.StartTime)) &&
			(endEvent.Before(ev.EndTime) || endEvent.Equal(ev.EndTime)) {
			return repository.ErrDateBusy
		}
		if startEvent.Before(ev.StartTime) && endEvent.After(ev.EndTime) {
			return repository.ErrDateBusy
		}
	}
	return nil
}
