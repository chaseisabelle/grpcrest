package rest

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"grpcrest/gen/pb"
	"grpcrest/pkg/config"
	"grpcrest/pkg/logger"
	"grpcrest/pkg/service"
	"net"
	"net/http"
)

type REST struct {
	server  *http.Server
	config  config.Server
	logger  logger.Logger
	service service.Service
}

func New(cfg config.Config, lgr logger.Logger, ser service.Service) (*REST, error) {
	mux := http.NewServeMux()

	rmx := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(func(hdr string) (string, bool) {
		return hdr, true
	}), runtime.WithForwardResponseOption(func(ctx context.Context, res http.ResponseWriter, msg proto.Message) error {
		msg.ProtoReflect().Range(func(d protoreflect.FieldDescriptor, v protoreflect.Value) bool {
			if d.Name() != "status" {
				return true
			}

			i := v.Interface()
			s, ok := i.(int)

			if !ok {
				lgr.Error(fmt.Errorf("invalid status field value in grpc response"), map[string]any{
					"value": i,
				})
			} else {
				res.WriteHeader(s)
			}

			return false
		})
		return nil
	}))

	mux.Handle("/", rmx)

	err := pb.RegisterServiceHandlerFromEndpoint(context.Background(), rmx, cfg.GRPC().Address(), []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to register service handler: %w", err)
	}

	srv := &http.Server{
		Handler: mux,
	}

	return &REST{
		server:  srv,
		config:  cfg.REST(),
		logger:  lgr,
		service: ser,
	}, nil
}

func (r *REST) Serve(ctx context.Context) error {
	adr := r.config.Address()
	lmd := map[string]any{
		"network": "tcp",
		"address": adr,
	}

	if ctx == nil {
		ctx = context.Background()
	}

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		r.logger.Info("starting rest server", lmd)

		lis, err := net.Listen("tcp", adr)

		if err != nil {
			return fmt.Errorf("rest server listener failure: %w", err)
		}

		r.server.BaseContext = func(_ net.Listener) context.Context {
			return ctx
		}

		err = r.server.Serve(lis)

		if err != nil && err != http.ErrServerClosed {
			err = fmt.Errorf("failed to start rest server: %w", err)
		} else {
			err = nil
		}

		return err
	})

	eg.Go(func() error {
		<-ctx.Done()

		r.logger.Info("stopping rest server", lmd)

		err := ctx.Err()

		if err != nil && err != context.Canceled {
			r.logger.Error(fmt.Errorf("rest server context error: %w", err), lmd)
		}

		err = r.server.Shutdown(context.Background())

		if err != nil {
			err = fmt.Errorf("failed to gracefully shutdown rest server: %w", err)
		} else {
			r.logger.Info("stopped rest server", lmd)
		}

		return err
	})

	err := eg.Wait()

	if err != nil {
		err = fmt.Errorf("rest server failure: %w", err)
	}

	return err
}
