package main

import (
	"fmt"
	"grpcrest/pkg/config"
	"grpcrest/pkg/logger"
	"grpcrest/pkg/server"
	"grpcrest/pkg/service"
)

func main() {
	cfg, err := config.New()

	if err != nil {
		panic(fmt.Errorf("failed to initialize configurations: %w", err))
	}

	lgr, err := logger.New(cfg)

	if err != nil {
		panic(fmt.Errorf("failed to initialize logger: %w", err))
	}

	ser, err := service.New(lgr)

	if err != nil {
		lgr.Fatal(fmt.Errorf("failed to initialize service: %w", err), nil)
	}

	srv, err := server.New(cfg, lgr, ser)

	if err != nil {
		lgr.Fatal(fmt.Errorf("failed to initialize server: %w", err), nil)
	}

	err = srv.Serve(nil)

	if err != nil {
		lgr.Fatal(fmt.Errorf("failed to start server: %w", err), nil)
	}
}
