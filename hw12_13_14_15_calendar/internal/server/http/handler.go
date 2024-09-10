package internalhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/model"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/service"
	"github.com/gorilla/mux"
)

const (
	ErrorResponseFieldKey              = "error"
	ErrorMandatoryFieldsMissingOrWrong = "mandatory fields are missing or wrong error"
	EventID                            = "eventID"
	Date                               = "date"
)

type EventsHandler struct {
	ctx     context.Context
	log     *logger.Logger
	service service.EventService
}

func NewEventsHandler(ctx context.Context, logger *logger.Logger, eventService service.EventService) *EventsHandler {
	return &EventsHandler{ctx: ctx, log: logger, service: eventService}
}

func (s *EventsHandler) HelloWorld(w http.ResponseWriter, r *http.Request) {
	resp := "Hello, World!"
	if r.Method != http.MethodGet {
		s.writeMethodNotAllowedResponse(w, r)
		return
	}
	*r = *(r.WithContext(context.WithValue(r.Context(), STATUS, http.StatusOK)))
	s.respondWithJSON(w, r, http.StatusOK, resp)
}

func (s *EventsHandler) InsertEvent(w http.ResponseWriter, r *http.Request) {
	var event *model.Event
	now := time.Now()
	decoder := json.NewDecoder(r.Body)
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		s.writeMethodNotAllowedResponse(w, r)
		return
	}
	if err := decoder.Decode(&event); err != nil {
		s.respondWithError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if len(event.Title) == 0 || event.StartTime.Before(now.Truncate(time.Hour)) ||
		event.EndTime.Before(now.Truncate(time.Hour)) || len(event.UserID) == 0 {
		s.respondWithError(w, r, http.StatusBadRequest, ErrorMandatoryFieldsMissingOrWrong)
		return
	}
	rEvent, err := s.service.InsertEvent(s.ctx, event)
	if err != nil {
		s.respondWithError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	s.respondWithJSON(w, r, http.StatusOK, rEvent)
}

func (s *EventsHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	var event *model.Event
	now := time.Now()
	decoder := json.NewDecoder(r.Body)
	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		s.writeMethodNotAllowedResponse(w, r)
		return
	}
	if err := decoder.Decode(&event); err != nil {
		s.respondWithError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	if len(event.Title) == 0 || event.StartTime.Before(now.Truncate(time.Hour)) ||
		event.EndTime.Before(now.Truncate(time.Hour)) || len(event.UserID) == 0 {
		s.respondWithError(w, r, http.StatusBadRequest, ErrorMandatoryFieldsMissingOrWrong)
		return
	}
	rEvent, err := s.service.UpdateEvent(s.ctx, event)
	if err != nil {
		s.respondWithError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	*r = *(r.WithContext(context.WithValue(r.Context(), STATUS, http.StatusOK)))
	s.respondWithJSON(w, r, http.StatusOK, rEvent)
}

func (s *EventsHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete && r.Method != http.MethodGet {
		s.writeMethodNotAllowedResponse(w, r)
		return
	}
	muxVar := mux.Vars(r)
	id := muxVar["id"]
	if len(id) == 0 {
		s.respondWithError(w, r, http.StatusBadRequest, "path variable id is missing")
		return
	}
	err := s.service.RemoveEvent(s.ctx, id)
	if err != nil {
		s.respondWithError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	s.respondWithJSON(w, r, http.StatusOK, map[string]string{EventID: id})
}

func (s *EventsHandler) FindEventByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete && r.Method != http.MethodGet {
		s.writeMethodNotAllowedResponse(w, r)
		return
	}
	muxVar := mux.Vars(r)
	id := muxVar["id"]
	if len(id) == 0 {
		s.respondWithError(w, r, http.StatusBadRequest, "path variable id is missing")
		return
	}
	event, err := s.service.FindEvent(s.ctx, id)
	if err != nil {
		s.respondWithError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	if event == nil {
		s.respondWithJSON(w, r, http.StatusOK, "")
		return
	}
	s.respondWithJSON(w, r, http.StatusOK, event)
}

func (s *EventsHandler) FindAllEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.writeMethodNotAllowedResponse(w, r)
		return
	}
	events, err := s.service.FindAllEvents(s.ctx)
	if err != nil {
		s.respondWithError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	s.respondWithJSON(w, r, http.StatusOK, events)
}

func (s *EventsHandler) FindAllEventByDay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.writeMethodNotAllowedResponse(w, r)
		return
	}
	muxVar := mux.Vars(r)
	date := muxVar[Date]
	if len(date) == 0 {
		s.respondWithError(w, r, http.StatusBadRequest, "path variable date is missing")
		return
	}
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		s.respondWithError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	events, err := s.service.FindAllEventsByDay(s.ctx, t)
	if err != nil {
		s.respondWithError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	s.respondWithJSON(w, r, http.StatusOK, events)
}

func (s *EventsHandler) FindAllEventByWeek(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.writeMethodNotAllowedResponse(w, r)
		return
	}
	muxVar := mux.Vars(r)
	date := muxVar[Date]
	if len(date) == 0 {
		s.respondWithError(w, r, http.StatusBadRequest, "path variable date is missing")
		return
	}
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		s.respondWithError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	events, err := s.service.FindAllEventsByWeek(s.ctx, t)
	if err != nil {
		s.respondWithError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	s.respondWithJSON(w, r, http.StatusOK, events)
}

func (s *EventsHandler) FindAllEventByMonth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.writeMethodNotAllowedResponse(w, r)
		return
	}
	muxVar := mux.Vars(r)
	date := muxVar[Date]
	if len(date) == 0 {
		s.respondWithError(w, r, http.StatusBadRequest, "path variable date is missing")
		return
	}
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		s.respondWithError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	events, err := s.service.FindAllEventsByMonth(s.ctx, t)
	if err != nil {
		s.respondWithError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	s.respondWithJSON(w, r, http.StatusOK, events)
}

func (s *EventsHandler) writeMethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	resp := fmt.Sprintf("method %s not not supported on uri %s", r.Method, r.URL.Path)
	*r = *(r.WithContext(context.WithValue(r.Context(), STATUS, http.StatusMethodNotAllowed)))
	s.respondWithError(w, r, http.StatusMethodNotAllowed, resp)
}

func (s *EventsHandler) respondWithError(w http.ResponseWriter, r *http.Request, code int, message string) {
	status := r.Context().Value(STATUS)
	if status == nil {
		*r = *(r.WithContext(context.WithValue(r.Context(), STATUS, code)))
	}
	s.respondWithJSON(w, r, code, map[string]string{ErrorResponseFieldKey: message})
}

func (s *EventsHandler) respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	status := r.Context().Value(STATUS)
	if status == nil {
		*r = *(r.WithContext(context.WithValue(r.Context(), STATUS, code)))
	}
	response, err := json.Marshal(payload)
	if err != nil {
		s.log.Errorf("response marshal error: %v", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(response)
	if err != nil {
		s.log.Errorf("response write error: %v", err)
	}
}
