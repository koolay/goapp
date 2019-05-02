package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

// NewLogger new logger instance
func NewLogger(levelName string) *logrus.Logger {

	level, err := logrus.ParseLevel(levelName)
	if err != nil {
		level = logrus.ErrorLevel
	}
	logger := logrus.New()
	logger.SetLevel(level)
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})

	return logger
}

func NewUnitLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	return logger
}
