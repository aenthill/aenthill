package jobs

import (
	"testing"

	"github.com/aenthill/aenthill/internal/pkg/manifest"
	"github.com/aenthill/aenthill/test"

	"github.com/spf13/afero"
)

func TestNewRunJob(t *testing.T) {
	t.Run("calling NewRunJob with an invalid event", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		if _, err := NewRunJob("aent/foo", "%FOO%", "", nil, m); err == nil {
			t.Error("NewRunJob should have thrown an error as given event is not valid")
		}
	})
	t.Run("calling NewRunJob with a broken manifest", func(t *testing.T) {
		m := manifest.New(test.BrokenManifestAbsPath(t), afero.NewOsFs())
		if _, err := NewRunJob("aent/foo", "FOO", "", nil, m); err == nil {
			t.Error("NewRunJob should have thrown an error as given manifest is broken")
		}
	})
	t.Run("calling NewRunJob with a key", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := test.Context(t)
		ctx.ID = m.AddAent("aent/foo")
		if _, err := m.AddDependency(ctx.ID, "aent/foo", "FOO"); err != nil {
			t.Fatalf(`An unexpected error occurred while trying to add a dependency: got "%s"`, err.Error())
		}
		job, err := NewRunJob("FOO", "BAR", "", ctx, m)
		if err != nil {
			t.Errorf(`NewRunJob should not have thrown an error: got "%s"`, err.Error())
		}
		j, ok := job.(*runJob)
		if !ok {
			t.Fatal("Casting a Job to runJob should not have failed")
		}
		if j.image != "aent/foo" {
			t.Errorf(`NewRunJob has a wrong image: got "%s" want "%s"`, j.image, "aent/foo")
		}
	})
	t.Run("calling NewRunJob with an image", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		if _, err := NewRunJob("aent/foo", "FOO", "", test.Context(t), m); err != nil {
			t.Errorf(`NewRunJob should not have thrown an error: got "%s"`, err.Error())
		}
	})
}

func TestRunJobExecute(t *testing.T) {
	m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
	j, err := NewRunJob("aent/foo", "FOO", "", test.Context(t), m)
	if err != nil {
		t.Fatalf(`An unexpected error occurred while creating a reply job: got "%s"`, err.Error())
	}
	if err := j.Execute(); err == nil {
		t.Error("Execute should have thrown an error as given image should not exist")
	}
}
