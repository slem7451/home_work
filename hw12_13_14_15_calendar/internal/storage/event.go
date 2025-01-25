package storage

import "time"

type Event struct {
	ID         int       `db:"id"`
	Title      string    `db:"title"`
	EventDate  time.Time `db:"event_date"`
	DateSince  time.Time `db:"date_since"`
	Descr      string    `db:"descr"`
	UserID     int       `db:"user_id"`
	NotifyDate time.Time `db:"notify_date"`
	IsSended   bool      `db:"is_sended"`
}

type Notification struct {
	ID        int       `db:"id"`
	Title     string    `db:"title"`
	EventDate time.Time `db:"event_date"`
	UserID    int       `db:"user_id"`
}
