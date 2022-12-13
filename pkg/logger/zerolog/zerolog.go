package zerolog

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"grpcrest/pkg/config"
	"os"
	"strings"
)

type Zerolog struct {
	stdout zerolog.Logger
	stderr zerolog.Logger
}

func New(cfg config.Config) (*Zerolog, error) {
	if cfg.Logger().Stack() {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	sol := zerolog.New(os.Stdout).With().Timestamp().Logger()
	sel := zerolog.New(os.Stderr).With().Timestamp().Logger()

	if cfg.Logger().Stack() {
		sel = sel.With().Stack().Logger()
	}

	lvl := strings.ToLower(cfg.Logger().Level())

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

	sol = sol.Level(zll)
	sel = sel.Level(zll)

	return &Zerolog{
		stdout: sol,
		stderr: sel,
	}, nil
}

func (z *Zerolog) Fatal(err error, md map[string]any) {
	z.stderr.Fatal().Err(err).Fields(md).Msg("")
}

func (z *Zerolog) Error(err error, md map[string]any) {
	z.stderr.Error().Err(err).Fields(md).Msg("")
}

func (z *Zerolog) Info(msg string, md map[string]any) {
	z.stdout.Info().Fields(md).Msg(msg)
}

func (z *Zerolog) Debug(msg string, md map[string]any) {
	z.stdout.Debug().Fields(md).Msg(msg)
}
