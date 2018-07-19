// Package log provides a simple logging solution.
package log

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

const (
	debugLevel = iota
	infoLevel
	warnLevel
	errorLevel
)

var levels = map[string]int{
	"DEBUG": debugLevel,
	"INFO":  infoLevel,
	"WARN":  warnLevel,
	"ERROR": errorLevel,
}

var currentLevel = errorLevel

var (
	debug = color.New(color.FgMagenta).FprintlnFunc()
	info  = color.New(color.FgCyan).FprintlnFunc()
	warn  = color.New(color.FgYellow).FprintlnFunc()
	error = color.New(color.FgRed).FprintlnFunc()
)

// SetLevel sets the log level.
func SetLevel(lvl string) {
	l, ok := levels[lvl]
	if ok {
		currentLevel = l
	}
}

// Debugf level formatted message.
func Debugf(format string, a ...interface{}) {
	Debug(fmt.Sprintf(format, a...))
}

// Debug level message.
func Debug(message string) {
	if currentLevel == debugLevel {
		debug(os.Stdout, fmt.Sprintf("%s: %s", "DEBUG", message))
	}
}

// Infof level formatted message.
func Infof(format string, a ...interface{}) {
	Info(fmt.Sprintf(format, a...))
}

// Info level message.
func Info(message string) {
	if currentLevel <= infoLevel {
		info(os.Stdout, fmt.Sprintf("%s: %s", "INFO", message))
	}
}

// Warnf level formatted message.
func Warnf(format string, a ...interface{}) {
	Warn(fmt.Sprintf(format, a...))
}

// Warn level message.
func Warn(message string) {
	if currentLevel <= warnLevel {
		warn(os.Stdout, fmt.Sprintf("%s: %s", "WARN", message))
	}
}

// Errorf level formatted message.
func Errorf(format string, a ...interface{}) {
	Error(fmt.Sprintf(format, a...))
}

// Error level message.
func Error(message string) {
	if currentLevel <= errorLevel {
		error(os.Stderr, fmt.Sprintf("%s: %s", "ERROR", message))
	}
}
