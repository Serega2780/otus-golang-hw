package model

import "time"

type Event struct {
	ID                string        `json:"id"`
	Title             string        `json:"title"`
	StartTime         time.Time     `json:"startTime"`
	EndTime           time.Time     `json:"endTime"`
	Description       string        `json:"description"`
	UserID            string        `json:"userId"`
	NotifyBeforeEvent time.Duration `json:"notifyBeforeEvent"`
}

func NewEvent() *Event {
	return &Event{}
}
