package jobs

import (
	"errors"
	"fmt"
	"testing"

	"github.com/aenthill/aenthill/tests"

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
		ctx := tests.NewAppContext()
		if _, err := NewRemoveJob(nil, m, ctx); err == nil {
			t.Error("NewRemoveJob should have thrown an error as there are no images in arguments")
		}
	})

	t.Run("calling NewRemoveJob with a non-existing manifest", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := tests.NewAppContext()
		if _, err := NewRemoveJob([]string{"aent/foo"}, m, ctx); err == nil {
			t.Error("NewRemoveJob should have thrown an error as the given manifest should not exist")
		}
	})

	t.Run("calling NewRemoveJob with a broken manifest", func(t *testing.T) {
		m := manifest.New("../../tests/aenthill-broken.json", afero.NewOsFs())
		ctx := tests.NewAppContext()
		if _, err := NewRemoveJob([]string{"aent/foo"}, m, ctx); err == nil {
			t.Error("NewRemoveJob should have thrown an error as the given manifest should be broken")
		}
	})

	t.Run("calling NewRemoveJob with a valid manifest", func(t *testing.T) {
		m, err := tests.NewEmptyInMemoryManifest()
		if err != nil {
			t.Error(err)
		}
		ctx := tests.NewAppContext()
		if _, err := NewRemoveJob([]string{"aent/foo"}, m, ctx); err != nil {
			t.Error("NewRemoveJob should not have thrown an error as the given manifest should be valid")
		}
	})
}

func TestRemoveJobRun(t *testing.T) {
	t.Run("calling Run with a fake image", func(t *testing.T) {
		image := "aent/bar"
		m, err := tests.NewEmptyInMemoryManifest()
		if err != nil {
			t.Error(err)
		}
		ctx := tests.NewAppContext()
		job := &removeJob{[]string{image}, m, ctx}
		if err := job.Run(); err == nil {
			t.Errorf("Run should have thrown an error as the image %s is invalid", image)
		}
	})

	t.Run("calling run with a fake image which does exist in the manifest", func(t *testing.T) {
		image := "aent/bar"
		m, err := tests.NewInMemoryManifestWithFakeImage()
		if err != nil {
			t.Error(err)
		}
		ctx := tests.NewAppContext()
		job := &removeJob{[]string{image}, m, ctx}
		if err := job.Run(); err == nil {
			t.Errorf("Run should have thrown an error as the image %s is invalid", image)
		}
	})

	t.Run("calling Run with a valid image", func(t *testing.T) {
		image := "aenthill/cassandra"
		m, err := tests.NewEmptyInMemoryManifest()
		if err != nil {
			t.Error(err)
		}
		ctx := tests.NewAppContext()
		job := &removeJob{[]string{image}, m, ctx}
		if err := job.Run(); err != nil {
			t.Errorf("Run should not have thrown an error as the image %s is valid", image)
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

func TestReAddAent(t *testing.T) {
	t.Run("calling reAddAent with an image which does exist in the manifest", func(t *testing.T) {
		image := "aent/bar"
		m, err := tests.NewInMemoryManifestWithFakeImage()
		if err != nil {
			t.Error(err)
		}
		ctx := tests.NewAppContext()
		job := &removeJob{[]string{image}, m, ctx}
		if err := job.reAddAent(image); err == nil {
			t.Errorf("reAddAent should have thrown an error as the image %s should exist in given manifest", image)
		}
	})
}
