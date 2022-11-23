package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpcrest/gen/pbgen"
	"grpcrest/pkg/config"
	"grpcrest/pkg/logger"
	"grpcrest/pkg/service"
	"net"
)

type GRPC struct {
	server  *grpc.Server
	config  config.Server
	logger  logger.Logger
	service service.Service
}

func New(cfg config.Server, lgr logger.Logger, ser service.Service) (*GRPC, error) {
	srv := grpc.NewServer()

	pbgen.RegisterUserServiceServer(srv, ser)

	return &GRPC{
		server:  srv,
		config:  cfg,
		logger:  lgr,
		service: ser,
	}, nil
}

func (g *GRPC) Serve() error {
	lis, err := net.Listen("tcp", g.config.Address())

	if err != nil {
		return fmt.Errorf("failed to listen on tcp address %s: %w", g.config.Address(), err)
	}

	err = g.server.Serve(lis)

	if err != grpc.ErrServerStopped {
		err = fmt.Errorf("failed to start grpc server: %w", err)
	} else {
		err = nil
	}

	return err
}

func (g *GRPC) Shutdown(ctx context.Context) error {
	g.server.GracefulStop()

	return nil
}
