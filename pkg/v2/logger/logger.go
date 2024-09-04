package logger

import "fmt"

type LoggerFactory interface {
	CreateLogger() Logger
}

type Logger interface {
	Log(string, ...any)
}

type loggerFactory struct{}

func NewLoggerFactory() LoggerFactory {
	return &loggerFactory{}
}

func (l *loggerFactory) CreateLogger() Logger {
	return &logger{}
}

type logger struct{}

func (l *logger) Log(message string, args ...any) {
	out := fmt.Sprintf(message, args...)
	fmt.Printf("%v\n", out)
}
