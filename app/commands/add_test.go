package commands

import (
	"testing"

	"github.com/aenthill/aenthill/tests"

	"github.com/aenthill/manifest"
	"github.com/spf13/afero"
)

func TestAddCmd(t *testing.T) {
	t.Run("calling RunE without images as argument", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := tests.NewAppContext()
		cmd := NewAddCmd(m, ctx)
		if err := cmd.RunE(nil, nil); err == nil {
			t.Error("RunE should have thrown an error as no image has been given")
		}
	})

	t.Run("calling RunE with a fake image as argument", func(t *testing.T) {
		image := "aent/foo"
		m, err := tests.NewEmptyInMemoryManifest()
		if err != nil {
			t.Error(err)
		}
		ctx := tests.NewAppContext()
		cmd := NewAddCmd(m, ctx)
		if err := cmd.RunE(nil, []string{image}); err == nil {
			t.Errorf("RunE should have thrown an error as the image %s is not valid", image)
		}
	})

	t.Run("calling RunE with a valid image as argument", func(t *testing.T) {
		image := "aenthill/cassandra"
		m, err := tests.NewEmptyInMemoryManifest()
		if err != nil {
			t.Error(err)
		}
		ctx := tests.NewAppContext()
		cmd := NewAddCmd(m, ctx)
		if err := cmd.RunE(nil, []string{image}); err != nil {
			t.Errorf("RunE should not have thrown an error as the image %s is valid", image)
		}
	})
}
