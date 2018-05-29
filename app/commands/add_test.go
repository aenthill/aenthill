package commands

import (
	"os"
	"testing"

	"github.com/aenthill/aenthill/app/context"

	"github.com/aenthill/manifest"
	"github.com/spf13/afero"
)

func TestNoImagesToAddError(t *testing.T) {
	err := &noImagesToAddError{}
	if err.Error() != noImagesToAddErrorMessage {
		t.Errorf("error returned a wrong message: got %s want %s", err.Error(), noImagesToAddErrorMessage)
	}
}

func TestAddCmd(t *testing.T) {
	t.Run("calling RunE without arguments", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := &context.AppContext{}
		cmd := NewAddCmd(m, ctx)
		if err := cmd.RunE(nil, nil); err == nil {
			t.Error("RunE should have thrown an error as there are no arguments")
		}
	})

	t.Run("calling RunE with a non-existing manifest file", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := &context.AppContext{}
		cmd := NewAddCmd(m, ctx)
		if err := cmd.RunE(nil, []string{"aenthill/cassandra"}); err == nil {
			t.Error("RunE should have thrown an error as the manifest file does not exist")
		}
	})

	t.Run("calling RunE with a broken manifest file", func(t *testing.T) {
		m := manifest.New("../../tests/aenthill-broken.json", afero.NewOsFs())
		ctx := &context.AppContext{}
		cmd := NewAddCmd(m, ctx)
		if err := cmd.RunE(nil, []string{"aenthill/cassandra"}); err == nil {
			t.Error("RunE should have thrown an error as the manifest file is broken")
		}
	})

	t.Run("calling RunE with a wrong application context", func(t *testing.T) {
		m := manifest.New("../../tests/aenthill.json", afero.NewOsFs())
		ctx := &context.AppContext{}
		cmd := NewAddCmd(m, ctx)
		if err := cmd.RunE(nil, []string{"aenthill/cassandra"}); err == nil {
			t.Error("RunE should have thrown an error as the application context is invalid")
		}
	})

	t.Run("calling RunE with a non-existing image as argument", func(t *testing.T) {
		image := "aenthill/cassandra"
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		if err := m.Flush(); err != nil {
			t.Errorf("an unexpected error occurred while trying to flush the given manifest: %s", err.Error())
		}
		ctx := &context.AppContext{LogLevel: "DEBUG", ProjectDir: os.Getenv("HOST_PROJECT_DIR")}
		cmd := NewAddCmd(m, ctx)
		if err := cmd.RunE(nil, []string{image}); err != nil {
			t.Errorf("RunE should not have thrown an error as the image %s should not exist in given manifest", image)
		}
	})

	t.Run("calling RunE with an existing image as argument", func(t *testing.T) {
		image := "aenthill/cassandra"
		m := manifest.New("../../tests/aenthill.json", afero.NewOsFs())
		ctx := &context.AppContext{LogLevel: "DEBUG", ProjectDir: os.Getenv("HOST_PROJECT_DIR")}
		cmd := NewAddCmd(m, ctx)
		if err := cmd.RunE(nil, []string{image}); err != nil {
			t.Errorf("RunE should not have thrown an error as the image %s should exist in given manifest", image)
		}
	})
}
