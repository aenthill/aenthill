package jobs

import (
	"fmt"
	"os"
	"testing"

	"github.com/aenthill/aenthill/app/context"
	"github.com/aenthill/log"
	"github.com/aenthill/manifest"
	"github.com/spf13/afero"
)

func TestManifestFileAlreadyExistingError(t *testing.T) {
	err := manifestFileAlreadyExistingError{manifest.DefaultManifestFileName}
	expected := fmt.Sprintf(manifestFileAlreadyExistingErrorMessage, manifest.DefaultManifestFileName)
	if err.Error() != expected {
		t.Errorf("error returned a wrong message: got %s want %s", err.Error(), expected)
	}
}

func TestNewInitJob(t *testing.T) {
	t.Run("calling NewInitJob with an existing manifest", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := &context.AppContext{ProjectDir: os.Getenv("HOST_PROJECT_DIR"), LogLevel: "DEBUG", EntryContext: &log.EntryContext{Source: "test"}}
		if err := m.Flush(); err != nil {
			t.Errorf("an unexpected error occurred while flushing the given manifest: %s", err.Error())
		}
		if _, err := NewInitJob(m, ctx); err == nil {
			t.Error("NewInitJob should have thrown an error as the given manifest should exist")
		}
	})

	t.Run("calling NewInitJob with a non-existing manifest", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := &context.AppContext{ProjectDir: os.Getenv("HOST_PROJECT_DIR"), LogLevel: "DEBUG", EntryContext: &log.EntryContext{Source: "test"}}
		if _, err := NewInitJob(m, ctx); err != nil {
			t.Error("NewInitJob should not have thrown an error as the given manifest should not exist")
		}
	})
}

func TestInitJobRun(t *testing.T) {
	t.Run("calling Run with a non-existing manifest", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := &context.AppContext{ProjectDir: os.Getenv("HOST_PROJECT_DIR"), LogLevel: "DEBUG", EntryContext: &log.EntryContext{Source: "test"}}
		job := &initJob{m, ctx}
		if err := job.Run(); err != nil {
			t.Error("Run should not have thrown an error as the given manifest should not exist")
		}
	})
}
