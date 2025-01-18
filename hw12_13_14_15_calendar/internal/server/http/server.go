package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/config" //nolint:depguard
	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/storage"
)

type Server struct {
	http.Server
	logger Logger
	handler *calendarHandler
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

func NewServer(logger Logger, app Application, httpConf config.HTTPConf) *Server {
	handler := &calendarHandler{
		app: app,
	}
	mux := http.NewServeMux()
	addr := fmt.Sprintf("%s:%d", httpConf.Host, httpConf.Port)

	mux.HandleFunc("/hello", handler.hello)
	mux.HandleFunc("POST /event", handler.create)
	mux.HandleFunc("PUT /event/{id}", handler.update)
	mux.HandleFunc("DELETE /event/{id}", handler.delete)
	mux.HandleFunc("GET /event/{mode}", handler.find)

	return &Server{
		Server: http.Server{
			Addr:              addr,
			Handler:           loggingMiddleware(mux),
			ReadHeaderTimeout: 5 * time.Second,
		},
		logger: logger,
		handler: handler,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.handler.ctx = ctx
	return s.Server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}
