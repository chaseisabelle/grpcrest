package server

import (
	"fmt"
	"grpcrest/pkg/config"
	"grpcrest/pkg/logger"
	"grpcrest/pkg/server/grpc"
	"grpcrest/pkg/server/rest"
	"grpcrest/pkg/service"
)

type Server interface {
	Serve() error
}

type server struct {
	grpc Server
	rest Server
}

func New(cfg config.Config, lgr logger.Logger, ser service.Service) (Server, error) {
	rpc, err := grpc.New(cfg.GRPC())

	if err != nil {
		return nil, fmt.Errorf("failed to initialize grpc server: %w", err)
	}

	rst, err := rest.New(cfg.REST())

	if err != nil {
		return nil, fmt.Errorf("failed to initialize rest server: %w", err)
	}

	return &server{
		grpc: rpc,
		rest: rst,
	}, nil
}

func (s *server) Serve() error {
	ech := make(chan error)

	go func(chan error) {
		ech <- s.grpc.Serve()
	}(ech)

	go func(chan error) {
		ech <- s.rest.Serve()
	}(ech)

}
