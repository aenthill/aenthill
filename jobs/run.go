package jobs

import (
	"strings"

	"github.com/aenthill/aenthill/context"
	"github.com/aenthill/aenthill/docker"
	"github.com/aenthill/aenthill/errors"
	"github.com/aenthill/aenthill/manifest"
)

type runJob struct {
	aent    *manifest.Aent
	key     string
	event   string
	payload string
	docker  *docker.Docker
}

// NewRunJob creates a new Job instance.
func NewRunJob(target, event, payload string, ctx *context.Context, m *manifest.Manifest) (Job, error) {
	d, err := docker.New(ctx)
	if err != nil {
		return nil, errors.Wrap("run job", err)
	}
	j := &runJob{event: strings.ToUpper(event), payload: payload, docker: d}
	aent, err := m.Aent(target)
	if err == nil {
		j.aent = aent
		j.key = target
		return j, nil
	}
	j.aent = &manifest.Aent{Image: target}
	return j, nil
}

func (j *runJob) Execute() error {
	return errors.Wrap("run job", j.docker.Run(j.key, j.aent.Image, j.event, j.payload, j.aent.Metadata, j.aent.Dependencies))
}
