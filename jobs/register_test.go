package jobs

import (
	"testing"

	"github.com/aenthill/aenthill/manifest"
	"github.com/aenthill/aenthill/tests"

	"github.com/spf13/afero"
)

func TestNewRegisterJob(t *testing.T) {
	t.Run("calling NewRegisterJob with a non-existing manifest", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		if _, err := NewRegisterJob("aent/foo", "FOO", nil, nil, m); err == nil {
			t.Error("NewRegisterJob should have thrown an error as given manifest does not exist")
		}
	})
	t.Run("calling NewRegisterJob with an existing maxnifest", func(t *testing.T) {
		m := manifest.New("../tests/aenthill.json", afero.NewOsFs())
		if _, err := NewRegisterJob("aent/foo", "FOO", nil, nil, m); err != nil {
			t.Errorf(`NewDispatchJob should not have thrown an error as given manifest does exist: got "%s"`, err.Error())
		}
	})
}

func TestRegisterJobExecute(t *testing.T) {
	t.Run("calling Execute from register job with no ID in context", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		if err := m.Flush(); err != nil {
			t.Fatalf(`An unexpected error occurred while flushing the manifest: got "%s"`, err.Error())
		}
		m.AddAent("aent/bar")
		j, err := NewRegisterJob("aent/foo", "BAR", nil, tests.MakeTestContext(t), m)
		if err != nil {
			t.Fatalf(`An unexpected error occurred while creating a register job: got "%s"`, err.Error())
		}
		if err := j.Execute(); err == nil {
			t.Error("Execute should have thrown an error as there is no ID in context")
		}
	})
	t.Run("calling Execute from register job with an invalid metadata", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		if err := m.Flush(); err != nil {
			t.Fatalf(`An unexpected error occurred while flushing the manifest: got "%s"`, err.Error())
		}
		ctx := tests.MakeTestContext(t)
		ctx.ID = m.AddAent("aent/foo")
		j, err := NewRegisterJob("aent/bar", "BAR", []string{"FOO:bar"}, ctx, m)
		if err != nil {
			t.Fatalf(`An unexpected error occurred while creating a register job: got "%s"`, err.Error())
		}
		if err := j.Execute(); err == nil {
			t.Error("Execute should have thrown an error as given metadata is not valid")
		}
	})
	t.Run("calling Execute from register job with valid parameters", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		if err := m.Flush(); err != nil {
			t.Fatalf(`An unexpected error occurred while flushing the manifest: got "%s"`, err.Error())
		}
		ctx := tests.MakeTestContext(t)
		ctx.ID = m.AddAent("aent/foo")
		j, err := NewRegisterJob("aent/bar", "BAR", []string{"FOO=bar"}, ctx, m)
		if err != nil {
			t.Fatalf(`An unexpected error occurred while creating a register job: got "%s"`, err.Error())
		}
		if err := j.Execute(); err != nil {
			t.Errorf(`Execute should not have thrown an error: got "%s"`, err.Error())
		}
	})
}
