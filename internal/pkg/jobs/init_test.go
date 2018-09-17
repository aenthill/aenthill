package jobs

import (
	"testing"

	"github.com/aenthill/aenthill/internal/pkg/manifest"
	"github.com/aenthill/aenthill/test"

	"github.com/spf13/afero"
)

func TestNewInitJob(t *testing.T) {
	t.Run("calling NewInitJob", func(t *testing.T) {
		ctx := test.Context(t)
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		if _, err := NewInitJob(ctx, m); err != nil {
			t.Errorf(`NewInitJob should have not thrown an error: got "%s"`, err.Error())
		}
	})
}
