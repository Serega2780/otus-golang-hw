package memorystorage

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/model"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/storage/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var (
	ctx     = context.Background()
	tm, _   = time.Parse(time.RFC3339, "2024-08-28T10:00:00Z")
	tm2, _  = time.Parse(time.RFC3339, "2024-08-28T12:00:00Z")
	tm3, _  = time.Parse(time.RFC3339, "2024-08-28T14:00:00Z")
	tmE, _  = time.Parse(time.RFC3339, "2024-08-28T11:00:00Z")
	tmE2, _ = time.Parse(time.RFC3339, "2024-08-28T13:00:00Z")
	tmE3, _ = time.Parse(time.RFC3339, "2024-08-28T15:00:00Z")
	events  = []model.DBEvent{
		{
			Title:       "1st event",
			StartTime:   tm,
			EndTime:     tmE,
			Description: "long description for 1st event",
			UserID:      uuid.New().String(),
		},
		{
			Title:       "2st event",
			StartTime:   tm2,
			EndTime:     tmE2,
			Description: "long description for 2nd event",
			UserID:      uuid.New().String(),
		},
		{
			Title:       "3st event",
			StartTime:   tm3,
			EndTime:     tmE3,
			Description: "long description for 3rd event",
			UserID:      uuid.New().String(),
		},
	}
)

func TestStorage(t *testing.T) {
	t.Run("concurrent insert", func(t *testing.T) {
		var wg sync.WaitGroup
		storage := New()
		for _, e := range events {
			wg.Add(1)
			go func(wg *sync.WaitGroup, e model.DBEvent) {
				defer wg.Done()
				o, err := storage.Insert(ctx, &e)
				require.Nil(t, err)
				require.NotEqual(t, "", o.ID)
			}(&wg, e)
		}
		wg.Wait()
		require.Equal(t, 3, len(events))
		require.Equal(t, 3, len(storage.db))
	})
	t.Run("concurrent update", func(t *testing.T) {
		var wg sync.WaitGroup
		storage := New()
		for _, e := range events {
			wg.Add(1)
			go func(wg *sync.WaitGroup, e model.DBEvent) {
				defer wg.Done()
				_, _ = storage.Insert(ctx, &e)
			}(&wg, e)
		}
		wg.Wait()

		for _, e := range storage.db {
			wg.Add(1)
			go func(wg *sync.WaitGroup, e *model.DBEvent) {
				defer wg.Done()
				e.Title = "new short title"
				o, err := storage.Update(ctx, e)
				require.Nil(t, err)
				require.Equal(t, "new short title", o.Title)
			}(&wg, e)
		}
		wg.Wait()
	})

	t.Run("concurrent delete", func(t *testing.T) {
		var wg sync.WaitGroup
		storage := New()
		for _, e := range events {
			wg.Add(1)
			go func(wg *sync.WaitGroup, e model.DBEvent) {
				defer wg.Done()
				_, _ = storage.Insert(ctx, &e)
			}(&wg, e)
		}
		wg.Wait()

		tmp := make(map[string]*model.DBEvent)
		for k, v := range storage.db {
			tmp[k] = v
		}

		for _, e := range tmp {
			wg.Add(1)
			go func(wg *sync.WaitGroup, e *model.DBEvent) {
				defer wg.Done()
				e.Title = "new short title"
				require.Nil(t, storage.Remove(ctx, e.ID))
			}(&wg, e)
		}
		wg.Wait()
		require.Equal(t, 0, len(storage.db))
	})

	t.Run("find all", func(t *testing.T) {
		var wg sync.WaitGroup
		storage := New()
		for _, e := range events {
			wg.Add(1)
			go func(wg *sync.WaitGroup, e model.DBEvent) {
				defer wg.Done()
				_, _ = storage.Insert(ctx, &e)
			}(&wg, e)
		}
		wg.Wait()

		evnts, err := storage.FindAll(ctx)
		require.Nil(t, err)
		require.Equal(t, 3, len(evnts))
	})

	t.Run("find by id", func(t *testing.T) {
		var wg sync.WaitGroup
		var id atomic.Pointer[string]
		storage := New()
		for _, e := range events {
			wg.Add(1)
			go func(wg *sync.WaitGroup, e model.DBEvent) {
				defer wg.Done()
				o, _ := storage.Insert(ctx, &e)
				id.Store(&o.ID)
			}(&wg, e)
		}
		wg.Wait()

		e, err := storage.Find(ctx, *id.Load())
		require.Nil(t, err)
		require.Equal(t, *id.Load(), e.ID)
	})

	t.Run("insert with non null id", func(t *testing.T) {
		var wg sync.WaitGroup
		storage := New()
		for _, e := range events {
			wg.Add(1)
			go func(wg *sync.WaitGroup, e model.DBEvent) {
				defer wg.Done()
				_, _ = storage.Insert(ctx, &e)
			}(&wg, e)
		}
		wg.Wait()

		tmp := tm.Add(15 * time.Minute)
		newEvent := model.DBEvent{
			ID:          uuid.New().String(),
			Title:       "4th event",
			StartTime:   tmp,
			EndTime:     tm.Add(time.Hour),
			Description: "long description for 4th event",
			UserID:      uuid.New().String(),
		}
		_, err := storage.Insert(ctx, &newEvent)
		require.Error(t, err, repository.ErrIDExists)
	})

	t.Run("insert with already occupied time range", func(t *testing.T) {
		var wg sync.WaitGroup
		storage := New()
		for _, e := range events {
			wg.Add(1)
			go func(wg *sync.WaitGroup, e model.DBEvent) {
				defer wg.Done()
				_, _ = storage.Insert(ctx, &e)
			}(&wg, e)
		}
		wg.Wait()

		tmp := tm.Add(15 * time.Minute)
		newEvent := model.DBEvent{
			Title:       "4th event",
			StartTime:   tmp,
			EndTime:     tm.Add(time.Hour),
			Description: "long description for 4th event",
			UserID:      uuid.New().String(),
		}
		_, err := storage.Insert(ctx, &newEvent)
		require.Error(t, err, repository.ErrDateBusy)
	})

	t.Run("update with already occupied time range", func(t *testing.T) {
		var ev *model.DBEvent
		tm4, _ := time.Parse(time.RFC3339, "2024-08-28T12:59:00Z")
		storage := New()
		for _, e := range events {
			ev, _ = storage.Insert(ctx, &e)
		}
		ev.StartTime = tm4
		ev.EndTime, _ = time.Parse(time.RFC3339, "2024-08-28T14:01:00Z")

		_, err := storage.Update(ctx, ev)
		require.Error(t, err, repository.ErrDateBusy)
	})
}
