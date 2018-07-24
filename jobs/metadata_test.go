package jobs

import (
	"testing"

	"github.com/aenthill/aenthill/manifest"
	"github.com/aenthill/aenthill/tests"

	"github.com/spf13/afero"
)

func TestNewMetadataJob(t *testing.T) {
	t.Run("calling NewMetadataJob with a non-existing manifest", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		if _, err := NewMetadataJob("", nil, m); err == nil {
			t.Error("NewMetadataJob should have thrown an error as given manifest does not exist")
		}
	})
	t.Run("calling NewMetadataJob with empty ID in context", func(t *testing.T) {
		m := manifest.New("../tests/aenthill.json", afero.NewOsFs())
		ctx := tests.MakeTestContext(t)
		if _, err := NewMetadataJob("", ctx, m); err == nil {
			t.Error("NewMetadataJob should have thrown an error as context has no ID")
		}
	})
	t.Run("calling NewMetadataJob with valid parameters", func(t *testing.T) {
		m := manifest.New("../tests/aenthill.json", afero.NewOsFs())
		ctx := tests.MakeTestContext(t)
		ctx.ID = "FOO"
		if _, err := NewMetadataJob("", ctx, m); err != nil {
			t.Errorf(`NewMetadataJob should not have thrown an error: got "%s"`, err.Error())
		}
	})
}

// nolint: gocyclo
func TestMetadataJobExecute(t *testing.T) {
	t.Run("calling Execute from metadata job with an empty ID", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := tests.MakeTestContext(t)
		ctx.ID = m.AddAent("aent/foo")
		if err := m.Flush(); err != nil {
			t.Fatalf(`An unexpected error occurred while flushing manifest: got "%s"`, err.Error())
		}
		j, err := NewMetadataJob("", ctx, m)
		if err != nil {
			t.Fatalf(`An unexpected error occurred while creating a metadata job: got "%s"`, err.Error())
		}
		ctx.ID = ""
		if err := j.Execute(); err == nil {
			t.Error("Execute should have thrown an error with an empty ID")
		}
	})
	t.Run("calling Execute from metadata job with an invalid key", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := tests.MakeTestContext(t)
		ctx.ID = m.AddAent("aent/foo")
		if err := m.Flush(); err != nil {
			t.Fatalf(`An unexpected error occurred while flushing manifest: got "%s"`, err.Error())
		}
		j, err := NewMetadataJob("FOO", ctx, m)
		if err != nil {
			t.Fatalf(`An unexpected error occurred while creating a metadata job: got "%s"`, err.Error())
		}
		if err := j.Execute(); err == nil {
			t.Error("Execute should have thrown an error with an invalid key")
		}
	})
	t.Run("calling Execute from metadata job with an empty key", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := tests.MakeTestContext(t)
		ctx.ID = m.AddAent("aent/foo")
		if err := m.AddMetadata(ctx.ID, map[string]string{"BAR": "foo"}); err != nil {
			t.Fatalf(`An unexpected error occurred while adding a metadata: got "%s"`, err.Error())
		}
		if err := m.Flush(); err != nil {
			t.Fatalf(`An unexpected error occurred while flushing manifest: got "%s"`, err.Error())
		}
		j, err := NewMetadataJob("", ctx, m)
		if err != nil {
			t.Fatalf(`An unexpected error occurred while creating a demetadatapendency job: got "%s"`, err.Error())
		}
		if err := j.Execute(); err != nil {
			t.Errorf(`Execute should not have thrown an error with an empty key: got "%s"`, err.Error())
		}
	})
	t.Run("calling Execute from metadata job with a valid key", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := tests.MakeTestContext(t)
		ctx.ID = m.AddAent("aent/foo")
		if err := m.AddMetadata(ctx.ID, map[string]string{"BAR": "foo"}); err != nil {
			t.Fatalf(`An unexpected error occurred while adding a metadata: got "%s"`, err.Error())
		}
		if err := m.Flush(); err != nil {
			t.Fatalf(`An unexpected error occurred while flushing manifest: got "%s"`, err.Error())
		}
		j, err := NewMetadataJob("BAR", ctx, m)
		if err != nil {
			t.Fatalf(`An unexpected error occurred while creating a metadata job: got "%s"`, err.Error())
		}
		if err := j.Execute(); err != nil {
			t.Errorf(`Execute should not have thrown an error with valid key: got "%s"`, err.Error())
		}
	})
}
