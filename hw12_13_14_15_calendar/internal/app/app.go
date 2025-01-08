package app

import (
	"context"
	"strings"
	"time"

	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/config"                       //nolint:depguard
	storagelib "github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/storage"           //nolint:depguard
	memorystorage "github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/storage/memory" //nolint:depguard
	sqlstorage "github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/storage/sql"       //nolint:depguard
)

const (
	SQLStorage      = "sql"
	InMemoryStorage = "in-memory"
)

type App struct { // TODO
}

type Logger interface { // TODO
}

type Storage interface {
	Create(ctx context.Context, event storagelib.Event) (int, error)
	Update(ctx context.Context, id int, event storagelib.Event) error
	Delete(ctx context.Context, id int) error
	FindForDay(ctx context.Context, date time.Time) ([]storagelib.Event, error)
	FindForWeek(ctx context.Context, date time.Time) ([]storagelib.Event, error)
	FindForMonth(ctx context.Context, date time.Time) ([]storagelib.Event, error)
	FindBetweenDates(ctx context.Context, start time.Time, end time.Time) ([]storagelib.Event, error)
}

func NewStorage(config config.Config) Storage {
	switch strings.ToLower(config.Storage) {
	case SQLStorage:
		return sqlstorage.New(config.DB)
	case InMemoryStorage:
		return memorystorage.New()
	default:
		panic("Unknown storage")
	}
}

func New(_ Logger, _ Storage) *App {
	return &App{}
}

func (a *App) CreateEvent(_ context.Context, _ int, _ string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
