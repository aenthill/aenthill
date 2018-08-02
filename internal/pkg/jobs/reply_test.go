package jobs

import (
	"testing"

	"github.com/aenthill/aenthill/internal/pkg/manifest"
	"github.com/aenthill/aenthill/test"

	"github.com/spf13/afero"
)

func TestNewReplyJob(t *testing.T) {
	t.Run("calling NewReplyJob with an invalid event", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		if _, err := NewReplyJob("%FOO%", "", nil, m); err == nil {
			t.Error("NewReplyJob should have thrown an error as given event is not valid")
		}
	})
	t.Run("calling NewReplyJob with a valid event", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		if _, err := NewReplyJob("FOO", "", nil, m); err != nil {
			t.Errorf(`NewReplyJob should not have thrown an error as given event is valid: got "%s"`, err.Error())
		}
	})
}

func TestReplyJobExecute(t *testing.T) {
	m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
	j, err := NewReplyJob("FOO", "", test.Context(t), m)
	if err != nil {
		t.Fatalf(`An unexpected error occurred while creating a reply job: got "%s"`, err.Error())
	}
	if err := j.Execute(); err == nil {
		t.Error("Execute should have thrown an error as given container ID should not exist")
	}
}
