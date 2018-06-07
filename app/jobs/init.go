package jobs

import (
	"fmt"

	"github.com/aenthill/aenthill/app/context"

	"github.com/aenthill/log"
	"github.com/aenthill/manifest"
)

type initJob struct {
	manifest *manifest.Manifest
	appCtx   *context.AppContext
}

type manifestFileAlreadyExistingError struct {
	existingPath string
}

const manifestFileAlreadyExistingErrorMessage = "manifest %s already exists"

func (e *manifestFileAlreadyExistingError) Error() string {
	return fmt.Sprintf(manifestFileAlreadyExistingErrorMessage, e.existingPath)
}

// NewInitJob creates an initJob instance.
// If given arguments are not valid, throws an error.
func NewInitJob(m *manifest.Manifest, appCtx *context.AppContext) (Job, error) {
	if m.Exist() {
		return nil, &manifestFileAlreadyExistingError{m.GetPath()}
	}

	return &initJob{m, appCtx}, nil
}

// Run creates the manifest file if it does not exist.
func (job *initJob) Run() error {
	err := job.manifest.Flush()
	if err == nil {
		entryCtx := &log.EntryContext{Source: job.appCtx.Source}
		log.Infof(entryCtx, "%s created! May the swarm be with you", job.manifest.GetPath())
	}

	return err
}
