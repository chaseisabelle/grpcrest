package grpc

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
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

func (g *GRPC) Serve(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		adr := g.config.Address()
		lmd := map[string]any{
			"network": "tcp",
			"address": adr,
		}

		g.logger.Info("starting grpc server", lmd)

		lis, err := net.Listen("tcp", adr)

		if err != nil {
			return fmt.Errorf("grpc server listener failure: %w", err)
		}

		err = g.server.Serve(lis)

		if err != nil && err != grpc.ErrServerStopped {
			err = fmt.Errorf("failed to start grpc server: %w", err)
		} else {
			err = nil
		}

		return err
	})

	eg.Go(func() error {
		<-ctx.Done()

		g.server.GracefulStop()

		return nil
	})

	err := eg.Wait()

	if err != nil {
		err = fmt.Errorf("grpc server failure: %w", err)
	}

	return err
}
