package app

import (
	"context"
	"time"

	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/storage/model"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/storage/repository"
)

type App struct { // TODO
	logger  *logger.Logger
	storage repository.EventInterface
}

func New(logger *logger.Logger, storage repository.EventInterface) *App {
	return &App{logger: logger, storage: storage}
}

func (a *App) CreateEvent(ctx context.Context, title string, startTime time.Time,
	endTime time.Time, userID string,
) (*model.Event, error) {
	return a.storage.Insert(ctx, &model.Event{Title: title, StartTime: startTime, EndTime: endTime, UserID: userID})
}

// TODO
