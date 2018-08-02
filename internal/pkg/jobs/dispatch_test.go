package jobs

import (
	"testing"

	"github.com/aenthill/aenthill/internal/pkg/manifest"
	"github.com/aenthill/aenthill/test"

	"github.com/spf13/afero"
)

func TestNewDispatchJob(t *testing.T) {
	t.Run("calling NewDispatchJob with an invalid event", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		if _, err := NewDispatchJob("%FOO%", "", nil, m); err == nil {
			t.Error("NewDispatchJob should have thrown an error as given event is not valid")
		}
	})
	t.Run("calling NewDispatchJob with a non-existing manifest", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		if _, err := NewDispatchJob("FOO", "", nil, m); err == nil {
			t.Error("NewDispatchJob should have thrown an error as given manifest does not exist")
		}
	})
	t.Run("calling NewDispatchJob with an existing manifest", func(t *testing.T) {
		m := manifest.New(test.ValidManifestAbsPath(t), afero.NewOsFs())
		if _, err := NewDispatchJob("FOO", "", nil, m); err != nil {
			t.Errorf(`NewDispatchJob should not have thrown an error: got "%s"`, err.Error())
		}
	})
}

func TestDispatchJobExecute(t *testing.T) {
	t.Run("calling Execute from dispatch job with a non-existing image in manifest", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		m.AddAent("aent/foo")
		if err := m.Flush(); err != nil {
			t.Fatalf(`An unexpected error occurred while flushing the manifest: got "%s"`, err.Error())
		}
		j, err := NewDispatchJob("FOO", "", test.Context(t), m)
		if err != nil {
			t.Fatalf(`An unexpected error occurred while creating a dispatch job: got "%s"`, err.Error())
		}
		if err := j.Execute(); err == nil {
			t.Error("Execute should have thrown an error as given image should not exist")
		}
	})
	t.Run("calling Execute from dispatch job with an existing image in manifest", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := test.Context(t)
		ctx.ID = m.AddAent("aenthill/cassandra")
		m.AddAent("aenthill/cassandra")
		if err := m.Flush(); err != nil {
			t.Fatalf(`An unexpected error occurred while flushing the manifest: got "%s"`, err.Error())
		}
		j, err := NewDispatchJob("FOO", "", ctx, m)
		if err != nil {
			t.Fatalf(`An unexpected error occurred while creating a dispatch job: got "%s"`, err.Error())
		}
		if err := j.Execute(); err != nil {
			t.Errorf(`Execute should not have thrown an error as given image should exist: got "%s"`, err.Error())
		}
	})
}
