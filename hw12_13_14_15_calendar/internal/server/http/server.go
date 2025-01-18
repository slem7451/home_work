package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/config" //nolint:depguard
)

type Server struct {
	http.Server
	logger Logger
}

type Logger interface {
	Error(msg string)
	Warn(msg string)
	Info(msg string)
	Debug(msg string)
}

type Application interface { // TODO
}

type calendarHandler struct{}

func (h *calendarHandler) hello(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("hello, world"))
}

func NewServer(logger Logger, _ Application, httpConf config.HTTPConf) *Server {
	handler := &calendarHandler{}
	mux := http.NewServeMux()
	addr := fmt.Sprintf("%s:%d", httpConf.Host, httpConf.Port)

	mux.HandleFunc("/hello", handler.hello)

	return &Server{
		Server: http.Server{
			Addr:              addr,
			Handler:           loggingMiddleware(mux),
			ReadHeaderTimeout: 5 * time.Second,
		},
		logger: logger,
	}
}

func (s *Server) Start(_ context.Context) error {
	return s.Server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}
