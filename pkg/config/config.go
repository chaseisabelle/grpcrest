package config

import (
	"flag"
	"os"
)

type Config interface {
	GRPC() Server
	REST() Server
	Logger() Logger
}

type config struct {
	grpc   Server
	rest   Server
	logger Logger
}

func New() (Config, error) {
	fs := flag.NewFlagSet("", flag.ExitOnError)

	ga := fs.String("grpc-address", ":3333", "grpc server address")
	ra := fs.String("rest-address", ":8080", "rest server address")
	ll := fs.String("log-level", "debug", "log level")
	ls := fs.Bool("log-stack", false, "log stack trace with errors")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return nil, err
	}

	return &config{
		grpc: &server{
			address: *ga,
		},
		rest: &server{
			address: *ra,
		},
		logger: &logger{
			level: *ll,
			stack: *ls,
		},
	}, nil
}

func (c *config) GRPC() Server {
	return c.grpc
}

func (c *config) REST() Server {
	return c.rest
}

func (c *config) Logger() Logger {
	return c.logger
}
