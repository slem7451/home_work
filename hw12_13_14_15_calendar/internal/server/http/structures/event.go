package structures

type Event struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	EventDate  string `json:"event_date"`
	DateSince  string `json:"date_since"`
	Descr      string    `json:"descr"`
	UserID     int       `json:"user_id"`
	NotifyDate string `json:"notify_date"`
}
