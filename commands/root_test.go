package commands

import (
	"fmt"
	"os"
	"testing"

	"github.com/aenthill/manifest"
)

func copyManifest(fileName string) {
	oldPath := fmt.Sprintf("./../tests/%s", fileName)
	newPath := fmt.Sprintf("./%s", manifest.DefaultManifestFileName)
	os.Remove(newPath)
	os.Link(oldPath, newPath)
}

func setFlags() {
	logLevel = "DEBUG"
	projectDir = os.Getenv("HOST_PROJECT_DIR")
}

func unsetFlags() {
	logLevel = ""
	projectDir = ""
}

func TestManifestFileDoesNotExistError(t *testing.T) {
	err := &manifestFileDoesNotExistError{}
	expected := fmt.Sprintf(manifestFileDoesNotExistErrorMessage, manifest.DefaultManifestFileName, RootCmd.Use, InitCmd.Use)

	if err.Error() != expected {
		t.Errorf("error returned a wrong message: got %s want %s", err.Error(), expected)
	}
}

func TestWrongLogLevelError(t *testing.T) {
	err := &wrongLogLevelError{}
	if err.Error() != wrongLogLevelErrorMessage {
		t.Errorf("error returned a wrong message: got %s want %s", err.Error(), wrongLogLevelErrorMessage)
	}
}

func TestRootCmd(t *testing.T) {
	// case 1: uses a wrong log level.
	logLevel = "FOO"
	if err := RootCmd.PersistentPreRunE(nil, nil); err == nil {
		t.Errorf("root command should have thrown an error because the log level %s is not valid", logLevel)
	}

	// case 2: so far so good!
	logLevel = "DEBUG"
	if err := RootCmd.PersistentPreRunE(nil, nil); err != nil {
		t.Error("root command should have worked")
	}
}
