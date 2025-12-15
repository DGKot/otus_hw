package storage

import (
	"errors"
	"time"
)

type Event struct {
	ID                       string
	Title                    string
	Datetime                 time.Time
	Duration                 time.Duration
	Description              string
	UserID                   string
	NotificationTimeDuration time.Duration
}

var (
	ErrDateBusy      = errors.New("date and time already busy")
	ErrEventNotFound = errors.New("event not found")
)
