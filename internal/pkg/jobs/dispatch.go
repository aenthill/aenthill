package jobs

import (
	"github.com/aenthill/aenthill/internal/pkg/context"
	"github.com/aenthill/aenthill/internal/pkg/docker"
	"github.com/aenthill/aenthill/internal/pkg/errors"
	"github.com/aenthill/aenthill/internal/pkg/manifest"
)

type dispatchJob struct {
	event    string
	payload  string
	filters  string
	docker   *docker.Docker
	ctx      *context.Context
	manifest *manifest.Manifest
}

// NewDispatchJob creates a new Job instance.
func NewDispatchJob(event, payload, filters string, ctx *context.Context, m *manifest.Manifest) (Job, error) {
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
	return &dispatchJob{event, payload, filters, d, ctx, m}, nil
}

func (j *dispatchJob) Execute() error {
	aents, err := j.manifest.Aents(j.event, j.filters)
	if err != nil {
		return errors.Wrap("dispatch job", err)
	}
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
