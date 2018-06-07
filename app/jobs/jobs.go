// Package jobs provides core logic of the commands of the application.
package jobs

import (
	"fmt"

	"github.com/aenthill/manifest"
)

// Job is a job interface.
type Job interface {
	Run() error
}

type manifestFileDoesNotExistError struct{}

const manifestFileDoesNotExistErrorMessage = "manifest %s not found in current directory. Did you run aenthill init?"

func (e *manifestFileDoesNotExistError) Error() string {
	return fmt.Sprintf(manifestFileDoesNotExistErrorMessage, manifest.DefaultManifestFileName)
}
