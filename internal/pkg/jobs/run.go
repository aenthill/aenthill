package jobs

import (
	"github.com/aenthill/aenthill/internal/pkg/context"
	"github.com/aenthill/aenthill/internal/pkg/docker"
	"github.com/aenthill/aenthill/internal/pkg/errors"
	"github.com/aenthill/aenthill/internal/pkg/manifest"
)

type runJob struct {
	image   string
	ID      string
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
	image, ID := func(ID, target string, m *manifest.Manifest) (string, string) {
		aent, err := m.Dependency(ID, target)
		if err == nil {
			return aent.Image, target
		}
		return target, ""
	}(ctx.ID, target, m)
	j := &runJob{image, ID, event, payload, d}
	return j, nil
}

func (j *runJob) Execute() error {
	return errors.Wrap("run job", j.docker.Run(j.image, j.ID, j.event, j.payload))
}
