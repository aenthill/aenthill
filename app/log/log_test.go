package log

import "testing"

func TestSetLevel(t *testing.T) {
	var level string

	// case 1: uses a wrong log level.
	level = "FOO"
	if err := SetLevel(level); err == nil {
		t.Errorf("Log level %s should have thrown an error.", level)
	}

	// case 2: uses a correct log level.
	level = "DEBUG"
	if err := SetLevel(level); err != nil {
		t.Errorf("Log level %s should not have thrown an error.", level)
	}
}

func TestWrongLogLeveLError(t *testing.T) {
	err := &wrongLogsLevelError{}
	if err.Error() != wrongLogsLevelErrorMessage {
		t.Errorf("Error returned a wrong message: got %s want %s.", err.Error(), wrongLogsLevelErrorMessage)
	}
}

func TestDumb(t *testing.T) {
	var (
		message = "foo"
		format  = "foo %s"
		args    = "bar"
		err     = &wrongLogsLevelError{}
	)

	// let's improve code coverage hehe!
	Debug(message)
	Debugf(format, args)
	Info(message)
	Infof(format, args)
	Warn(message)
	Warnf(format, args)
	Error(err)
}
