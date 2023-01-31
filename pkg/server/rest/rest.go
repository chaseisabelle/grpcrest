package rest

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"grpcrest/gen/pb"
	"grpcrest/pkg/config"
	"grpcrest/pkg/logger"
	"grpcrest/pkg/service"
	"grpcrest/pkg/utility"
	"net"
	"net/http"
	"strconv"
)

type REST struct {
	server  *http.Server
	config  config.Server
	logger  logger.Logger
	service service.Service
}

func New(cfg config.Config, lgr logger.Logger, ser service.Service) (*REST, error) {
	mux := http.NewServeMux()

	rmx := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(headerMatcher), runtime.WithForwardResponseOption(responseForwarder), runtime.WithErrorHandler(func(ctx context.Context, mux *runtime.ServeMux, mar runtime.Marshaler, res http.ResponseWriter, req *http.Request, err error) {
		if err == nil {
			err = status.New(codes.Unknown, "unknown error").Err()
		}

		s, ok := status.FromError(err)

		if !ok {
			s = status.New(codes.Unknown, err.Error())
		}

		c := s.Code()
		hsc, err := utility.GRPCCodeToHTTPCode(c)
		
		if err != nil {
			lgr.Error(fmt.Errorf("failed to get HTTP status code: %w", err), map[string]any{
				"code": c,
			})
		}

		runtime.DefaultHTTPErrorHandler(ctx, mux, mar, res, req, &runtime.HTTPStatusError{
			HTTPStatus: hsc,
			Err:        err,
		})
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

func headerMatcher(hdr string) (string, bool) {
	return hdr, true
}

func responseForwarder(ctx context.Context, res http.ResponseWriter, msg proto.Message) error {
	md, ok := runtime.ServerMetadataFromContext(ctx)

	if !ok {
		return fmt.Errorf("failed to get metadata from context")
	}

	hdr := md.HeaderMD.Get("x-http-status-code")

	if len(hdr) != 1 {
		return fmt.Errorf("failed to get x-http-status-code response header")
	}

	hsc, err := strconv.Atoi(hdr[0])

	if err != nil {
		return fmt.Errorf("invalid x-http-status-code response header: %s", hdr[0])
	}

	res.WriteHeader(hsc)

	return nil
}
