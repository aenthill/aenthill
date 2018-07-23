package jobs

import (
	"github.com/aenthill/aenthill/context"
	"github.com/aenthill/aenthill/docker"
	"github.com/aenthill/aenthill/errors"
	"github.com/aenthill/aenthill/manifest"
)

type dispatchJob struct {
	event    string
	payload  string
	docker   *docker.Docker
	ctx      *context.Context
	manifest *manifest.Manifest
}

// NewDispatchJob creates a new Job instance.
func NewDispatchJob(event, payload string, ctx *context.Context, m *manifest.Manifest) (Job, error) {
	if err := m.Validate(event, "event"); err != nil {
		return nil, errors.Wrap("dispatch job", err)
	}
	if err := m.Parse(); err != nil {
		return nil, errors.Wrap("dispatch job", err)
	}
	d, err := docker.New(ctx)
	if err != nil {
		return nil, errors.Wrap("dispatch job", err)
	}
	return &dispatchJob{event, payload, d, ctx, m}, nil
}

func (j *dispatchJob) Execute() error {
	aents := j.manifest.Aents(j.event)
	for ID, aent := range aents {
		if err := j.sendEvent(aent.Image, ID); err != nil {
			return errors.Wrap("dispatch job", err)
		}
	}
	return nil
}

func (j *dispatchJob) sendEvent(image, ID string) error {
	if ID == j.ctx.ID {
		return nil
	}
	return j.docker.Run(image, ID, j.event, j.payload)
}
