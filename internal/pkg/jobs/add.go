package jobs

import (
	"github.com/aenthill/aenthill/internal/pkg/context"
	"github.com/aenthill/aenthill/internal/pkg/docker"
	"github.com/aenthill/aenthill/internal/pkg/errors"
	"github.com/aenthill/aenthill/internal/pkg/log"
	"github.com/aenthill/aenthill/internal/pkg/manifest"
)

type addJob struct {
	image    string
	docker   *docker.Docker
	ctx      *context.Context
	manifest *manifest.Manifest
}

// NewAddJob creates a new Job instance.
func NewAddJob(image string, ctx *context.Context, m *manifest.Manifest) (Job, error) {
	if err := m.ParseIfExist(); err != nil {
		return nil, errors.Wrap("add job", err)
	}
	d, err := docker.New(ctx)
	if err != nil {
		return nil, errors.Wrap("add job", err)
	}
	return &addJob{image, d, ctx, m}, nil
}

func (j *addJob) Execute() error {
	ID := j.manifest.AddAent(j.image)
	if err := j.manifest.Flush(); err != nil {
		return errors.Wrap("add job", err)
	}
	return errors.Wrap("add job", j.sendEvent(ID))
}

func (j *addJob) sendEvent(ID string) error {
	runErr := j.docker.Run(j.image, ID, "ADD", "")
	if runErr == nil {
		return nil
	}
	log.Warnf(`event "%s" failed, removing "%s" identified as "%s" from manifest`, "ADD", j.image, ID)
	if err := j.manifest.RemoveAent(ID); err != nil {
		return err
	}
	if err := j.manifest.Flush(); err != nil {
		return err
	}
	return runErr
}
