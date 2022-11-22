package logger

import (
	"grpcrest/pkg/config"
	"grpcrest/pkg/logger/zerolog"
)

type Logger interface {
	Fatal(error, map[string]any)
	Error(error, map[string]any)
	Info(string, map[string]any)
	Debug(string, map[string]any)
}

func New(cfg config.Logger) (Logger, error) {
	return zerolog.New(cfg)
}
