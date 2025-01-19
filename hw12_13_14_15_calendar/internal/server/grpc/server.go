package internalgrpc

import (
	"context"
	"fmt"
	"net"

	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/config"         //nolint:depguard
	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/server"         //nolint:depguard
	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/server/grpc/pb" //nolint:depguard
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	*grpc.Server
	logger server.Logger
	config config.GRPCConf
}

func NewServer(logger server.Logger, app server.Application, grpcConf config.GRPCConf) server.Server {
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(loggingInterceptor()))
	pb.RegisterEventServiceServer(grpcServer, &calendarService{app: app})
	reflection.Register(grpcServer)

	return &Server{
		Server: grpcServer,
		logger: logger,
		config: grpcConf,
	}
}

func (s *Server) Start(_ context.Context) error {
	lsn, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.config.Host, s.config.Port))
	if err != nil {
		return err
	}

	if err := s.Server.Serve(lsn); err != nil {
		return err
	}

	return nil
}

func (s *Server) Stop(_ context.Context) error {
	s.Server.GracefulStop()
	return nil
}

func (s *Server) Whoami() string {
	return "GRPC"
}
