package app

import (
	"testing"

	"github.com/aenthill/aenthill/app/context"

	"github.com/aenthill/manifest"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func TestNew(t *testing.T) {
	version := "snapshot"
	m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
	app := New(version, m)
	if app.version != version {
		t.Errorf("New should have instantiate an App with a correct version: got %s want %s", app.version, version)
	}
}

func TestExecute(t *testing.T) {
	m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
	app := New("snapshot", m)
	if err := app.Execute(); err != nil {
		t.Error("Execute should not have thrown an error as all parameters are OK")
	}
}

func TestEntrypoint(t *testing.T) {
	t.Run("calling entrypoint without cobra.Command", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		app := New("snapshot", m)
		app.ctx = &context.AppContext{}
		if err := app.entrypoint(nil, nil); err != nil {
			t.Error("entrypoint should not have thrown an error as no cobra.Command has been provided")
		}
	})

	t.Run("calling entrypoint with no verbosity", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		app := New("snapshot", m)
		app.ctx = &context.AppContext{}
		if err := app.entrypoint(&cobra.Command{}, nil); err != nil {
			t.Error("entrypoint should not have thrown an error with default verbosity")
		}
	})

	t.Run("calling entrypoint with info verbosity", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		app := New("snapshot", m)
		app.ctx = &context.AppContext{IsVerbose: true}
		if err := app.entrypoint(&cobra.Command{}, nil); err != nil {
			t.Error("entrypoint should not have thrown an error with info verbosity")
		}
	})

	t.Run("calling entrypoint with debug verbosity", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		app := New("snapshot", m)
		app.ctx = &context.AppContext{IsVeryVerbose: true}
		if err := app.entrypoint(&cobra.Command{}, nil); err != nil {
			t.Error("entrypoint should not have thrown an error with debug verbosity")
		}
	})
}
