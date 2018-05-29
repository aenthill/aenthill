package commands

import (
	"fmt"
	"testing"

	"github.com/aenthill/manifest"
)

func TestManifestFileDoesNotExistError(t *testing.T) {
	err := &manifestFileDoesNotExistError{}
	expected := fmt.Sprintf(manifestFileDoesNotExistErrorMessage, manifest.DefaultManifestFileName)
	if err.Error() != expected {
		t.Errorf("error returned a wrong message: got %s want %s", err.Error(), expected)
	}
}
