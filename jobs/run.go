package jobs

import (
	"github.com/aenthill/aenthill/context"
	"github.com/aenthill/aenthill/docker"
	"github.com/aenthill/aenthill/errors"
	"github.com/aenthill/aenthill/manifest"
)

type runJob struct {
	image   string
	key     string
	event   string
	payload string
	docker  *docker.Docker
}

// NewRunJob creates a new Job instance.
func NewRunJob(target, event, payload string, ctx *context.Context, m *manifest.Manifest) (Job, error) {
	if err := m.Validate(event, "event"); err != nil {
		return nil, errors.Wrap("run job", err)
	}
	if err := m.ParseIfExist(); err != nil {
		return nil, errors.Wrap("run job", err)
	}
	d, err := docker.New(ctx)
	if err != nil {
		return nil, errors.Wrap("run job", err)
	}
	image, key := func(target string, m *manifest.Manifest) (string, string) {
		aent, err := m.Aent(target)
		if err == nil {
			return aent.Image, target
		}
		return target, ""
	}(target, m)
	j := &runJob{image, key, event, payload, d}
	return j, nil
}

func (j *runJob) Execute() error {
	return errors.Wrap("run job", j.docker.Run(j.image, j.key, j.event, j.payload))
}
