package jobs

import (
	"testing"

	"github.com/aenthill/aenthill/manifest"
	"github.com/aenthill/aenthill/tests"

	"github.com/spf13/afero"
)

func TestNewUpdateJob(t *testing.T) {
	t.Run("calling NewUpdateJob with a non-existing manifest", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		if _, err := NewUpdateJob(nil, nil, nil, m); err == nil {
			t.Error("NewUpdateJob should have thrown an error as given manifest does not exist")
		}
	})
	t.Run("calling NewUpdateJob with empty key in context", func(t *testing.T) {
		m := manifest.New("../tests/aenthill.json", afero.NewOsFs())
		ctx := tests.MakeTestContext(t)
		if _, err := NewUpdateJob(nil, nil, ctx, m); err == nil {
			t.Error("NewUpdateJob should have thrown an error as context has no key")
		}
	})
	t.Run("calling NewUpdateJob with valid parameters", func(t *testing.T) {
		m := manifest.New("../tests/aenthill.json", afero.NewOsFs())
		ctx := tests.MakeTestContext(t)
		ctx.Key = "FOO"
		if _, err := NewUpdateJob(nil, nil, ctx, m); err != nil {
			t.Errorf(`NewUpdateJob should not have thrown an error: got "%s"`, err.Error())
		}
	})
}

// nolint: gocyclo
func TestUpdateJobExecute(t *testing.T) {
	t.Run("calling Execute from update job with an invalid metadata", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := tests.MakeTestContext(t)
		ctx.Key = m.AddAent("aent/foo")
		if err := m.Flush(); err != nil {
			t.Fatalf(`An unexpected error occurred while flushing manifest: got "%s"`, err.Error())
		}
		j, err := NewUpdateJob([]string{"FOO:bar"}, nil, ctx, m)
		if err != nil {
			t.Fatalf(`An unexpected error occurred while creating an update job: got "%s"`, err.Error())
		}
		if err := j.Execute(); err == nil {
			t.Error("Execute should have thrown an error as given metadata is not valid")
		}
	})
	t.Run("calling Execute from update job with an invalid event", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := tests.MakeTestContext(t)
		ctx.Key = m.AddAent("aent/foo")
		if err := m.Flush(); err != nil {
			t.Fatalf(`An unexpected error occurred while flushing manifest: got "%s"`, err.Error())
		}
		j, err := NewUpdateJob(nil, []string{"%FOO%"}, ctx, m)
		if err != nil {
			t.Fatalf(`An unexpected error occurred while creating an update job: got "%s"`, err.Error())
		}
		if err := j.Execute(); err == nil {
			t.Error("Execute should have thrown an error as given event is not valid")
		}
	})
	t.Run("calling Execute from update job with valid parameters", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := tests.MakeTestContext(t)
		ctx.Key = m.AddAent("aent/foo")
		if err := m.Flush(); err != nil {
			t.Fatalf(`An unexpected error occurred while flushing manifest: got "%s"`, err.Error())
		}
		j, err := NewUpdateJob([]string{"FOO=bar"}, []string{"FOO"}, ctx, m)
		if err != nil {
			t.Fatalf(`An unexpected error occurred while creating an update job: got "%s"`, err.Error())
		}
		if err := j.Execute(); err != nil {
			t.Errorf(`Execute should not have thrown an error: got "%s"`, err.Error())
		}
		j, err = NewUpdateJob(nil, []string{"FOO"}, ctx, m)
		if err != nil {
			t.Fatalf(`An unexpected error occurred while creating an update job: got "%s"`, err.Error())
		}
		if err := j.Execute(); err != nil {
			t.Errorf(`Execute should not have thrown an error: got "%s"`, err.Error())
		}
	})
}
