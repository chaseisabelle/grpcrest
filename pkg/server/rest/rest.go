package rest

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/sync/errgroup"
	"grpcrest/gen/pbgen"
	"grpcrest/pkg/config"
	"grpcrest/pkg/logger"
	"grpcrest/pkg/service"
	"net"
	"net/http"
)

type REST struct {
	server  *http.Server
	config  config.Server
	logger  logger.Logger
	service service.Service
}

func New(cfg config.Server, lgr logger.Logger, ser service.Service) (*REST, error) {
	mux := http.NewServeMux()
	rmx := runtime.NewServeMux()

	mux.Handle("/", rmx)

	err := pbgen.RegisterServiceHandlerServer(context.Background(), rmx, ser)

	if err != nil {
		return nil, fmt.Errorf("failed to register service handler: %w", err)
	}

	srv := &http.Server{
		Handler: mux,
	}

	return &REST{
		server:  srv,
		config:  cfg,
		logger:  lgr,
		service: ser,
	}, nil
}

func (r *REST) Serve(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		adr := r.config.Address()
		lmd := map[string]any{
			"network": "tcp",
			"address": adr,
		}

		r.logger.Info("starting rest server", lmd)

		lis, err := net.Listen("tcp", adr)

		if err != nil {
			return fmt.Errorf("rest server listener failure: %w", err)
		}

		r.server.BaseContext = func(_ net.Listener) context.Context {
			return ctx
		}

		err = r.server.Serve(lis)

		if err != nil && err != http.ErrServerClosed {
			err = fmt.Errorf("failed to start rest server: %w", err)
		} else {
			err = nil
		}

		return err
	})

	eg.Go(func() error {
		<-ctx.Done()

		r.logger.Info("stopping rest server", nil)

		err := r.server.Shutdown(context.Background())

		if err != nil {
			err = fmt.Errorf("failed to gracefully shutdown rest server: %w", err)
		} else {
			r.logger.Info("stopped rest server", nil)
		}

		return err
	})

	err := eg.Wait()

	if err != nil {
		err = fmt.Errorf("rest server failure: %w", err)
	}

	return err
}
