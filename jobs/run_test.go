package jobs

import (
	"testing"

	"github.com/aenthill/aenthill/manifest"
	"github.com/aenthill/aenthill/tests"

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
		m := manifest.New("../tests/aenthill-broken.json", afero.NewOsFs())
		if _, err := NewRunJob("aent/foo", "FOO", "", nil, m); err == nil {
			t.Error("NewRunJob should have thrown an error as given manifest is broken")
		}
	})
	t.Run("calling NewRunJob with an ID", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ID := m.AddAent("aent/foo")
		if _, err := NewRunJob(ID, "FOO", "", nil, m); err != nil {
			t.Errorf(`NewReplyJob should not have thrown an error: got "%s"`, err.Error())
		}
	})
	t.Run("calling NewRunJob with an image", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		if _, err := NewRunJob("aent/foo", "FOO", "", nil, m); err != nil {
			t.Errorf(`NewReplyJob should not have thrown an error: got "%s"`, err.Error())
		}
	})
}

func TestRunJobExecute(t *testing.T) {
	m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
	j, err := NewRunJob("aent/foo", "FOO", "", tests.MakeTestContext(t), m)
	if err != nil {
		t.Fatalf(`An unexpected error occurred while creating a reply job: got "%s"`, err.Error())
	}
	if err := j.Execute(); err == nil {
		t.Error("Execute should have thrown an error as given image should not exist")
	}
}
