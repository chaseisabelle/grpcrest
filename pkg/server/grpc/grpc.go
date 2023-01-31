package grpc

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"grpcrest/gen/pb"
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

func New(cfg config.Config, lgr logger.Logger, ser service.Service) (*GRPC, error) {
	uic := grpc.ChainUnaryInterceptor(func(ctx context.Context, req any, usi *grpc.UnaryServerInfo, han grpc.UnaryHandler) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)

		if !ok {
			md = metadata.New(map[string]string{})
		}

		res, err := han(ctx, req)

		s := res.(status.Status)

		md.Set("x-http-status-code", )

		grpc.SetHeader(ctx, md)
	})

	srv := grpc.NewServer(uic)

	pb.RegisterServiceServer(srv, ser)

	return &GRPC{
		server:  srv,
		config:  cfg.GRPC(),
		logger:  lgr,
		service: ser,
	}, nil
}

func (g *GRPC) Serve(ctx context.Context) error {
	adr := g.config.Address()
	lmd := map[string]any{
		"network": "tcp",
		"address": adr,
	}

	if ctx == nil {
		ctx = context.Background()
	}

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
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

		g.logger.Info("stopping grpc server", lmd)

		err := ctx.Err()

		if err != nil && err != context.Canceled {
			g.logger.Error(fmt.Errorf("grpc server context error: %w", err), lmd)
		}

		g.server.GracefulStop()
		g.logger.Info("stopped grpc server", lmd)

		return nil
	})

	err := eg.Wait()

	if err != nil {
		err = fmt.Errorf("grpc server failure: %w", err)
	}

	return err
}

