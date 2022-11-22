package main

import (
	"context"
	"flag"
	"fmt"
	"grpcrest/gen/pbgen"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {

	// Flags.
	//
	fs := flag.NewFlagSet("", flag.ExitOnError)
	grpcAddr := fs.String("grpc-addr", ":6565", "grpc address")
	httpAddr := fs.String("http-addr", ":8080", "http address")
	if err := fs.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}

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

type userServer struct {
	m     map[string]*pbgen.UserWrite
	mutex *sync.RWMutex
}

func NewUserGRPCServer() pbgen.UserServiceServer {
	return &userServer{
		m:     map[string]*pbgen.UserWrite{},
		mutex: &sync.RWMutex{},
	}
}

func (s *userServer) CreateUser(ctx context.Context, req *pbgen.CreateUserRequest) (*pbgen.CreateUserResponse, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil,
			status.Error(codes.Internal, err.Error())
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.m[id.String()] = req.User
	return &pbgen.CreateUserResponse{
		Id: id.String(),
	}, nil
}

func (s *userServer) GetUser(ctx context.Context, req *pbgen.GetUserRequest) (*pbgen.GetUserResponse, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	foundUser, ok := s.m[req.Id]
	if !ok {
		return nil,
			status.Error(codes.NotFound, fmt.Errorf("User not found by id %v", req.Id).Error())
	}
	return &pbgen.GetUserResponse{
		User: &pbgen.UserRead{
			Id:   req.Id,
			Name: foundUser.Name,
			Type: foundUser.Type,
		},
	}, nil
}
