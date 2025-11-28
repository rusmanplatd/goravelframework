package schedule

import (
	"github.com/rusmanplatd/goravelframework/contracts/log"
	"github.com/rusmanplatd/goravelframework/support/color"
)

type Logger struct {
	log   log.Log
	debug bool
}

func NewLogger(log log.Log, debug bool) *Logger {
	return &Logger{
		debug: debug,
		log:   log,
	}
}

func (log *Logger) Info(msg string, keysAndValues ...any) {
	if !log.debug {
		return
	}
	color.Successf("%s %v\n", msg, keysAndValues)
}

func (log *Logger) Error(err error, msg string, keysAndValues ...any) {
	log.log.Error(msg, keysAndValues)
}
