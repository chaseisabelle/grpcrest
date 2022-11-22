package logger

import (
	"fmt"
	"grpcrest/pkg/config"
	"strings"
)

type Logger interface {
	Fatal(error, map[string]any)
	Error(error, map[string]any)
	Info(string, map[string]any)
	Debug(string, map[string]any)
}

type logger struct {
	level string
}

func New(cfg config.Config) (Logger, error) {
	lvl := strings.ToLower(cfg.Logger().Level())

	switch lvl {
	case "error":
		fallthrough
	case "debug":
	default:
		return nil, fmt.Errorf("invalid log level: %s", lvl)
	}

	return &logger{
		level: lvl,
	}, nil
}

func (l *logger) Fatal(err error, md map[string]any) {

}

func (l *logger) Error(err error, md map[string]any) {

}

func (l *logger) Info(msg string, md map[string]any) {

}

func (l *logger) Debug(msg string, md map[string]any) {

}
