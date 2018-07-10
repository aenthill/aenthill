package jobs

import (
	"github.com/aenthill/aenthill/context"
	"github.com/aenthill/aenthill/docker"
	"github.com/aenthill/aenthill/errors"
)

type replyJob struct {
	event   string
	payload string
	docker  *docker.Docker
}

func NewReplyJob(event, payload string, ctx *context.Context) (Job, error) {
	d, err := docker.New(ctx)
	if err != nil {
		return nil, errors.Wrap("reply job", err)
	}
	return &replyJob{event, payload, d}, nil
}

func (j *replyJob) Execute() error {
	return errors.Wrap("reply job", j.docker.Reply(j.event, j.payload))
}
