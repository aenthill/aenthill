package jobs

import (
	"github.com/aenthill/aenthill/context"
	"github.com/aenthill/aenthill/errors"
	"github.com/aenthill/aenthill/manifest"
)

type addJob struct {
	image    string
	ctx      *context.Context
	manifest *manifest.Manifest
}

// NewAddJob creates a new Job instance.
func NewAddJob(image string, ctx *context.Context, m *manifest.Manifest) (Job, error) {
	if err := m.ParseIfExist(); err != nil {
		return nil, errors.Wrap("add job", err)
	}
	return &addJob{image, ctx, m}, nil
}

func (j *addJob) Execute() error {
	key := j.manifest.AddAent(j.image)
	if err := j.manifest.Flush(); err != nil {
		return errors.Wrap("add job", err)
	}
	return errors.Wrap("add job", j.sendEvent(key))
}

func (j *addJob) sendEvent(key string) error {
	runJob, err := NewRunJob(key, "ADD", "", j.ctx, j.manifest)
	if err != nil {
		return err
	}
	if err := runJob.Execute(); err == nil {
		return nil
	}
	return j.manifest.RemoveAent(key)
}
