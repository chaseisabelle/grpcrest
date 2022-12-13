package zap

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"grpcrest/pkg/config"
	"os"
)

type Zap struct {
	logger *zap.Logger
}

func New(cfg config.Config) (*Zap, error) {
	ile := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.InfoLevel
	})

	ele := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.ErrorLevel || lvl == zapcore.FatalLevel
	})

	sos := zapcore.Lock(os.Stdout)
	ses := zapcore.Lock(os.Stderr)
	zoc := zap.NewProductionEncoderConfig()
	zec := zap.NewProductionEncoderConfig()

	if !cfg.Logger().Stack() {
		zoc.StacktraceKey = ""
		zec.StacktraceKey = ""
	}

	zoe := zapcore.NewJSONEncoder(zoc)
	zee := zapcore.NewJSONEncoder(zec)
	ozc := zapcore.NewCore(zoe, sos, ile)
	ezc := zapcore.NewCore(zee, ses, ele)
	tee := zapcore.NewTee(ozc, ezc)
	lgr := zap.New(tee)

	return &Zap{
		logger: lgr,
	}, nil
}

func (z *Zap) Fatal(err error, md map[string]any) {
	z.logger.Fatal(err.Error(), fields(md)...)
}

func (z *Zap) Error(err error, md map[string]any) {
	z.logger.Error(err.Error(), fields(md)...)
}

func (z *Zap) Info(msg string, md map[string]any) {
	z.logger.Info(msg, fields(md)...)
}

func (z *Zap) Debug(msg string, md map[string]any) {
	z.logger.Debug(msg, fields(md)...)
}

func fields(md map[string]any) []zapcore.Field {
	fds := make([]zapcore.Field, len(md))
	ind := 0

	for k, v := range md {
		fds[ind] = zap.Any(k, v)

		ind++
	}

	return fds
}
