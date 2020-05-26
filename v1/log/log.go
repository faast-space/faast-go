package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()

	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.TraceLevel)
	logger.SetOutput(os.Stdout)
}

// access to underlying Logrus logger
func Log() logrus.FieldLogger {
	return logger
}

func WithField(key string, value interface{}) *logrus.Entry {
	return logger.WithField(key, value)
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return logger.WithFields(fields)
}
func WithError(err error) *logrus.Entry {
	return logger.WithError(err)
}

func Debug(msg string, args ...interface{}) {
	logger.Debugf(msg, args...)
}

func Info(msg string, args ...interface{}) {
	logger.Infof(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	logger.Warnf(msg, args)
}

func Error(msg string, args ...interface{}) {
	logger.Errorf(msg, args...)
}

func Fatal(msg string, args ...interface{}) {
	logger.Fatalf(msg, args...)
}
