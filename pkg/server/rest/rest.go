package rest

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"grpcrest/gen/pbgen"
	"grpcrest/pkg/config"
	"grpcrest/pkg/logger"
	"grpcrest/pkg/service"
	"net/http"
)

type REST struct {
}

func New(cfg config.Server, lgr logger.Logger, ser service.Service) (*REST, error) {
	rmx := runtime.NewServeMux()

	http.NewServeMux().Handle("/", rmx)

	err := pbgen.RegisterUserServiceHandlerServer(context.Background(), rmx, ser)
}

func (r *REST) Serve() error {

}
