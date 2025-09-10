package logger

import (
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	logrus.Logger
}

var instance *Logger
var once sync.Once

func GetLogger() *Logger {
	once.Do(func() {
		instance = &Logger{
			Logger: logrus.Logger{
				Out:       os.Stderr,
				Formatter: new(logrus.JSONFormatter),
				Hooks:     make(logrus.LevelHooks),
				Level:     logrus.InfoLevel,
				ExitFunc:  funcExit,
			},
		}
	})

	return instance
}

func funcExit(code int) {
	os.Exit(code)
}
