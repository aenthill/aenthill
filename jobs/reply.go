package jobs

import (
	"github.com/aenthill/aenthill/context"
	"github.com/aenthill/aenthill/docker"
	"github.com/aenthill/aenthill/errors"
	"github.com/aenthill/aenthill/manifest"
)

type replyJob struct {
	event   string
	payload string
	docker  *docker.Docker
}

// NewReplyJob creates a new Job instance.
func NewReplyJob(event, payload string, ctx *context.Context) (Job, error) {
	d, err := docker.New(ctx)
	if err != nil {
		return nil, errors.Wrap("reply job", err)
	}
	if !manifest.IsAlpha(event) {
		return nil, errors.Errorf("reply job", `"%s" is not a valid event name: only [A-Z0-9_] characters are authorized`, event)
	}
	return &replyJob{event, payload, d}, nil
}

func (j *replyJob) Execute() error {
	return errors.Wrap("reply job", j.docker.Reply(j.event, j.payload))
}
