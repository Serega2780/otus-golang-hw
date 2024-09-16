package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/model"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/service/memory"
	service "github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func Test_Events(t *testing.T) {
	ctx := context.Background()
	l := logger.New(&config.LoggerConf{Level: "info", Format: "json", LogToFile: false, LogToConsole: true})
	repo := service.New()
	service := memory.NewEventMemoryService(repo)
	h := NewEventsHandler(ctx, l, service)
	r := mux.NewRouter()
	r.HandleFunc("/v1/events", h.InsertEvent).Methods(http.MethodPost)
	r.HandleFunc("/v1/events", h.UpdateEvent).Methods(http.MethodPut)
	r.HandleFunc("/v1/events", h.FindAllEvent).Methods(http.MethodGet)
	r.HandleFunc("/v1/events/{id}", h.FindEventByID).Methods(http.MethodGet)
	r.HandleFunc("/v1/events/{id}", h.DeleteEvent).Methods(http.MethodDelete)
	ts := httptest.NewServer(r)
	defer ts.Close()
	t.Run("not found", func(t *testing.T) {
		res, err := http.Get(ts.URL + "/v1/events/c39d8763-34ad-4759-954f-b288649dac73") //nolint:noctx, bodyclose
		require.NoError(t, err)
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(res.Body)
		bodyBytes, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)
		require.Equal(t, "\"\"", string(bodyBytes))
	})
	t.Run("insert&update", func(t *testing.T) {
		var event model.Event
		var uEvent model.Event
		e, _ := json.Marshal(getEvent())
		b := bytes.NewReader(e)
		res, err := http.Post(ts.URL+"/v1/events", "application/json", b) //nolint:noctx, bodyclose
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(res.Body)
		require.Equal(t, http.StatusOK, res.StatusCode)
		require.NoError(t, err)
		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(&event)
		require.NoError(t, err)
		require.Greater(t, len(event.ID), 0)

		event.Title = "UPDATED"
		res, err = updateEvent(ctx, ts, &event) //nolint:bodyclose
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				fmt.Printf("%v\n", err)
			}
		}(res.Body)
		require.Equal(t, http.StatusOK, res.StatusCode)
		require.NoError(t, err)
		decoder = json.NewDecoder(res.Body)
		err = decoder.Decode(&uEvent)
		require.NoError(t, err)
		require.Equal(t, "UPDATED", uEvent.Title)

		clearEvents(ctx, ts, []*model.Event{&uEvent})
	})
	t.Run("findAll", func(t *testing.T) {
		e := getEvent()
		e2 := getEvent()
		e2.StartTime = e2.StartTime.Add(2 * time.Hour)
		e2.EndTime = e2.EndTime.Add(3 * time.Hour)
		insertEvent(ts, &e)
		insertEvent(ts, &e2)
		var events []*model.Event
		res, err := http.Get(ts.URL + "/v1/events") //nolint:noctx
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(res.Body)
		require.Equal(t, http.StatusOK, res.StatusCode)
		require.NoError(t, err)
		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(&events)
		res.Body.Close()
		require.NoError(t, err)
		require.Equal(t, 2, len(events))
		clearEvents(ctx, ts, events)
	})
}

func getEvent() model.Event {
	return model.Event{
		Title:       "4th event",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(1 * time.Hour),
		Description: "long description for 4th event",
		UserID:      uuid.New().String(),
	}
}

func clearEvents(ctx context.Context, ts *httptest.Server, events []*model.Event) {
	client := &http.Client{}
	for _, e := range events {
		req, _ := http.NewRequestWithContext(ctx, "DELETE", ts.URL+"/v1/events/"+e.ID, nil)
		_, _ = client.Do(req) //nolint:bodyclose
	}
}

func insertEvent(ts *httptest.Server, event *model.Event) {
	e, _ := json.Marshal(event)
	_, _ = http.Post(ts.URL+"/v1/events", "application/json", bytes.NewReader(e)) //nolint:noctx,bodyclose
}

func updateEvent(ctx context.Context, ts *httptest.Server, event *model.Event) (*http.Response, error) {
	e, _ := json.Marshal(event)
	b := bytes.NewReader(e)
	client := &http.Client{}
	req, _ := http.NewRequestWithContext(ctx, "PUT", ts.URL+"/v1/events", b)
	return client.Do(req)
}
