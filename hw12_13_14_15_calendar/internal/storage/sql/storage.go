package sqlstorage

import (
	"context"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	db   *sqlx.DB
	conn *sqlx.Conn
}

func New(dbConf config.DbConf) *Storage {
	db, err := sqlx.Open("pgx", fmt.Sprintf("postgres://%s:%s@%s:%d/%s", dbConf.User, dbConf.Password, dbConf.Host, dbConf.Port, dbConf.Name))
	if err != nil {
		panic(err)
	}

	return &Storage{db: db}
}

func (s *Storage) Connect(ctx context.Context) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}

	s.conn = conn
	return nil
}

func (s *Storage) Close(_ context.Context) error {
	return s.conn.Close()
}

func (s *Storage) Create(ctx context.Context, event storage.Event) (int, error) {
	var query string

	if event.ID == 0 {
		query = `insert into events (title, event_date, date_since, descr, user_id, notify_date) 
					values (:title, :event_date, :date_since, :descr, :user_id, :notify_date) returning id`
	} else {
		query = `insert into events (id, title, event_date, date_since, descr, user_id, notify_date) 
					values (:id, :title, :event_date, :date_since, :descr, :user_id, :notify_date) returning id`
	}

	res, err := s.db.NamedExecContext(ctx, query, event)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s *Storage) Update(ctx context.Context, id int, event storage.Event) error {
	query := `update events set title = :title, event_date = :event_date, date_since = :date_since, descr = :descr, user_id = :user_id, notify_date = :notify_date
				where id = :id`
	event.ID = id

	if _, err := s.db.NamedExecContext(ctx, query, event); err != nil {
		return err
	}

	return nil
}

func (s *Storage) Delete(ctx context.Context, id int) error {
	query := `delete from event where id = $1`

	if _, err := s.db.ExecContext(ctx, query, id); err != nil {
		return err
	}
	
	return nil
}

func (s *Storage) FindForDay(ctx context.Context, date time.Time) ([]storage.Event, error) {
	start := date.Truncate(time.Hour * 24)
	end := start.Add(time.Hour * 24).Add(-time.Nanosecond)

	return s.FindBetweenDates(ctx, start, end)
}

func (s *Storage) FindForWeek(ctx context.Context, date time.Time) ([]storage.Event, error) {
	start := date.Truncate(time.Hour * 24 * 7)
	end := start.Add(time.Hour * 24 * 7).Add(-time.Nanosecond)

	return s.FindBetweenDates(ctx, start, end)
}

func (s *Storage) FindForMonth(ctx context.Context, date time.Time) ([]storage.Event, error) {
	start := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)

	return s.FindBetweenDates(ctx, start, end)
}

func (s *Storage) FindBetweenDates(ctx context.Context, start time.Time, end time.Time) ([]storage.Event, error) {
	rows, err := s.db.NamedQueryContext(ctx, `select * from event where event_date between :start and :end`, map[string]time.Time{"start": start, "end": end})

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]storage.Event, 0)

	for rows.Next() {
		var e storage.Event

		if err := rows.StructScan(&e); err != nil {
			return make([]storage.Event, 0), err
		}

		events = append(events, e)
	}

	return events, nil
}