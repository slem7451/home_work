package structures

type Event struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	EventDate  string `json:"event_date"` //nolint:tagliatelle
	DateSince  string `json:"date_since"` //nolint:tagliatelle
	Descr      string `json:"descr"`
	UserID     int    `json:"user_id"`     //nolint:tagliatelle
	NotifyDate string `json:"notify_date"` //nolint:tagliatelle
}
