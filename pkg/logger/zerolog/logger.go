package zerolog

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"grpcrest/pkg/config"
	"os"
	"strings"
)

type Logger struct {
	stdout zerolog.Logger
	stderr zerolog.Logger
}

func New(cfg config.Logger) (*Logger, error) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if cfg.Stack() {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	}

	sol := zerolog.New(os.Stdout)
	sel := zerolog.New(os.Stderr)
	lvl := strings.ToLower(cfg.Level())

	var zll zerolog.Level

	switch lvl {
	case "fatal":
		zll = zerolog.FatalLevel
	case "error":
		zll = zerolog.ErrorLevel
	case "info":
		zll = zerolog.InfoLevel
	case "debug":
		zll = zerolog.DebugLevel
	default:
		return nil, fmt.Errorf("invalid log level: %s", lvl)
	}

	return &Logger{
		stdout: sol.Level(zll),
		stderr: sel.Level(zll),
	}, nil
}

func (z *Logger) Fatal(err error, md map[string]any) {
	z.stderr.Fatal().Stack().Err(err).Fields(md).Msg("")
}

func (z *Logger) Error(err error, md map[string]any) {
	z.stderr.Error().Stack().Err(err).Fields(md).Msg("")
}

func (z *Logger) Info(msg string, md map[string]any) {
	z.stdout.Info().Fields(md).Msg(msg)
}

func (z *Logger) Debug(msg string, md map[string]any) {
	z.stdout.Debug().Fields(md).Msg(msg)
}
