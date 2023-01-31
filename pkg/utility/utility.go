package utility

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"net/http"
)

var grpcCodesToHTTPCodes = map[codes.Code]int{
	codes.OK:                 http.StatusOK,
	codes.Canceled:           http.StatusRequestTimeout,
	codes.Unknown:            http.StatusInternalServerError,
	codes.InvalidArgument:    http.StatusBadRequest,
	codes.DeadlineExceeded:   http.StatusGatewayTimeout,
	codes.NotFound:           http.StatusNotFound,
	codes.AlreadyExists:      http.StatusConflict,
	codes.PermissionDenied:   http.StatusForbidden,
	codes.Unauthenticated:    http.StatusUnauthorized,
	codes.ResourceExhausted:  http.StatusTooManyRequests,
	codes.FailedPrecondition: http.StatusPreconditionFailed,
	codes.Aborted:            http.StatusConflict,
	codes.OutOfRange:         http.StatusBadRequest,
	codes.Unimplemented:      http.StatusNotImplemented,
	codes.Internal:           http.StatusInternalServerError,
	codes.Unavailable:        http.StatusServiceUnavailable,
	codes.DataLoss:           http.StatusInternalServerError,
}

func GRPCCodeToHTTPCode(c codes.Code) (int, error) {
	var err error

	hsc, ok := grpcCodesToHTTPCodes[c]

	if !ok {
		err = fmt.Errorf("failed to map gRPC code %d to HTTP status code", c)
		hsc = http.StatusInternalServerError
	}

	return hsc, err
}
