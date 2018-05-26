package commands

import (
	"fmt"
	"os"
	"testing"

	"github.com/aenthill/manifest"
)

func TestManifestFileAlreadyExistingError(t *testing.T) {
	err := &manifestFileAlreadyExistingError{}
	expected := fmt.Sprintf(manifestFileAlreadyExistingErrorMessage, manifest.DefaultManifestFileName)

	if err.Error() != expected {
		t.Errorf("error returned a wrong message: got %s want %s", err.Error(), expected)
	}
}

func TestInitCmd(t *testing.T) {
	// case 1: the manifest does already exist.
	copyManifest("aenthill.json")
	if err := InitCmd.RunE(nil, nil); err == nil {
		t.Error("init command should not have worked because there is no manifest file")
	}
	os.Remove(manifest.DefaultManifestFileName)

	// TODO find a way to test user inputs.
}
