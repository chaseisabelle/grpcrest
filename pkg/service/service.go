package service

import (
	"context"
	"grpcrest/gen/pb"
	"grpcrest/pkg/logger"
	"net/http"
)

type Service pb.ServiceServer

type service struct {
	logger logger.Logger
}

func New(lgr logger.Logger) (Service, error) {
	return &service{
		logger: lgr,
	}, nil
}

func (s *service) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	return &pb.GetResponse{
		Status: http.StatusOK,
		Name:   "foo",
	}, nil
}

func (s *service) Put(ctx context.Context, req *pb.PutRequest) (*pb.PutResponse, error) {
	return &pb.PutResponse{
		Status: http.StatusOK,
	}, nil
}

func (s *service) Post(ctx context.Context, req *pb.PostRequest) (*pb.PostResponse, error) {
	return &pb.PostResponse{
		Status: http.StatusCreated,
	}, nil
}

func (s *service) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	return &pb.DeleteResponse{
		Status: http.StatusOK,
	}, nil
}
