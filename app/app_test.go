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

	t.Run("calling entrypoint with a wrong log level", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		app := New("snapshot", m)
		app.ctx = &context.AppContext{LogLevel: "FOO"}
		if err := app.entrypoint(&cobra.Command{}, nil); err == nil {
			t.Errorf("entrypoint should have thrown an error as log level %s does not exist", app.ctx.LogLevel)
		}
	})

	t.Run("calling entrypoint with a valid log level", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		app := New("snapshot", m)
		app.ctx = &context.AppContext{LogLevel: "DEBUG"}
		if err := app.entrypoint(&cobra.Command{}, nil); err != nil {
			t.Errorf("entrypoint should not have thrown an error as log level %s does exist", app.ctx.LogLevel)
		}
	})
}

func TestInitialize(t *testing.T) {
	t.Run("calling initialize with a wrong log level", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		app := New("snapshot", m)
		app.ctx = &context.AppContext{LogLevel: "FOO"}
		if err := app.initialize(); err == nil {
			t.Errorf("initialize should have thrown an error as log level %s does not exist", app.ctx.LogLevel)
		}
	})

	t.Run("calling initialize with a valid log level", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		app := New("snapshot", m)
		app.ctx = &context.AppContext{LogLevel: "DEBUG"}
		if err := app.initialize(); err != nil {
			t.Errorf("initialize should not have thrown an error as log level %s does exist", app.ctx.LogLevel)
		}
	})

	t.Run("calling initialize with all parameters OK", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		app := New("snapshot", m)
		app.ctx = &context.AppContext{}
		if err := app.initialize(); err != nil {
			t.Error("initialize should not have thrown an error as all parameters are OK")
		}
	})
}
