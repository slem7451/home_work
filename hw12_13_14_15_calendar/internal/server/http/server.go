package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/config" //nolint:depguard
	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/server"
)

type Server struct {
	http.Server
	logger server.Logger
	handler *calendarHandler
}

func NewServer(logger server.Logger, app server.Application, httpConf config.HTTPConf) server.Server {
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

func (s *Server) Whoami() string {
	return "HTTP"
}