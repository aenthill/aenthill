package jobs

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/aenthill/aenthill/app/context"

	"github.com/aenthill/manifest"
	"github.com/spf13/afero"
)

func TestNoImageToAddError(t *testing.T) {
	err := &noImageToAddError{}
	if err.Error() != noImageToAddErrorMessage {
		t.Errorf("error returned a wrong message: got %s want %s", err.Error(), noImageToAddErrorMessage)
	}
}

func TestNewAddJob(t *testing.T) {
	t.Run("calling NewAddJob without images", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := &context.AppContext{ProjectDir: os.Getenv("HOST_PROJECT_DIR")}
		if _, err := NewAddJob(nil, m, ctx); err == nil {
			t.Error("NewAddJob should have thrown an error as there are no images in arguments")
		}
	})

	t.Run("calling NewAddJob with a non-existing manifest", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := &context.AppContext{ProjectDir: os.Getenv("HOST_PROJECT_DIR")}
		if _, err := NewAddJob([]string{"aent/foo"}, m, ctx); err == nil {
			t.Error("NewAddJob should have thrown an error as the given manifest should not exist")
		}
	})

	t.Run("calling NewAddJob with a broken manifest", func(t *testing.T) {
		m := manifest.New("../../tests/aenthill-broken.json", afero.NewOsFs())
		ctx := &context.AppContext{ProjectDir: os.Getenv("HOST_PROJECT_DIR")}
		if _, err := NewAddJob([]string{"aent/foo"}, m, ctx); err == nil {
			t.Error("NewAddJob should have thrown an error as the given manifest should be broken")
		}
	})

	t.Run("calling NewAddJob with a valid manifest", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := &context.AppContext{ProjectDir: os.Getenv("HOST_PROJECT_DIR")}
		if err := m.Flush(); err != nil {
			t.Errorf("an unexpected error occurred while flushing the given manifest: %s", err.Error())
		}
		if _, err := NewAddJob([]string{"aent/foo"}, m, ctx); err != nil {
			t.Error("NewAddJob should not have thrown an error as the given manifest should be valid")
		}
	})
}

func TestAddJobRun(t *testing.T) {
	t.Run("calling Run with a fake image", func(t *testing.T) {
		image := "aent/foo"
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := &context.AppContext{ProjectDir: os.Getenv("HOST_PROJECT_DIR")}
		if err := m.Flush(); err != nil {
			t.Errorf("an unexpected error occurred while flushing the given manifest: %s", err.Error())
		}
		job := &addJob{[]string{image}, m, ctx}
		if err := job.Run(); err == nil {
			t.Errorf("Run should have thrown an error as the image %s is invalid", image)
		}
	})

	t.Run("calling Run with a valid image", func(t *testing.T) {
		image := "aenthill/cassandra"
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := &context.AppContext{ProjectDir: os.Getenv("HOST_PROJECT_DIR")}
		if err := m.Flush(); err != nil {
			t.Errorf("an unexpected error occurred while flushing the given manifest: %s", err.Error())
		}
		job := &addJob{[]string{image}, m, ctx}
		if err := job.Run(); err != nil {
			t.Errorf("Run should not have thrown an error as the image %s is valid", image)
		}
	})
}

func TestAddJobHandle(t *testing.T) {
	t.Run("calling handle with a fake image", func(t *testing.T) {
		image := "aent/foo"
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := &context.AppContext{ProjectDir: os.Getenv("HOST_PROJECT_DIR")}
		if err := m.Flush(); err != nil {
			t.Errorf("an unexpected error occurred while flushing the given manifest: %s", err.Error())
		}
		job := &addJob{[]string{image}, m, ctx}
		if err := job.Run(); err == nil {
			t.Errorf("handle should have thrown an error as the image %s is invalid", image)
		}
	})

	t.Run("calling handle with a fake image which does exist in the manifest", func(t *testing.T) {
		image := "aent/foo"
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := &context.AppContext{ProjectDir: os.Getenv("HOST_PROJECT_DIR")}
		if err := m.AddAent(image); err != nil {
			t.Errorf("an unexpected error occurred while adding aent %s in the given manifest: %s", image, err.Error())
		}
		if err := m.Flush(); err != nil {
			t.Errorf("an unexpected error occurred while flushing the given manifest: %s", err.Error())
		}
		job := &addJob{[]string{image}, m, ctx}
		if err := job.Run(); err == nil {
			t.Errorf("handle should have thrown an error as the image %s is invalid", image)
		}
	})

	t.Run("calling handle with a valid image", func(t *testing.T) {
		image := "aenthill/cassandra"
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := &context.AppContext{ProjectDir: os.Getenv("HOST_PROJECT_DIR")}
		if err := m.Flush(); err != nil {
			t.Errorf("an unexpected error occurred while flushing the given manifest: %s", err.Error())
		}
		job := &addJob{[]string{image}, m, ctx}
		if err := job.Run(); err != nil {
			t.Errorf("handle should not have thrown an error as the image %s is valid", image)
		}
	})
}

func TestAddJobAddAent(t *testing.T) {
	t.Run("calling addAent with a non-existing image", func(t *testing.T) {
		image := "aent/foo"
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := &context.AppContext{ProjectDir: os.Getenv("HOST_PROJECT_DIR")}
		if err := m.Flush(); err != nil {
			t.Errorf("an unexpected error occurred while flushing the given manifest: %s", err.Error())
		}
		job := &addJob{[]string{image}, m, ctx}
		if _, err := job.addAent(image); err != nil {
			t.Error("addAent should not have thrown an error")
		}
	})

	t.Run("calling addAent with an existing image", func(t *testing.T) {
		image := "aent/foo"
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := &context.AppContext{ProjectDir: os.Getenv("HOST_PROJECT_DIR")}
		if err := m.AddAent(image); err != nil {
			t.Errorf("an unexpected error occurred while adding aent %s in the given manifest: %s", image, err.Error())
		}
		if err := m.Flush(); err != nil {
			t.Errorf("an unexpected error occurred while flushing the given manifest: %s", err.Error())
		}
		job := &addJob{[]string{image}, m, ctx}
		if _, err := job.addAent(image); err != nil {
			t.Error("addAent should not have thrown an error")
		}
	})
}

func TestEventAddFailedError(t *testing.T) {
	image := "aent/foo"
	err := errors.New("")
	eventFailedErr := eventAddFailedError{image, err}
	expected := fmt.Sprintf(eventAddFailedErrorMessage, image, err)
	if eventFailedErr.Error() != expected {
		t.Errorf("error returned a wrong message: got %s want %s", eventFailedErr.Error(), expected)
	}
}

func TestAddJobSendEvent(t *testing.T) {
	t.Run("calling sendEvent with a fake image", func(t *testing.T) {
		image := "aent/foo"
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := &context.AppContext{ProjectDir: os.Getenv("HOST_PROJECT_DIR")}
		job := &addJob{[]string{image}, m, ctx}
		if err := job.sendEvent(image); err == nil {
			t.Errorf("sendEvent should have thrown an error as the image %s is invalid", image)
		}
	})

	t.Run("calling sendEvent with a valid image", func(t *testing.T) {
		image := "aenthill/cassandra"
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := &context.AppContext{ProjectDir: os.Getenv("HOST_PROJECT_DIR")}
		job := &addJob{[]string{image}, m, ctx}
		if err := job.sendEvent(image); err != nil {
			t.Errorf("sendEvent should not have thrown an error as the image %s is valid", image)
		}
	})
}

func TestAddJobRemoveAent(t *testing.T) {
	t.Run("calling removeAent with a non-existing image", func(t *testing.T) {
		image := "aent/foo"
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := &context.AppContext{ProjectDir: os.Getenv("HOST_PROJECT_DIR")}
		if err := m.Flush(); err != nil {
			t.Errorf("an unexpected error occurred while flushing the given manifest: %s", err.Error())
		}
		job := &addJob{[]string{image}, m, ctx}
		if err := job.removeAent(image); err == nil {
			t.Errorf("removeAent should have thrown an error as the image %s should exist", image)
		}
	})

	t.Run("calling removeAent with an existing image", func(t *testing.T) {
		image := "aent/foo"
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := &context.AppContext{ProjectDir: os.Getenv("HOST_PROJECT_DIR")}
		if err := m.AddAent(image); err != nil {
			t.Errorf("an unexpected error occurred while adding aent %s in the given manifest: %s", image, err.Error())
		}
		if err := m.Flush(); err != nil {
			t.Errorf("an unexpected error occurred while flushing the given manifest: %s", err.Error())
		}
		job := &addJob{[]string{image}, m, ctx}
		if err := job.removeAent(image); err != nil {
			t.Errorf("removeAent should not have thrown an error as the image %s should not exist", image)
		}
	})
}
