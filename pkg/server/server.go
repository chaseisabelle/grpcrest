package server

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"grpcrest/pkg/config"
	"grpcrest/pkg/logger"
	"grpcrest/pkg/server/grpc"
	"grpcrest/pkg/server/rest"
	"grpcrest/pkg/service"
	"os"
	"os/signal"
	"syscall"
)

type Server interface {
	Serve() error
	Shutdown(context.Context) error
}

type server struct {
	grpc   Server
	rest   Server
	logger logger.Logger
}

func New(cfg config.Config, lgr logger.Logger, ser service.Service) (Server, error) {
	rpc, err := grpc.New(cfg.GRPC(), lgr, ser)

	if err != nil {
		return nil, fmt.Errorf("failed to initialize grpc server: %w", err)
	}

	rst, err := rest.New(cfg.REST(), lgr, ser)

	if err != nil {
		return nil, fmt.Errorf("failed to initialize rest server: %w", err)
	}

	return &server{
		grpc:   rpc,
		rest:   rst,
		logger: lgr,
	}, nil
}

func (s *server) Serve() error {
	ctx, cxl := context.WithCancel(context.Background())

	go func(cxl context.CancelFunc) {
		chn := make(chan os.Signal, 1)

		signal.Notify(chn, os.Interrupt, syscall.SIGTERM)

		<-chn

		cxl()
	}(cxl)

	eg, gCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return s.grpc.Serve()
	})

	eg.Go(func() error {
		<-gCtx.Done()
		return httpServer.Shutdown(context.Background())
	})

	if err := eg.Wait(); err != nil {
		fmt.Printf("exit reason: %s \n", err)
	}
}
