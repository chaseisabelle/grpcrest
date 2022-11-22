package config

type Logger interface {
	Level() string
}

type logger struct {
	level string
}

func (l *logger) Level() string {
	return l.level
}
