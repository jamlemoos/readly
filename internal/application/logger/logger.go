package logger

import "log"

type Logger interface {
	Warn(msg string, err error)
}

type stdLogger struct{}

func NewStdLogger() Logger { return &stdLogger{} }

func (l *stdLogger) Warn(msg string, err error) {
	log.Printf("[WARN] %s: %v", msg, err)
}
