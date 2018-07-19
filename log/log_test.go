package log

import "testing"

func TestSetLevel(t *testing.T) {
	SetLevel("INFO")
	if currentLevel != infoLevel {
		t.Error("SetLevel should have updated the log level correctly")
	}
	currentLevel = errorLevel
}

// dumb test.
func TestDebug(t *testing.T) {
	currentLevel = debugLevel
	Debugf("%s", "FOO")
	currentLevel = errorLevel
}

// dumb test.
func TestInfo(t *testing.T) {
	currentLevel = infoLevel
	Infof("%s", "FOO")
	currentLevel = errorLevel
}

// dumb test.
func TestWarn(t *testing.T) {
	currentLevel = warnLevel
	Warnf("%s", "FOO")
	currentLevel = errorLevel
}

// dumb test.
func TestError(t *testing.T) {
	currentLevel = errorLevel
	Errorf("%s", "FOO")
}
