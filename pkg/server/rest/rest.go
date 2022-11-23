package rest

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"grpcrest/gen/pbgen"
	"grpcrest/pkg/config"
	"grpcrest/pkg/logger"
	"grpcrest/pkg/service"
	"net"
	"net/http"
)

type REST struct {
	server  *http.ServeMux
	config  config.Server
	logger  logger.Logger
	service service.Service
}

func New(cfg config.Server, lgr logger.Logger, ser service.Service) (*REST, error) {
	mux := http.NewServeMux()
	rmx := runtime.NewServeMux()

	mux.Handle("/", rmx)

	err := pbgen.RegisterUserServiceHandlerServer(context.Background(), rmx, ser)

	http.Server{
		Handler: mux,
	}

	return &REST{
		server:  mux,
		config:  cfg,
		logger:  lgr,
		service: ser,
	}, err
}

func (r *REST) Serve() error {
	lis, err := net.Listen("tcp", r.config.Address())

	if err != nil {
		return fmt.Errorf("failed to listen on tcp address %s: %w", r.config.Address(), err)
	}

	err = http.Serve(lis, r.server)

	if err != nil {
		err = fmt.Errorf("failed to start rest server: %w", err)
	}

	return err
}

func (r *REST) Shutdown(ctx context.Context) error {
	return
}
