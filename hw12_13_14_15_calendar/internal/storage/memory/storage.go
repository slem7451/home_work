package memorystorage

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
)

var (
	ErrNotUniqueID   = errors.New("ID must be unique")
	ErrEventNotFound = errors.New("event with this ID not found")
)

type Storage struct {
	mu      sync.RWMutex
	storage map[int]storage.Event
	eventID int
}

func New() *Storage {
	return &Storage{storage: make(map[int]storage.Event), eventID: 1}
}

func (s *Storage) Create(_ context.Context, event storage.Event) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if event.ID == 0 {
		event.ID = s.eventID
		s.eventID++
	}

	if _, ok := s.storage[event.ID]; ok {
		return 0, ErrNotUniqueID
	}

	s.storage[event.ID] = event
	return event.ID, nil
}

func (s *Storage) Update(_ context.Context, id int, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.storage[id]; !ok {
		return ErrEventNotFound
	}

	delete(s.storage, id)

	s.storage[event.ID] = event
	return nil
}

func (s *Storage) Delete(_ context.Context, id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.storage, id)
	return nil
}

func (s *Storage) FindForDay(ctx context.Context, date time.Time) ([]storage.Event, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(time.Hour * 24).Add(-time.Nanosecond)

	return s.FindBetweenDates(ctx, start, end)
}

func (s *Storage) FindForWeek(ctx context.Context, date time.Time) ([]storage.Event, error) {
	start := date.Add(-time.Hour * 24 * time.Duration(date.Weekday())).Add(time.Hour * 24)
	end := start.Add(time.Hour * 24 * 7).Add(-time.Nanosecond)

	return s.FindBetweenDates(ctx, start, end)
}

func (s *Storage) FindForMonth(ctx context.Context, date time.Time) ([]storage.Event, error) {
	start := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)

	return s.FindBetweenDates(ctx, start, end)
}

func (s *Storage) FindBetweenDates(_ context.Context, start time.Time, end time.Time) ([]storage.Event, error) {
	res := make([]storage.Event, 0)

	for _, v := range s.storage {
		if v.EventDate.Compare(start) >= 0 && v.EventDate.Compare(end) <= 0 {
			res = append(res, v)
		}
	}

	return res, nil
}
