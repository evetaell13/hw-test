package app

import (
	"context"
	"time"

	storage "github.com/evetaell13/hw-test/hw12_13_14_15_calendar/internal/storage"
)

type App struct { // TODO
}

type Logger interface { // TODO
}

type Storage interface { // TODO
	WriteEvent(storage.Event) error
	WriteUser(storage.User) error
	UpdateEvent(storage.Event) error
	UpdateUser(storage.User) error
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
	GetEventsByPeriod(start time.Time, finish time.Time) ([]storage.Event, error)
	GetEvent(id string) (storage.Event, error)
	GetUser(id string) (storage.User, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	return nil
}

// TODO
