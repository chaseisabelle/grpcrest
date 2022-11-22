package config

type Server interface {
	Address() string
}

type server struct {
	address string
}

func (s *server) Address() string {
	return s.address
}
