package service

import (
	"context"
	"grpcrest/gen/pbgen"
	"grpcrest/pkg/logger"
)

type Service pbgen.UserServiceServer

type service struct {
	logger logger.Logger
}

func New(lgr logger.Logger) (Service, error) {
	return &service{
		logger: lgr,
	}, nil
}

func (s *service) CreateUser(ctx context.Context, req *pbgen.CreateUserRequest) (*pbgen.CreateUserResponse, error) {
	return &pbgen.CreateUserResponse{
		Id: "todo",
	}, nil
}

func (s *service) GetUser(ctx context.Context, req *pbgen.GetUserRequest) (*pbgen.GetUserResponse, error) {
	return &pbgen.GetUserResponse{
		User: &pbgen.UserRead{
			Id:   req.Id,
			Name: "todo",
			Type: pbgen.UserType_USER_TYPE_ADMIN,
		},
	}, nil
}
