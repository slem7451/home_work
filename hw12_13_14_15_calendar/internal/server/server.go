package server

import (
	"context"
	"time"

	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
)

type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Whoami() string
}

type Logger interface {
	Error(msg string)
	Warn(msg string)
	Info(msg string)
	Debug(msg string)
}

type Application interface {
	CreateEvent(ctx context.Context, event storage.Event) (int, error)
	UpdateEvent(ctx context.Context, id int, event storage.Event) error
	DeleteEvent(ctx context.Context, id int) error
	FindEventsForDay(ctx context.Context, date time.Time) ([]storage.Event, error)
	FindEventsForWeek(ctx context.Context, date time.Time) ([]storage.Event, error)
	FindEventsForMonth(ctx context.Context, date time.Time) ([]storage.Event, error)
	FindEventsBetweenDates(ctx context.Context, start time.Time, end time.Time) ([]storage.Event, error)
}
