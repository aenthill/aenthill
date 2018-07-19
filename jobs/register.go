package jobs

import (
	"github.com/aenthill/aenthill/context"
	"github.com/aenthill/aenthill/errors"
	"github.com/aenthill/aenthill/manifest"
)

type registerJob struct {
	image         string
	dependencyKey string
	metadata      []string
	ctx           *context.Context
	manifest      *manifest.Manifest
}

// NewRegisterJob creates a new Job instance.
func NewRegisterJob(image, dependencyKey string, metadata []string, ctx *context.Context, m *manifest.Manifest) (Job, error) {
	if err := m.Parse(); err != nil {
		return nil, err
	}
	return &registerJob{image, dependencyKey, metadata, ctx, m}, nil
}

func (j *registerJob) Execute() error {
	key, err := j.manifest.AddDependency(j.ctx.Key, j.image, j.dependencyKey)
	if err != nil {
		return errors.Wrap("register job", err)
	}
	if err := j.manifest.AddMetadataFromFlags(key, j.metadata); err != nil {
		return errors.Wrap("register job", err)
	}
	return errors.Wrap("register job", j.manifest.Flush())
}
