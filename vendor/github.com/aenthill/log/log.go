// Package log is a simple wrapper around the https://github.com/apex/log library.
package log

import (
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	isatty "github.com/mattn/go-isatty"
)

func init() {
	if isatty.IsTerminal(os.Stdin.Fd()) {
		log.SetHandler(cli.Default)
	}

	log.SetLevel(log.InfoLevel)
}

const (
	// DebugLevel debug level.
	DebugLevel string = "DEBUG"
	// InfoLevel info level.
	InfoLevel string = "INFO"
	// WarnLevel warn level.
	WarnLevel string = "WARN"
	// ErrorLevel error level.
	ErrorLevel string = "ERROR"
)

// levels associates log levels as used within Aenthill ecosystem
// with its counterpart from the https://github.com/apex/log library.
var levels = map[string]log.Level{
	DebugLevel: log.DebugLevel,
	InfoLevel:  log.InfoLevel,
	WarnLevel:  log.WarnLevel,
	ErrorLevel: log.ErrorLevel,
}

type wrongLogLevelError struct {
	providedLogLevel string
}

const wrongLogLevelErrorMessage = "wrong log level provided (got %s); accepted values for log level: %s, %s, %s, %s"

func (e *wrongLogLevelError) Error() string {
	return fmt.Sprintf(wrongLogLevelErrorMessage, e.providedLogLevel, DebugLevel, InfoLevel, WarnLevel, ErrorLevel)
}

// SetLevel sets the log level.
func SetLevel(logLevel string) error {
	l, ok := levels[logLevel]
	if !ok {
		return &wrongLogLevelError{logLevel}
	}
	log.SetLevel(l)

	return nil
}

// EntryContext is an implementation of log.Fielder.
type EntryContext struct {
	Event         string
	Payload       string
	FromImageName string
	ToImageName   string
}

// Fields implements log.Fielder.
func (ctx *EntryContext) Fields() log.Fields {
	fields := make(log.Fields)
	fields["src"] = os.Args[0]

	if ctx.FromImageName != "" {
		fields["aent"] = ctx.FromImageName
	}

	if ctx.Event != "" {
		fields["event"] = ctx.Event

		if ctx.ToImageName != "" {
			fields["recipient_aent"] = ctx.ToImageName
		}

		if len(ctx.Payload) > 0 {
			fields["with_payload"] = "yes"
		} else {
			fields["with_payload"] = "no"
		}
	}

	return fields
}

// Debug level message.
func Debug(ctx *EntryContext, msg string) {
	log.WithFields(ctx).Debug(msg)
}

// Debugf level formatted message.
func Debugf(ctx *EntryContext, msg string, v ...interface{}) {
	log.WithFields(ctx).Debugf(msg, v...)
}

// Info level message.
func Info(ctx *EntryContext, msg string) {
	log.WithFields(ctx).Info(msg)
}

// Infof level formatted message.
func Infof(ctx *EntryContext, msg string, v ...interface{}) {
	log.WithFields(ctx).Infof(msg, v...)
}

// Warn level message.
func Warn(ctx *EntryContext, msg string) {
	log.WithFields(ctx).Warn(msg)
}

// Warnf level formatted message.
func Warnf(ctx *EntryContext, msg string, v ...interface{}) {
	log.WithFields(ctx).Warnf(msg, v...)
}

// Error level message.
func Error(ctx *EntryContext, err error, msg string) {
	log.WithFields(ctx).WithError(err).Error(msg)
}

// Errorf level formatted message.
func Errorf(ctx *EntryContext, err error, msg string, v ...interface{}) {
	log.WithFields(ctx).WithError(err).Errorf(msg, v...)
}
