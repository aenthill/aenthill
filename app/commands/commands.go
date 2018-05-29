// Package commands contains all commands of the application.
package commands

import (
	"fmt"

	"github.com/aenthill/manifest"
)

type manifestFileDoesNotExistError struct{}

const manifestFileDoesNotExistErrorMessage = "manifest %s not found in current directory. Did you run aenthill init?"

func (e *manifestFileDoesNotExistError) Error() string {
	return fmt.Sprintf(manifestFileDoesNotExistErrorMessage, manifest.DefaultManifestFileName)
}
