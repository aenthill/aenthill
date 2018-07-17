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
	manifest *manifest.Manifest
}

// NewDispatchJob creates a new Job instance.
func NewDispatchJob(event, payload string, ctx *context.Context, m *manifest.Manifest) (Job, error) {
	d, err := docker.New(ctx)
	if err != nil {
		return nil, errors.Wrap("dispatch job", err)
	}
	return &dispatchJob{event, payload, d, m}, nil
}

func (j *dispatchJob) Execute() error {
	aents := j.manifest.Aents(j.event)
	for key, aent := range aents {
		if err := j.docker.Run(key, aent.Image, j.event, j.payload, aent.Metadata, aent.Dependencies); err != nil {
			return errors.Wrap("dispatch job", err)
		}
	}
	return nil
}
