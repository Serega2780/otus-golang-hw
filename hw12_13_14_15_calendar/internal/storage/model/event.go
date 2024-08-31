package model

import "time"

type Event struct {
	ID                string
	Title             string
	StartTime         time.Time
	EndTime           time.Time
	Description       string
	UserID            string
	NotifyBeforeEvent time.Duration
}

func New() *Event {
	return &Event{}
}
