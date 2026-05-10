package log

import (
	"github.com/kkrypt0nn/tangra/v2"
)

var logger tangra.Logger

func init() {
	logger = tangra.NewLogger()
}

func Debug(message string) {
	logger.Debug(message)
}

func Debugf(format string, args ...any) {
	logger.Debugf(format, args...)
}

func Info(message string) {
	logger.Info(message)
}

func Infof(format string, args ...any) {
	logger.Infof(format, args...)
}

func Warn(message string) {
	logger.Warn(message)
}

func Warnf(format string, args ...any) {
	logger.Warnf(format, args...)
}

func Error(message string) {
	logger.Error(message)
}

func Errorf(format string, args ...any) {
	logger.Errorf(format, args...)
}

func Fatal(message string) {
	logger.Fatal(message)
}

func Fatalf(format string, args ...any) {
	logger.Fatalf(format, args...)
}

func Trace(message string) {
	logger.Trace(message)
}

func Tracef(format string, args ...any) {
	logger.Tracef(format, args...)
}
