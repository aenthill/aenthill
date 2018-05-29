package app

import (
	"testing"

	"github.com/aenthill/aenthill/app/context"

	"github.com/aenthill/manifest"
	"github.com/spf13/afero"
)

func TestNew(t *testing.T) {
	version := "snapshot"
	m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
	app := New(m, version)
	if app.version != version {
		t.Errorf("New should have instantiate an App with a correct version: got %s want %s", app.version, version)
	}
}

func TestExecute(t *testing.T) {
	m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
	app := New(m, "snapshot")
	if err := app.Execute(); err != nil {
		t.Error("Execute should not have thrown an error as all parameters are OK")
	}
}

func TestWrongLogLevelError(t *testing.T) {
	err := &wrongLogLevelError{}
	if err.Error() != wrongLogLevelErrorMessage {
		t.Errorf("error returned a wrong message: got %s want %s", err.Error(), wrongLogLevelErrorMessage)
	}
}

func TestInitialize(t *testing.T) {
	t.Run("calling initialize with a wrong log level", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		app := New(m, "snapshot")
		app.ctx = &context.AppContext{LogLevel: "FOO"}
		if err := app.initialize(); err == nil {
			t.Errorf("initialize should have thrown an error as log level %s does not exist", app.ctx.LogLevel)
		}
	})

	t.Run("calling initialize with a valid log level", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		app := New(m, "snapshot")
		app.ctx = &context.AppContext{LogLevel: "DEBUG"}
		if err := app.initialize(); err != nil {
			t.Errorf("initialize should not have thrown an error as log level %s does exist", app.ctx.LogLevel)
		}
	})

	t.Run("calling initialize with all parameters OK", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		app := New(m, "snapshot")
		app.ctx = &context.AppContext{}
		if err := app.initialize(); err != nil {
			t.Error("initialize should not have thrown an error as all parameters are OK")
		}
	})
}
