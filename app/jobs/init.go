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

// Run implements Job.
func (job *initJob) Run() error {
	err := job.manifest.Flush()
	if err == nil {
		// set the log level to INFO to display the following message to user.
		log.SetLevel(log.InfoLevel)
		log.Infof(job.appCtx.EntryContext, "%s created! May the swarm be with you", job.manifest.GetPath())
		// reset the log level to the level defined by the user.
		log.SetLevel(job.appCtx.LogLevel)
	}

	return err
}
