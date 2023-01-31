package mderr

import (
	"google.golang.org/grpc/codes"
)

type mderr struct {
	message  string
	cause    error
	code     codes.Code
	metadata map[string]any
}

func New(msg string, c codes.Code, md map[string]any) error {
	return Wrap(msg, c, md, nil)
}

func Wrap(msg string, c codes.Code, md map[string]any, err error) error {
	return &mderr{
		message:  msg,
		cause:    err,
		code:     c,
		metadata: md,
	}
}

func (e *mderr) Error() string {
	return e.message
}

func (e *mderr) Code() codes.Code {
	return e.code
}

func (e *mderr) Metadata() map[string]any {
	return e.metadata
}

func (e *mderr) Unwrap() error {
	return e.cause
}

func (e *mderr) Root() error {
	c := e.Unwrap()

	if c == nil {
		return e
	}

	as, is := AsIs(c)
}

func As(err error) *mderr {
	mde, ok := err.(*mderr)

	if !ok {
		return nil
	}

	return mde
}

func Is(err error) bool {
	_, ok := err.(*mderr)

	return ok
}

func AsIs(err error) (*mderr, bool) {
	mde, ok := err.(*mderr)

	return mde, ok
}
