package jobs

import (
	"testing"

	"github.com/aenthill/aenthill/manifest"
	"github.com/aenthill/aenthill/tests"

	"github.com/spf13/afero"
)

func TestNewAddJob(t *testing.T) {
	t.Run("calling NewAddJob with a broken manifest", func(t *testing.T) {
		m := manifest.New("../tests/aenthill-broken.json", afero.NewOsFs())
		if _, err := NewAddJob("aent/foo", nil, m); err == nil {
			t.Error("NewAddJob should have thrown an error as given manifest is broken")
		}
	})
	t.Run("calling NewAddJob with valid parameters", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		if _, err := NewAddJob("aent/foo", nil, m); err != nil {
			t.Errorf(`NewAddJob should not have thrown an error: got "%s"`, err.Error())
		}
	})
}

func TestAddJobExecute(t *testing.T) {
	t.Run("calling Execute from add job with a non-existing image", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		j, err := NewAddJob("aent/foo", tests.MakeTestContext(t), m)
		if err != nil {
			t.Fatalf(`An unexpected error occurred while creating an install job: got "%s"`, err.Error())
		}
		if err := j.Execute(); err != nil {
			t.Errorf(`Execute should not have thrown an error: got "%s"`, err.Error())
		}
	})
	t.Run("calling Execute from add job with an existing image", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		j, err := NewAddJob("aenthill/cassandra", tests.MakeTestContext(t), m)
		if err != nil {
			t.Fatalf(`An unexpected error occurred while creating an install job: got "%s"`, err.Error())
		}
		if err := j.Execute(); err != nil {
			t.Errorf(`Execute should not have thrown an error as image should exist: got "%s"`, err.Error())
		}
	})
}
