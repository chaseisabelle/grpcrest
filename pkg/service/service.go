package service

import (
	"context"
	"grpcrest/gen/pbgen"
	"grpcrest/pkg/logger"
)

type Service pbgen.ServiceServer

type service struct {
	logger logger.Logger
}

func New(lgr logger.Logger) (Service, error) {
	return &service{
		logger: lgr,
	}, nil
}

func (s *service) Create(ctx context.Context, req *pbgen.CreateRequest) (*pbgen.CreateResponse, error) {
	return &pbgen.CreateResponse{
		Id: 1,
	}, nil
}

func (s *service) Read(ctx context.Context, req *pbgen.ReadRequest) (*pbgen.ReadResponse, error) {
	return &pbgen.ReadResponse{
		Model: &pbgen.Model{
			Id:   1,
			Name: "test",
		},
	}, nil
}
