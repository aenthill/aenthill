package jobs

import (
	"os"
	"testing"

	"github.com/aenthill/aenthill/context"
	"github.com/aenthill/aenthill/manifest"
	"github.com/spf13/afero"
)

func TestNewDispatchJob(t *testing.T) {
	if _, err := NewDispatchJob("FOO", "", nil, nil); err != nil {
		t.Errorf(`NewDispatchJob should not have thrown an error: got "%s"`, err.Error())
	}
}

func TestDispatchJobExecute(t *testing.T) {
	t.Run("calling Execute from dispatch job with a non-existing image in manifest", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		m.AddAent("aent/foo")
		ctx, err := context.New()
		if err != nil {
			t.Errorf(`An unexpected error occurred while creating the context: got "%s"`, err.Error())
		}

		ctx.HostProjectDir = os.Getenv("HOST_PROJECT_DIR")
		j, err := NewDispatchJob("FOO", "", ctx, m)
		if err != nil {
			t.Errorf(`An unexpected error occurred while creating a dispatch job: got "%s"`, err.Error())
		}
		if err := j.Execute(); err == nil {
			t.Error("Execute should have thrown an error as given image should not exist")
		}
	})
	t.Run("calling Execute from dispatch job with an existing image in manifest", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		m.AddAent("aenthill/cassandra")
		ctx, err := context.New()
		if err != nil {
			t.Errorf(`An unexpected error occurred while creating the context: got "%s"`, err.Error())
		}
		ctx.HostProjectDir = os.Getenv("HOST_PROJECT_DIR")
		j, err := NewDispatchJob("FOO", "", ctx, m)
		if err != nil {
			t.Errorf(`An unexpected error occurred while creating a dispatch job: got "%s"`, err.Error())
		}
		if err := j.Execute(); err != nil {
			t.Errorf(`Execute should not have thrown an error as given image should exist: got "%s"`, err.Error())
		}
	})
}
