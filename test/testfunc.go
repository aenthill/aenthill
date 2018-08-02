// Package test contains useful functions used across tests.
package test

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/aenthill/aenthill/internal/pkg/context"
)

// Context creates a new context for tests.
func Context(t *testing.T) *context.Context {
	ctx, err := context.New()
	if err != nil {
		t.Fatalf(`An unexpected error occurred while creating the context: got "%s"`, err.Error())
	}
	ctx.HostProjectDir = os.Getenv("HOST_PROJECT_DIR")
	return ctx
}

// BrokenManifestAbsPath returns the absolute path of the broken manifest file.
func BrokenManifestAbsPath(t *testing.T) string {
	return abs(t, "aenthill-broken.json")
}

// ValidManifestAbsPath returns the absolute path of the valid manifest file.
func ValidManifestAbsPath(t *testing.T) string {
	return abs(t, "aenthill.json")
}

func abs(t *testing.T, filename string) string {
	_, gofilename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf(`Got no caller information while trying to get the absolute path of manifest "%s"`, filename)
	}
	path, err := filepath.Abs(fmt.Sprintf("%s/testdata/%s", path.Dir(gofilename), filename))
	if err != nil {
		t.Fatalf(`An unexpected error occurred while getting the absolute path of the manifest "%s": got "%s"`, filename, err.Error())
	}
	return path
}
