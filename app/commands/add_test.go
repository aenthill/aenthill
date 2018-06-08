package commands

import (
	"os"
	"testing"

	"github.com/aenthill/aenthill/app/context"

	"github.com/aenthill/log"
	"github.com/aenthill/manifest"
	"github.com/spf13/afero"
)

func TestAddCmd(t *testing.T) {
	t.Run("calling RunE without images as argument", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := &context.AppContext{ProjectDir: os.Getenv("HOST_PROJECT_DIR"), LogLevel: "DEBUG", EntryContext: &log.EntryContext{Source: "test"}}
		cmd := NewAddCmd(m, ctx)
		if err := cmd.RunE(nil, nil); err == nil {
			t.Error("RunE should have thrown an error as no image has been given")
		}
	})

	t.Run("calling RunE with a fake image as argument", func(t *testing.T) {
		image := "aent/foo"
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		if err := m.Flush(); err != nil {
			t.Errorf("an unexpected error occurred while trying to flush the given manifest: %s", err.Error())
		}
		ctx := &context.AppContext{ProjectDir: os.Getenv("HOST_PROJECT_DIR"), LogLevel: "DEBUG", EntryContext: &log.EntryContext{Source: "test"}}
		cmd := NewAddCmd(m, ctx)
		if err := cmd.RunE(nil, []string{image}); err == nil {
			t.Errorf("RunE should have thrown an error as the image %s is not valid", image)
		}
	})

	t.Run("calling RunE with a valid image as argument", func(t *testing.T) {
		image := "aenthill/cassandra"
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		if err := m.Flush(); err != nil {
			t.Errorf("an unexpected error occurred while trying to flush the given manifest: %s", err.Error())
		}
		ctx := &context.AppContext{ProjectDir: os.Getenv("HOST_PROJECT_DIR"), LogLevel: "DEBUG", EntryContext: &log.EntryContext{Source: "test"}}
		cmd := NewAddCmd(m, ctx)
		if err := cmd.RunE(nil, []string{image}); err != nil {
			t.Errorf("RunE should not have thrown an error as the image %s is valid", image)
		}
	})
}
