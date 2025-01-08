package internalhttp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/config"
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

type calendarHandler struct {}

func (h *calendarHandler)hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello, world"))
}

func NewServer(logger Logger, app Application, httpConf config.HttpConf) *Server {
	handler := &calendarHandler{}
	mux := http.NewServeMux()
	addr := fmt.Sprintf("%s:%d", httpConf.Host, httpConf.Port)

	mux.HandleFunc("/hello", handler.hello)

	return &Server{
		Server: http.Server{
			Addr:         addr,
			Handler:      loggingMiddleware(mux),
		},
		logger: logger,
	}
}

func (s *Server) Start(ctx context.Context) error {
	return s.Server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}