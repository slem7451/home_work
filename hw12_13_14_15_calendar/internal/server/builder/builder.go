package builder

import (
	calendarconfig "github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/config/calendar"
	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/server"                   //nolint:depguard
	internalgrpc "github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/server/grpc" //nolint:depguard
	internalhttp "github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/server/http" //nolint:depguard
)

func NewServers(logger server.Logger, app server.Application, config calendarconfig.Config) []server.Server {
	grpc := internalgrpc.NewServer(logger, app, config.GRPC)
	http := internalhttp.NewServer(logger, app, config.HTTP)

	return []server.Server{grpc, http}
}
