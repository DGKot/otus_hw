package app

import (
	"context"
)

type App struct {
	storage Storage
	logger  Logger
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type Storage interface { // TODO
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(_ context.Context, _, _ string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}
