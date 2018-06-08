package jobs

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/aenthill/aenthill/app/context"

	"github.com/aenthill/log"
	"github.com/aenthill/manifest"
	"github.com/spf13/afero"
)

func TestNoImageToRemoveError(t *testing.T) {
	err := &noImageToRemoveError{}
	if err.Error() != noImageToRemoveErrorMessage {
		t.Errorf("error returned a wrong message: got %s want %s", err.Error(), noImageToAddErrorMessage)
	}
}

func TestNewRemoveJob(t *testing.T) {
	t.Run("calling NewRemoveJob without images", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := &context.AppContext{ProjectDir: os.Getenv("HOST_PROJECT_DIR"), LogLevel: "DEBUG", EntryContext: &log.EntryContext{Source: "test"}}
		if _, err := NewRemoveJob(nil, m, ctx); err == nil {
			t.Error("NewRemoveJob should have thrown an error as there are no images in arguments")
		}
	})

	t.Run("calling NewRemoveJob with a non-existing manifest", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := &context.AppContext{ProjectDir: os.Getenv("HOST_PROJECT_DIR"), LogLevel: "DEBUG", EntryContext: &log.EntryContext{Source: "test"}}
		if _, err := NewRemoveJob([]string{"aent/foo"}, m, ctx); err == nil {
			t.Error("NewRemoveJob should have thrown an error as the given manifest should not exist")
		}
	})

	t.Run("calling NewRemoveJob with a broken manifest", func(t *testing.T) {
		m := manifest.New("../../tests/aenthill-broken.json", afero.NewOsFs())
		ctx := &context.AppContext{ProjectDir: os.Getenv("HOST_PROJECT_DIR"), LogLevel: "DEBUG", EntryContext: &log.EntryContext{Source: "test"}}
		if _, err := NewRemoveJob([]string{"aent/foo"}, m, ctx); err == nil {
			t.Error("NewRemoveJob should have thrown an error as the given manifest should be broken")
		}
	})

	t.Run("calling NewRemoveJob with a valid manifest", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := &context.AppContext{ProjectDir: os.Getenv("HOST_PROJECT_DIR"), LogLevel: "DEBUG", EntryContext: &log.EntryContext{Source: "test"}}
		if err := m.Flush(); err != nil {
			t.Errorf("an unexpected error occurred while flushing the given manifest: %s", err.Error())
		}
		if _, err := NewRemoveJob([]string{"aent/foo"}, m, ctx); err != nil {
			t.Error("NewRemoveJob should not have thrown an error as the given manifest should be valid")
		}
	})
}

func TestRemoveJobRun(t *testing.T) {
	t.Run("calling Run with a fake image", func(t *testing.T) {
		image := "aent/foo"
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := &context.AppContext{ProjectDir: os.Getenv("HOST_PROJECT_DIR"), LogLevel: "DEBUG", EntryContext: &log.EntryContext{Source: "test"}}
		if err := m.Flush(); err != nil {
			t.Errorf("an unexpected error occurred while flushing the given manifest: %s", err.Error())
		}
		job := &removeJob{[]string{image}, m, ctx}
		if err := job.Run(); err == nil {
			t.Errorf("Run should have thrown an error as the image %s is invalid", image)
		}
	})

	t.Run("calling Run with a valid image", func(t *testing.T) {
		image := "aenthill/cassandra"
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := &context.AppContext{ProjectDir: os.Getenv("HOST_PROJECT_DIR"), LogLevel: "DEBUG", EntryContext: &log.EntryContext{Source: "test"}}
		if err := m.Flush(); err != nil {
			t.Errorf("an unexpected error occurred while flushing the given manifest: %s", err.Error())
		}
		job := &removeJob{[]string{image}, m, ctx}
		if err := job.Run(); err != nil {
			t.Errorf("Run should not have thrown an error as the image %s is valid", image)
		}
	})
}

func TestRemoveJobHandle(t *testing.T) {
	t.Run("calling handle with a fake image", func(t *testing.T) {
		image := "aent/foo"
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := &context.AppContext{ProjectDir: os.Getenv("HOST_PROJECT_DIR"), LogLevel: "DEBUG", EntryContext: &log.EntryContext{Source: "test"}}
		if err := m.Flush(); err != nil {
			t.Errorf("an unexpected error occurred while flushing the given manifest: %s", err.Error())
		}
		job := &removeJob{[]string{image}, m, ctx}
		if err := job.Run(); err == nil {
			t.Errorf("handle should have thrown an error as the image %s is invalid", image)
		}
	})

	t.Run("calling handle with a fake image which does exist in the manifest", func(t *testing.T) {
		image := "aent/foo"
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := &context.AppContext{ProjectDir: os.Getenv("HOST_PROJECT_DIR"), LogLevel: "DEBUG", EntryContext: &log.EntryContext{Source: "test"}}
		if err := m.AddAent(image); err != nil {
			t.Errorf("an unexpected error occurred while adding aent %s in the given manifest: %s", image, err.Error())
		}
		if err := m.Flush(); err != nil {
			t.Errorf("an unexpected error occurred while flushing the given manifest: %s", err.Error())
		}
		job := &removeJob{[]string{image}, m, ctx}
		if err := job.Run(); err == nil {
			t.Errorf("handle should have thrown an error as the image %s is invalid", image)
		}
	})

	t.Run("calling handle with a valid image", func(t *testing.T) {
		image := "aenthill/cassandra"
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := &context.AppContext{ProjectDir: os.Getenv("HOST_PROJECT_DIR"), LogLevel: "DEBUG", EntryContext: &log.EntryContext{Source: "test"}}
		if err := m.Flush(); err != nil {
			t.Errorf("an unexpected error occurred while flushing the given manifest: %s", err.Error())
		}
		job := &removeJob{[]string{image}, m, ctx}
		if err := job.Run(); err != nil {
			t.Errorf("handle should not have thrown an error as the image %s is valid", image)
		}
	})
}

func TestRemoveJobRemoveAent(t *testing.T) {
	t.Run("calling removeAent with a non-existing image", func(t *testing.T) {
		image := "aent/foo"
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := &context.AppContext{ProjectDir: os.Getenv("HOST_PROJECT_DIR"), LogLevel: "DEBUG", EntryContext: &log.EntryContext{Source: "test"}}
		if err := m.Flush(); err != nil {
			t.Errorf("an unexpected error occurred while flushing the given manifest: %s", err.Error())
		}
		job := &removeJob{[]string{image}, m, ctx}
		if _, err := job.removeAent(image); err != nil {
			t.Error("removeAent should not have thrown an error")
		}
	})

	t.Run("calling removeAent with an existing image", func(t *testing.T) {
		image := "aent/foo"
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := &context.AppContext{ProjectDir: os.Getenv("HOST_PROJECT_DIR"), LogLevel: "DEBUG", EntryContext: &log.EntryContext{Source: "test"}}
		if err := m.AddAent(image); err != nil {
			t.Errorf("an unexpected error occurred while adding aent %s in the given manifest: %s", image, err.Error())
		}
		if err := m.Flush(); err != nil {
			t.Errorf("an unexpected error occurred while flushing the given manifest: %s", err.Error())
		}
		job := &removeJob{[]string{image}, m, ctx}
		if _, err := job.removeAent(image); err != nil {
			t.Error("removeAent should not have thrown an error")
		}
	})
}

func TestEventRemoveFailedError(t *testing.T) {
	image := "aent/foo"
	err := errors.New("")
	eventFailedErr := eventRemoveFailedError{image, err}
	expected := fmt.Sprintf(eventRemoveFailedErrorMessage, image, err)
	if eventFailedErr.Error() != expected {
		t.Errorf("error returned a wrong message: got %s want %s", eventFailedErr.Error(), expected)
	}
}

func TestRemoveJobSendEvent(t *testing.T) {
	t.Run("calling sendEvent with a fake image", func(t *testing.T) {
		image := "aent/foo"
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := &context.AppContext{ProjectDir: os.Getenv("HOST_PROJECT_DIR"), LogLevel: "DEBUG", EntryContext: &log.EntryContext{Source: "test"}}
		job := &removeJob{[]string{image}, m, ctx}
		if err := job.sendEvent(image); err == nil {
			t.Errorf("sendEvent should have thrown an error as the image %s is invalid", image)
		}
	})

	t.Run("calling sendEvent with a valid image", func(t *testing.T) {
		image := "aenthill/cassandra"
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := &context.AppContext{ProjectDir: os.Getenv("HOST_PROJECT_DIR"), LogLevel: "DEBUG", EntryContext: &log.EntryContext{Source: "test"}}
		job := &removeJob{[]string{image}, m, ctx}
		if err := job.sendEvent(image); err != nil {
			t.Errorf("sendEvent should not have thrown an error as the image %s is valid", image)
		}
	})
}

func TestRemoveJobReAddAent(t *testing.T) {
	t.Run("calling reAddAent with an existing image", func(t *testing.T) {
		image := "aent/foo"
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := &context.AppContext{ProjectDir: os.Getenv("HOST_PROJECT_DIR"), LogLevel: "DEBUG", EntryContext: &log.EntryContext{Source: "test"}}
		if err := m.AddAent(image); err != nil {
			t.Errorf("an unexpected error occurred while adding aent %s in the given manifest: %s", image, err.Error())
		}
		if err := m.Flush(); err != nil {
			t.Errorf("an unexpected error occurred while flushing the given manifest: %s", err.Error())
		}
		job := &removeJob{[]string{image}, m, ctx}
		if err := job.reAddAent(image); err == nil {
			t.Errorf("reAddAent should have thrown an error as the image %s should exist", image)
		}
	})

	t.Run("calling reAddAent with a non-existing image", func(t *testing.T) {
		image := "aent/foo"
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := &context.AppContext{ProjectDir: os.Getenv("HOST_PROJECT_DIR"), LogLevel: "DEBUG", EntryContext: &log.EntryContext{Source: "test"}}
		if err := m.Flush(); err != nil {
			t.Errorf("an unexpected error occurred while flushing the given manifest: %s", err.Error())
		}
		job := &removeJob{[]string{image}, m, ctx}
		if err := job.reAddAent(image); err != nil {
			t.Errorf("reAddAent should not have thrown an error as the image %s should not exist", image)
		}
	})
}
