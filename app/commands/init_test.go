package commands

import (
	"fmt"
	"testing"

	"github.com/aenthill/aenthill/app/context"

	"github.com/aenthill/manifest"
	"github.com/spf13/afero"
)

func TestManifestFileAlreadyExistingError(t *testing.T) {
	err := &manifestFileAlreadyExistingError{}
	expected := fmt.Sprintf(manifestFileAlreadyExistingErrorMessage, manifest.DefaultManifestFileName)
	if err.Error() != expected {
		t.Errorf("error returned a wrong message: got %s want %s", err.Error(), expected)
	}
}

func TestInitCmd(t *testing.T) {
	t.Run("calling RunE with an existing manifest file", func(t *testing.T) {
		m := manifest.New("../../tests/aenthill.json", afero.NewOsFs())
		ctx := &context.AppContext{}
		cmd := NewInitCmd(m, ctx)
		if err := cmd.RunE(nil, nil); err == nil {
			t.Error("RunE should have thrown an error as the manifest file does exist")
		}
	})

	t.Run("calling RunE with a non-existing manifest file", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := &context.AppContext{}
		cmd := NewInitCmd(m, ctx)
		if err := cmd.RunE(nil, nil); err != nil {
			t.Error("RunE should not have thrown an error as the manifest file does not exist")
		}
	})
}
