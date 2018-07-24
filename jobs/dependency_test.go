package jobs

import (
	"testing"

	"github.com/aenthill/aenthill/manifest"
	"github.com/aenthill/aenthill/tests"

	"github.com/spf13/afero"
)

func TestNewDependencyJob(t *testing.T) {
	t.Run("calling NewDependencyJob with a non-existing manifest", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		if _, err := NewDependencyJob("", nil, m); err == nil {
			t.Error("NewDependencyJob should have thrown an error as given manifest does not exist")
		}
	})
	t.Run("calling NewDependencyJob with empty ID in context", func(t *testing.T) {
		m := manifest.New("../tests/aenthill.json", afero.NewOsFs())
		ctx := tests.MakeTestContext(t)
		if _, err := NewDependencyJob("", ctx, m); err == nil {
			t.Error("NewDependencyJob should have thrown an error as context has no ID")
		}
	})
	t.Run("calling NewDependencyJob with valid parameters", func(t *testing.T) {
		m := manifest.New("../tests/aenthill.json", afero.NewOsFs())
		ctx := tests.MakeTestContext(t)
		ctx.ID = "FOO"
		if _, err := NewDependencyJob("", ctx, m); err != nil {
			t.Errorf(`NewDependencyJob should not have thrown an error: got "%s"`, err.Error())
		}
	})
}

// nolint: gocyclo
func TestDependencyJobExecute(t *testing.T) {
	t.Run("calling Execute from dependency job with an empty key", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := tests.MakeTestContext(t)
		ctx.ID = m.AddAent("aent/foo")
		if err := m.Flush(); err != nil {
			t.Fatalf(`An unexpected error occurred while flushing manifest: got "%s"`, err.Error())
		}
		j, err := NewDependencyJob("", ctx, m)
		if err != nil {
			t.Fatalf(`An unexpected error occurred while creating a dependency job: got "%s"`, err.Error())
		}
		if err := j.Execute(); err != nil {
			t.Errorf(`Execute should not have thrown an error with an empty key: got "%s"`, err.Error())
		}
	})
	t.Run("calling Execute from dependency job with an invalid key", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := tests.MakeTestContext(t)
		ctx.ID = m.AddAent("aent/foo")
		if err := m.Flush(); err != nil {
			t.Fatalf(`An unexpected error occurred while flushing manifest: got "%s"`, err.Error())
		}
		j, err := NewDependencyJob("FOO", ctx, m)
		if err != nil {
			t.Fatalf(`An unexpected error occurred while creating a dependency job: got "%s"`, err.Error())
		}
		if err := j.Execute(); err == nil {
			t.Error("Execute should have thrown an error with an invalid key")
		}
	})
	t.Run("calling Execute from dependency job with a valid key", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := tests.MakeTestContext(t)
		ctx.ID = m.AddAent("aent/foo")
		if _, err := m.AddDependency(ctx.ID, "aent/bar", "BAR"); err != nil {
			t.Fatalf(`An unexpected error occurred while adding a dependency: got "%s"`, err.Error())
		}
		if err := m.Flush(); err != nil {
			t.Fatalf(`An unexpected error occurred while flushing manifest: got "%s"`, err.Error())
		}
		j, err := NewDependencyJob("BAR", ctx, m)
		if err != nil {
			t.Fatalf(`An unexpected error occurred while creating a dependency job: got "%s"`, err.Error())
		}
		if err := j.Execute(); err != nil {
			t.Errorf(`Execute should not have thrown an error with valid key: got "%s"`, err.Error())
		}
	})
}
