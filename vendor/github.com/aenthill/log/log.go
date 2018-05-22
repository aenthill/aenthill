// Package log implements a simple wrapper of the logrus library.
package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

// logger is our logger instance used across the application.
var logger = newLogger()

// newLogger instantiates a logger instance with default values.
func newLogger() *logrus.Logger {
	l := logrus.New()
	l.Out = os.Stdout
	l.Level = logrus.InfoLevel

	return l
}

// levels associates log levels as used with the --logLevel -l flag
// with its counterpart from the logrus library.
var levels = map[string]logrus.Level{
	"DEBUG": logrus.DebugLevel,
	"INFO":  logrus.InfoLevel,
	"WARN":  logrus.WarnLevel,
	"ERROR": logrus.ErrorLevel,
	"FATAL": logrus.FatalLevel,
	"PANIC": logrus.PanicLevel,
}

type wrongLogsLevelError struct{}

const wrongLogsLevelErrorMessage = "accepted values for logs level: DEBUG, INFO, WARN, ERROR, FATAL, PANIC"

func (e *wrongLogsLevelError) Error() string {
	return wrongLogsLevelErrorMessage
}

// SetLevel updates the level of messages which will be logged.
func SetLevel(level string) error {
	l, ok := levels[level]
	if !ok {
		return &wrongLogsLevelError{}
	}

	logger.Level = l
	return nil
}

// Debug is a wrapper of the logrus Debug function.
func Debug(message string) {
	logger.Debug(message)
}

// Debugf is a wrapper of the logrus Debugf function.
func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

// Info is a wrapper of the logrus Info function.
func Info(message string) {
	logger.Info(message)
}

// Infof is a wrapper of the logrus Infof function.
func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

// Warn is a wrapper of the logrus Warn function.
func Warn(message string) {
	logger.Warn(message)
}

// Warnf is wrapper of the logrus Warnf function.
func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

// Error is a wrapper of the logrus Error function.
func Error(err error) {
	logger.Error(err.Error())
}
