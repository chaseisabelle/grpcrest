package client

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpcrest/gen/pbgen"
	"sync"
)

type Client interface {
	Open() error
	Close() error
	Reopen() error
	Service() (pbgen.UserServiceClient, error)
}

type client struct {
	config     Config
	mutex      sync.Mutex
	service    pbgen.UserServiceClient
	connection *grpc.ClientConn
}

func New(cfg Config) (Client, error) {
	return &client{
		config:     cfg,
		mutex:      sync.Mutex{},
		service:    nil,
		connection: nil,
	}, nil
}

func (c *client) Open() error {
	c.mutex.Lock()

	defer c.mutex.Unlock()

	if c.connection != nil {
		return fmt.Errorf("grpc client already connected")
	}

	con, err := grpc.Dial(c.config.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return fmt.Errorf("failed to connect grpc client: %w", err)
	}

	c.service = pbgen.NewUserServiceClient(con)
	c.connection = con

	return nil
}

func (c *client) Close() error {
	c.mutex.Lock()

	defer c.mutex.Unlock()

	if c.connection == nil {
		return nil
	}

	err := c.connection.Close()

	if err != nil {
		err = fmt.Errorf("failed to close grpc client connection: %w", err)
	}

	return err
}

func (c *client) Reopen() error {
	err := c.Close()

	if err != nil {
		return fmt.Errorf("failed to reopen grpc client: %w", err)
	}

	err = c.Open()

	if err != nil {
		err = fmt.Errorf("failed to reopen grpc client: %w", err)
	}

	return err
}

func (c *client) Service() (pbgen.UserServiceClient, error) {
	c.mutex.Lock()

	defer c.mutex.Unlock()

	if c.service == nil {
		return nil, fmt.Errorf("nil client service")
	}

	return c.service, nil
}
