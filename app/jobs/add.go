package jobs

import (
	"fmt"

	"github.com/aenthill/aenthill/app/context"

	"github.com/aenthill/docker"
	"github.com/aenthill/log"
	"github.com/aenthill/manifest"
)

type addJob struct {
	images   []string
	manifest *manifest.Manifest
	appCtx   *context.AppContext
}

type noImageToAddError struct{}

const noImageToAddErrorMessage = "expecting at least one image to add"

func (e *noImageToAddError) Error() string {
	return noImageToAddErrorMessage
}

// NewAddJob creates an addJob instance.
// If given arguments are not valid, throws an error.
func NewAddJob(images []string, m *manifest.Manifest, appCtx *context.AppContext) (Job, error) {
	if len(images) == 0 {
		return nil, &noImageToAddError{}
	}

	if !m.Exist() {
		return nil, &manifestFileDoesNotExistError{}
	}

	if err := m.Parse(); err != nil {
		return nil, err
	}

	return &addJob{images, m, appCtx}, nil
}

// Run sends ADD event to all its images.
func (job *addJob) Run() error {
	for _, image := range job.images {
		if err := job.handle(image); err != nil {
			return err
		}
	}

	return nil
}

func (job *addJob) handle(image string) error {
	if err := job.addAent(image); err != nil {
		return err
	}

	if eventFailedErr := job.sendEvent(image); eventFailedErr != nil {
		if err := job.removeAent(image); err != nil {
			entryCtx := &log.EntryContext{Source: job.appCtx.Source}
			log.Errorf(entryCtx, err, "an unexpected error happened while removing aent %s from manifest %s", image, job.manifest.GetPath())
		}

		return eventFailedErr
	}

	return nil
}

func (job *addJob) addAent(image string) error {
	entryCtx := &log.EntryContext{Source: job.appCtx.Source}

	if err := job.manifest.AddAent(image); err != nil {
		log.Warnf(entryCtx, "aent %s does already exist in manifest %s", image, job.manifest.GetPath())
		return nil
	}

	if err := job.manifest.Flush(); err != nil {
		return err
	}

	log.Infof(entryCtx, "aent %s has been added to manifest %s", image, job.manifest.GetPath())

	return nil
}

type eventAddFailedError struct {
	image string
	err   error
}

const eventAddFailedErrorMessage = "event ADD sent to aent %s has failed: %s"

func (e *eventAddFailedError) Error() string {
	return fmt.Sprintf(eventAddFailedErrorMessage, e.image, e.err.Error())
}

func (job *addJob) sendEvent(image string) error {
	ctx := &docker.EventContext{
		Source:         job.appCtx.Source,
		To:             image,
		HostProjectDir: job.appCtx.ProjectDir,
		LogLevel:       job.appCtx.LogLevel,
	}

	if err := docker.Execute("ADD", "", ctx); err != nil {
		return &eventAddFailedError{image, err}
	}

	return nil
}

func (job *addJob) removeAent(image string) error {
	if err := job.manifest.RemoveAent(image); err != nil {
		return err
	}

	if err := job.manifest.Flush(); err != nil {
		return err
	}

	entryCtx := &log.EntryContext{Source: job.appCtx.Source}
	log.Warnf(entryCtx, "aent %s has been removed from manifest %s", image, job.manifest.GetPath())

	return nil
}