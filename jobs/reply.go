package jobs

import (
	"strings"

	"github.com/aenthill/aenthill/context"
	"github.com/aenthill/aenthill/docker"
	"github.com/aenthill/aenthill/errors"
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
	return &replyJob{strings.ToUpper(event), payload, d}, nil
}

func (j *replyJob) Execute() error {
	return errors.Wrap("reply job", j.docker.Reply(j.event, j.payload))
}
