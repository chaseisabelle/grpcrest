package main

import (
	"context"
	"fmt"
	"grpcrest/gen/pbgen"
	"grpcrest/pkg/config"
	"grpcrest/pkg/logger"
	"grpcrest/pkg/service"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.New()

	if err != nil {
		panic(fmt.Errorf("failed to initialize configurations: %w", err))
	}

	lgr, err := logger.New(cfg.Logger())

	if err != nil {
		panic(fmt.Errorf("failed to initialize logger: %w", err))
	}

	ser, err := service.New(lgr)

	if err != nil {
		lgr.Fatal(fmt.Errorf("failed to initialize service: %w", err), nil)
	}

	///////-----------

	// Setup gRPC servers.
	//
	baseGrpcServer := grpc.NewServer()
	userGrpcServer := NewUserGRPCServer()
	pbgen.RegisterUserServiceServer(baseGrpcServer, userGrpcServer)

	// Setup gRPC gateway.
	//
	ctx := context.Background()
	rmux := runtime.NewServeMux()
	mux := http.NewServeMux()
	mux.Handle("/", rmux)

	err := pbgen.RegisterUserServiceHandlerServer(ctx, rmux, userGrpcServer)
	if err != nil {
		log.Fatal(err)
	}

	// Serve.
	//

	grpcListener, err := net.Listen("tcp", *grpcAddr)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		err := baseGrpcServer.Serve(grpcListener)

		if err != grpc.ErrServerStopped {
			panic(err)
		}
	}()

	httpListener, err := net.Listen("tcp", *httpAddr)
	if err != nil {
		log.Fatal(err)
	}
	g.Add(func() error {
		log.Printf("Serving http address %s", *httpAddr)
		return http.Serve(httpListener, mux)
	}, func(err error) {
		httpListener.Close()
	})

	cancelInterrupt := make(chan struct{})
	g.Add(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-c:
			return fmt.Errorf("received signal %s", sig)
		case <-cancelInterrupt:
			return nil
		}
	}, func(error) {
		close(cancelInterrupt)
	})

	if err := g.Run(); err != nil {
		log.Fatal(err)
	}
}
