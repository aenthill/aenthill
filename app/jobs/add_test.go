package jobs

import (
	"errors"
	"fmt"
	"testing"

	"github.com/aenthill/aenthill/tests"

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
		ctx := tests.NewAppContext()
		if _, err := NewAddJob(nil, m, ctx); err == nil {
			t.Error("NewAddJob should have thrown an error as there are no images in arguments")
		}
	})

	t.Run("calling NewAddJob with a non-existing manifest", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := tests.NewAppContext()
		if _, err := NewAddJob([]string{"aent/foo"}, m, ctx); err == nil {
			t.Error("NewAddJob should have thrown an error as the given manifest should not exist")
		}
	})

	t.Run("calling NewAddJob with a broken manifest", func(t *testing.T) {
		m := manifest.New("../../tests/aenthill-broken.json", afero.NewOsFs())
		ctx := tests.NewAppContext()
		if _, err := NewAddJob([]string{"aent/foo"}, m, ctx); err == nil {
			t.Error("NewAddJob should have thrown an error as the given manifest should be broken")
		}
	})

	t.Run("calling NewAddJob with a valid manifest", func(t *testing.T) {
		m, err := tests.NewEmptyInMemoryManifest()
		if err != nil {
			t.Error(err)
		}
		ctx := tests.NewAppContext()
		if _, err := NewAddJob([]string{"aent/foo"}, m, ctx); err != nil {
			t.Error("NewAddJob should not have thrown an error as the given manifest should be valid")
		}
	})
}

func TestAddJobRun(t *testing.T) {
	t.Run("calling Run with a fake image", func(t *testing.T) {
		image := "aent/foo"
		m, err := tests.NewEmptyInMemoryManifest()
		if err != nil {
			t.Error(err)
		}
		ctx := tests.NewAppContext()
		job := &addJob{[]string{image}, m, ctx}
		if err := job.Run(); err == nil {
			t.Errorf("Run should have thrown an error as the image %s is invalid", image)
		}
	})

	t.Run("calling Run with a fake image which does exist in the manifest", func(t *testing.T) {
		image := "aent/bar"
		m, err := tests.NewInMemoryManifestWithFakeImage()
		if err != nil {
			t.Error(err)
		}
		ctx := tests.NewAppContext()
		job := &addJob{[]string{image}, m, ctx}
		if err := job.Run(); err == nil {
			t.Errorf("Run should not have thrown an error as the image %s is valid", image)
		}
	})

	t.Run("calling Run with a valid image", func(t *testing.T) {
		image := "aenthill/cassandra"
		m, err := tests.NewEmptyInMemoryManifest()
		if err != nil {
			t.Error(err)
		}
		ctx := tests.NewAppContext()
		if err := m.Flush(); err != nil {
			t.Errorf("an unexpected error occurred while flushing the given manifest: %s", err.Error())
		}
		job := &addJob{[]string{image}, m, ctx}
		if err := job.Run(); err != nil {
			t.Errorf("Run should not have thrown an error as the image %s is valid", image)
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
