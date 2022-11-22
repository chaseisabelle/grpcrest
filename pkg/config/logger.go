package config

type Logger interface {
	Level() string
	Stack() bool
}

type logger struct {
	level string
	stack bool
}

func (l *logger) Level() string {
	return l.level
}

func (l *logger) Stack() bool {
	return l.stack
}
