package log

import (
	"fmt"

	"github.com/kkrypt0nn/tangra"
)

var logger *tangra.Logger

func init() {
	logger = tangra.NewLogger()
}

func Debug(message string) {
	logger.Debug(message)
}

func Info(message string) {
	logger.Info(message)
}

func Success(message string) {
	logger.Println(fmt.Sprintf("${datetime} ${fg:green}${effect:bold}SUCCESS${reset}: %s", message))
}

func Warn(message string) {
	logger.Warn(message)
}

func Error(message string) {
	logger.Error(message)
}

func Fatal(message string) {
	logger.Fatal(message)
}

func Trace(message string) {
	logger.Trace(message)
}
