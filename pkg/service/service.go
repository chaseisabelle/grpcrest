package service

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpcrest/gen/pb"
	"grpcrest/pkg/logger"
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
		Name: "foo",
	}, status.Errorf(codes.NotFound, "not found error weee")
}

func (s *service) Put(ctx context.Context, req *pb.PutRequest) (*pb.PutResponse, error) {
	return &pb.PutResponse{}, nil
}

func (s *service) Post(ctx context.Context, req *pb.PostRequest) (*pb.PostResponse, error) {
	return &pb.PostResponse{}, nil
}

func (s *service) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	return &pb.DeleteResponse{}, nil
}
