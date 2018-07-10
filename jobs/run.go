package jobs

import (
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

func NewRunJob(target, event, payload string, ctx *context.Context, m *manifest.Manifest) (Job, error) {
	d, err := docker.New(ctx)
	if err != nil {
		return nil, errors.Wrap("run job", err)
	}
	j := &runJob{event: event, payload: payload, docker: d}
	j.aent = &manifest.Aent{}
	for key, aent := range m.Aents("") {
		if key == target {
			j.aent = aent
			j.key = key
			break
		}
	}
	return j, nil
}

func (j *runJob) Execute() error {
	return errors.Wrap("run job", j.docker.Run(j.key, j.aent.Image, j.event, j.payload, j.aent.Metadata))
}
