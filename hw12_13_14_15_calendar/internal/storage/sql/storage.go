package sqlstorage

import (
	"context"
	"time"

	_ "github.com/jackc/pgx/stdlib" //nolint:depguard
	"github.com/jmoiron/sqlx"       //nolint:depguard
	"github.com/pressly/goose/v3" //nolint:depguard
	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/config"  //nolint:depguard
	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
	"github.com/slem7451/home_work/hw12_13_14_15_calendar/migrations" //nolint:depguard
)

type Storage struct {
	db   *sqlx.DB
	conn *sqlx.Conn
}

func New(dbConf config.DBConf) *Storage {
	db, err := sqlx.Open("pgx", dbConf.DSN())
	if err != nil {
		panic(err)
	}

	goose.SetBaseFS(migrations.EmbedMigrations)
	
	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(db.DB, "."); err != nil {
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

	stmt, err := s.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return 0, err
	}

	var id int
	err = stmt.Get(&id, event)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Storage) Update(ctx context.Context, id int, event storage.Event) error {
	query := `update events set 
				title = :title, 
				event_date = :event_date, 
				date_since = :date_since, 
				descr = :descr, 
				user_id = :user_id, 
				notify_date = :notify_date
				where id = :id`
	event.ID = id

	if _, err := s.db.NamedExecContext(ctx, query, event); err != nil {
		return err
	}

	return nil
}

func (s *Storage) Delete(ctx context.Context, id int) error {
	query := `delete from events where id = $1`

	if _, err := s.db.ExecContext(ctx, query, id); err != nil {
		return err
	}

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

func (s *Storage) FindBetweenDates(ctx context.Context, start time.Time, end time.Time) ([]storage.Event, error) {
	rows, err := s.db.QueryxContext(ctx, `select * from events where event_date between $1 and $2`, start, end)
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

func (s *Storage) FindEventsForNotify(ctx context.Context) ([]storage.Notification, error) {
	var emptyTime time.Time

	rows, err := s.db.QueryxContext(ctx,
		`select id, title, event_date, user_id from events 
		where is_sended = false 
		and ((event_date < now() and notify_date = $1) or (notify_date <> $1 and notify_date < now()))`,
		emptyTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notifications := make([]storage.Notification, 0)

	for rows.Next() {
		var n storage.Notification

		if err := rows.StructScan(&n); err != nil {
			return make([]storage.Notification, 0), err
		}

		notifications = append(notifications, n)
	}

	return notifications, nil
}

func (s *Storage) RemoveOldEvents(ctx context.Context) error {
	query := `delete from events where is_sended = true and event_date < now() - interval '1 year'`

	if _, err := s.db.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}

func (s *Storage) MarkSendedEvent(ctx context.Context, id int) error {
	query := `update events set is_sended = true where id = $1`

	if _, err := s.db.ExecContext(ctx, query, id); err != nil {
		return err
	}

	return nil
}
