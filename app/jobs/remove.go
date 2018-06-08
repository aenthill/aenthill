package jobs

import (
	"fmt"
	"time"

	"github.com/aenthill/aenthill/app/context"

	"github.com/aenthill/docker"
	"github.com/aenthill/log"
	"github.com/aenthill/manifest"
)

type removeJob struct {
	images   []string
	manifest *manifest.Manifest
	appCtx   *context.AppContext
}

type noImageToRemoveError struct{}

const noImageToRemoveErrorMessage = "expecting at least one image to remove"

func (e *noImageToRemoveError) Error() string {
	return noImageToRemoveErrorMessage
}

// NewRemoveJob creates a removeJob instance.
// If given arguments are not valid, throws an error.
func NewRemoveJob(images []string, m *manifest.Manifest, appCtx *context.AppContext) (Job, error) {
	if len(images) == 0 {
		return nil, &noImageToRemoveError{}
	}

	if !m.Exist() {
		return nil, &manifestFileDoesNotExistError{}
	}

	if err := m.Parse(); err != nil {
		return nil, err
	}

	return &removeJob{images, m, appCtx}, nil
}

// Run implements Job.
func (job *removeJob) Run() error {
	start := time.Now()

	err := job.run()
	if err != nil {
		log.Errorf(job.appCtx.EntryContext, err, "job has failed after %0.2fs", time.Since(start).Seconds())
	} else {
		log.Infof(job.appCtx.EntryContext, "job has successfully finished after %0.2fs", time.Since(start).Seconds())
	}

	return err
}

func (job *removeJob) run() error {
	for _, image := range job.images {
		if err := job.handle(image); err != nil {
			return err
		}
	}

	return nil
}

func (job *removeJob) handle(image string) error {
	aentExist, err := job.removeAent(image)
	if err != nil {
		return err
	}

	if eventFailedErr := job.sendEvent(image); eventFailedErr != nil {
		if !aentExist {
			log.Warnf(job.appCtx.EntryContext, "aent %s has not been re-added to manifest %s as it was not existing previously", image, job.manifest.GetPath())
			return eventFailedErr
		}

		if err := job.reAddAent(image); err != nil {
			log.Errorf(job.appCtx.EntryContext, err, "an unexpected error happened while re-adding aent %s to manifest %s", image, job.manifest.GetPath())
		}

		return eventFailedErr
	}

	return nil
}

func (job *removeJob) removeAent(image string) (bool, error) {
	if err := job.manifest.RemoveAent(image); err != nil {
		log.Warnf(job.appCtx.EntryContext, "aent %s does not exist in manifest %s", image, job.manifest.GetPath())
		return false, nil
	}

	if err := job.manifest.Flush(); err != nil {
		return true, err
	}

	log.Infof(job.appCtx.EntryContext, "aent %s has been removed from manifest %s", image, job.manifest.GetPath())

	return true, nil
}

type eventRemoveFailedError struct {
	image string
	err   error
}

const eventRemoveFailedErrorMessage = "event REMOVE sent to aent %s has failed: %s"

func (e *eventRemoveFailedError) Error() string {
	return fmt.Sprintf(eventRemoveFailedErrorMessage, e.image, e.err.Error())
}

func (job *removeJob) sendEvent(image string) error {
	ctx := &docker.EventContext{
		Source:         job.appCtx.EntryContext.Source,
		To:             image,
		HostProjectDir: job.appCtx.ProjectDir,
		LogLevel:       job.appCtx.LogLevel,
	}

	if err := docker.Execute("REMOVE", "", ctx); err != nil {
		return &eventRemoveFailedError{image, err}
	}

	return nil
}

func (job *removeJob) reAddAent(image string) error {
	if err := job.manifest.AddAent(image); err != nil {
		return err
	}

	if err := job.manifest.Flush(); err != nil {
		return err
	}

	log.Warnf(job.appCtx.EntryContext, "aent %s has been re-added to manifest %s", image, job.manifest.GetPath())

	return nil
}
