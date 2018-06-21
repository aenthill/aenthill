// Package tests contains functions used across our tests.
package tests

import (
	"fmt"
	"os"

	"github.com/aenthill/aenthill/app/context"

	"github.com/aenthill/log"
	"github.com/aenthill/manifest"
	"github.com/spf13/afero"
)

// NewAppContext creates an AppContext for a test.
func NewAppContext() *context.AppContext {
	return &context.AppContext{
		ProjectDir:    os.Getenv("HOST_PROJECT_DIR"),
		IsVerbose:     false,
		IsVeryVerbose: true,
		LogLevel:      "DEBUG",
		EntryContext:  &log.EntryContext{},
	}
}

// NewEmptyInMemoryManifest creates a in memory Manifest without images.
func NewEmptyInMemoryManifest() (*manifest.Manifest, error) {
	m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
	if err := m.Flush(); err != nil {
		return nil, fmt.Errorf("an unexpected error occurred while flushing the given manifest: %s", err.Error())
	}
	return m, nil
}

// NewInMemoryManifestWithFakeImage creates a in memory Manifest with a fake image.
func NewInMemoryManifestWithFakeImage() (*manifest.Manifest, error) {
	image := "aent/bar"
	m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
	if err := m.AddAent(image); err != nil {
		return nil, fmt.Errorf("an unexpected error occurred while adding aent %s in the given manifest: %s", image, err.Error())
	}
	if err := m.Flush(); err != nil {
		return nil, fmt.Errorf("an unexpected error occurred while flushing the given manifest: %s", err.Error())
	}
	return m, nil
}
