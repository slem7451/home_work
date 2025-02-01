package app

import (
	"context"
	"fmt"
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

type App struct {
	storage Storage
	logger  Logger
}

type Logger interface {
	Error(msg string)
	Warn(msg string)
	Info(msg string)
	Debug(msg string)
}

type Storage interface {
	Create(ctx context.Context, event storagelib.Event) (int, error)
	Update(ctx context.Context, id int, event storagelib.Event) error
	Delete(ctx context.Context, id int) error
	FindForDay(ctx context.Context, date time.Time) ([]storagelib.Event, error)
	FindForWeek(ctx context.Context, date time.Time) ([]storagelib.Event, error)
	FindForMonth(ctx context.Context, date time.Time) ([]storagelib.Event, error)
	FindBetweenDates(ctx context.Context, start time.Time, end time.Time) ([]storagelib.Event, error)

	FindEventsForNotify(ctx context.Context) ([]storagelib.Notification, error)
	RemoveOldEvents(ctx context.Context) error
	MarkSendedEvent(ctx context.Context, id int) error
}

func NewStorage(config config.Config) Storage {
	switch strings.ToLower(config.GetStorage()) {
	case SQLStorage:
		return sqlstorage.New(config.GetDB())
	case InMemoryStorage:
		return memorystorage.New()
	default:
		panic("Unknown storage")
	}
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, event storagelib.Event) (int, error) {
	a.logger.Info("Creating event")

	id, err := a.storage.Create(ctx, event)
	if err != nil {
		a.logger.Error(err.Error())
	}

	return id, err
}

func (a *App) UpdateEvent(ctx context.Context, id int, event storagelib.Event) error {
	a.logger.Info(fmt.Sprintf("Updating event with ID %d", id))

	err := a.storage.Update(ctx, id, event)
	if err != nil {
		a.logger.Error(err.Error())
	}

	return err
}

func (a *App) DeleteEvent(ctx context.Context, id int) error {
	a.logger.Info(fmt.Sprintf("Deleting event with ID %d", id))

	err := a.storage.Delete(ctx, id)
	if err != nil {
		a.logger.Error(err.Error())
	}

	return err
}

func (a *App) FindEventsForDay(ctx context.Context, date time.Time) ([]storagelib.Event, error) {
	a.logger.Info(fmt.Sprintf("Finding events for day %s", date.Format(time.DateOnly)))

	events, err := a.storage.FindForDay(ctx, date)
	if err != nil {
		a.logger.Error(err.Error())
	}

	return events, err
}

func (a *App) FindEventsForWeek(ctx context.Context, date time.Time) ([]storagelib.Event, error) {
	a.logger.Info(fmt.Sprintf("Finding events for week %s", date.Format(time.DateOnly)))

	events, err := a.storage.FindForWeek(ctx, date)
	if err != nil {
		a.logger.Error(err.Error())
	}

	return events, err
}

func (a *App) FindEventsForMonth(ctx context.Context, date time.Time) ([]storagelib.Event, error) {
	a.logger.Info(fmt.Sprintf("Finding events for month %s", date.Format(time.DateOnly)))

	events, err := a.storage.FindForMonth(ctx, date)
	if err != nil {
		a.logger.Error(err.Error())
	}

	return events, err
}

func (a *App) FindEventsBetweenDates(ctx context.Context, start time.Time, end time.Time) ([]storagelib.Event, error) {
	a.logger.Info(fmt.Sprintf("Finding events between %s and %s", start.Format(time.DateOnly), end.Format(time.DateOnly)))

	events, err := a.storage.FindBetweenDates(ctx, start, end)
	if err != nil {
		a.logger.Error(err.Error())
	}

	return events, err
}
