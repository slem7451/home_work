package builder

import (
	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/server"
	internalgrpc "github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/server/http"
)

func NewServers(logger server.Logger, app server.Application, config config.Config) []server.Server {
	grpc := internalgrpc.NewServer(logger, app, config.GRPC)
	http := internalhttp.NewServer(logger, app, config.HTTP)

	return []server.Server{grpc, http}
}