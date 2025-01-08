package app

import (
	"context"
	"strings"
	"time"

	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/config"
	storagelib "github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/storage/sql"
)

const (
	SqlStorage = "sql"
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
	case SqlStorage:
		return sqlstorage.New(config.Db)
	case InMemoryStorage:
		return memorystorage.New()
	default:
		panic("Unknown storage")
	}
}

func New(logger Logger, storage Storage) *App {
	return &App{}
}

func (a *App) CreateEvent(ctx context.Context, id int, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
