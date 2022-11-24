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
	Serve(context.Context) error
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

func (s *server) Serve(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	ctx, cxl := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	defer cxl()

	geg, gcx := errgroup.WithContext(ctx)
	reg, rcx := errgroup.WithContext(ctx)

	geg.Go(func() error {
		return s.grpc.Serve(ctx)
	})

	reg.Go(func() error {
		return s.rest.Serve(ctx)
	})

	geg.Go(func() error {
		<-gcx.Done()

		return nil
	})

	reg.Go(func() error {
		<-rcx.Done()

		return nil
	})

	ge := geg.Wait()
	re := reg.Wait()

	if ge != nil {
		ge = fmt.Errorf("server failure: %w", ge)

		s.logger.Error(ge, nil)
	}

	if re != nil {
		re = fmt.Errorf("server failure: %w", re)

		s.logger.Error(re, nil)
	}

	if ge != nil {
		return ge
	}

	return re
}
