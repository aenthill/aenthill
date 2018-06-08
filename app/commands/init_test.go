package commands

import (
	"os"
	"testing"

	"github.com/aenthill/aenthill/app/context"

	"github.com/aenthill/log"
	"github.com/aenthill/manifest"
	"github.com/spf13/afero"
)

func TestInitCmd(t *testing.T) {
	t.Run("calling RunE with an existing manifest file", func(t *testing.T) {
		m := manifest.New("../../tests/aenthill.json", afero.NewOsFs())
		ctx := &context.AppContext{ProjectDir: os.Getenv("HOST_PROJECT_DIR"), LogLevel: "DEBUG", EntryContext: &log.EntryContext{Source: "test"}}
		cmd := NewInitCmd(m, ctx)
		if err := cmd.RunE(nil, nil); err == nil {
			t.Error("RunE should have thrown an error as the manifest file does exist")
		}
	})

	t.Run("calling RunE with a non-existing manifest file", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := &context.AppContext{ProjectDir: os.Getenv("HOST_PROJECT_DIR"), LogLevel: "DEBUG", EntryContext: &log.EntryContext{Source: "test"}}
		cmd := NewInitCmd(m, ctx)
		if err := cmd.RunE(nil, nil); err != nil {
			t.Error("RunE should not have thrown an error as the manifest file does not exist")
		}
	})
}
