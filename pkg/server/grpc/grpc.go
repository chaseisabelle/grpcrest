package grpc

import (
	"google.golang.org/grpc"
	"grpcrest/gen/pbgen"
	"grpcrest/pkg/config"
	"grpcrest/pkg/logger"
	"grpcrest/pkg/service"
)

type GRPC struct {
}

func New(cfg config.Server, lgr logger.Logger, ser service.Service) (*GRPC, error) {
	pbgen.RegisterUserServiceServer(grpc.NewServer(), ser)
}

func (g *GRPC) Serve() error {

}
