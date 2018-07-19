package jobs

import (
	"github.com/aenthill/aenthill/context"
	"github.com/aenthill/aenthill/docker"
	"github.com/aenthill/aenthill/errors"
	"github.com/aenthill/aenthill/log"
	"github.com/aenthill/aenthill/manifest"
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
	key := j.manifest.AddAent(j.image)
	if err := j.manifest.Flush(); err != nil {
		return errors.Wrap("add job", err)
	}
	return errors.Wrap("add job", j.sendEvent(key))
}

func (j *addJob) sendEvent(key string) error {
	runErr := j.docker.Run(j.image, key, "ADD", "")
	if runErr == nil {
		return nil
	}
	log.Warnf(`event "%s" failed, removing "%s" identified as "%s" from manifest`, "ADD", j.image, key)
	if err := j.manifest.RemoveAent(key); err != nil {
		return err
	}
	if err := j.manifest.Flush(); err != nil {
		return err
	}
	return runErr
}
