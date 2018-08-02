package jobs

import (
	"github.com/aenthill/aenthill/internal/pkg/context"
	"github.com/aenthill/aenthill/internal/pkg/docker"
	"github.com/aenthill/aenthill/internal/pkg/errors"
	"github.com/aenthill/aenthill/internal/pkg/manifest"
)

type replyJob struct {
	event   string
	payload string
	docker  *docker.Docker
}

// NewReplyJob creates a new Job instance.
func NewReplyJob(event, payload string, ctx *context.Context, m *manifest.Manifest) (Job, error) {
	if err := m.Validate(event, "event"); err != nil {
		return nil, errors.Wrap("reply job", err)
	}
	d, err := docker.New(ctx)
	if err != nil {
		return nil, errors.Wrap("reply job", err)
	}
	return &replyJob{event, payload, d}, nil
}

func (j *replyJob) Execute() error {
	return errors.Wrap("reply job", j.docker.Reply(j.event, j.payload))
}
