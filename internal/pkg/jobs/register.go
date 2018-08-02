package jobs

import (
	"github.com/aenthill/aenthill/internal/pkg/context"
	"github.com/aenthill/aenthill/internal/pkg/errors"
	"github.com/aenthill/aenthill/internal/pkg/log"
	"github.com/aenthill/aenthill/internal/pkg/manifest"
)

type registerJob struct {
	image    string
	key      string
	metadata []string
	ctx      *context.Context
	manifest *manifest.Manifest
}

// NewRegisterJob creates a new Job instance.
func NewRegisterJob(image, key string, metadata []string, ctx *context.Context, m *manifest.Manifest) (Job, error) {
	if err := m.Parse(); err != nil {
		return nil, err
	}
	return &registerJob{image, key, metadata, ctx, m}, nil
}

func (j *registerJob) Execute() error {
	log.Infof(`adding "%s" with key "%s" as dependency of "%s" identified as "%s" in manifest`, j.image, j.key, j.ctx.Image, j.ctx.ID)
	ID, err := j.manifest.AddDependency(j.ctx.ID, j.image, j.key)
	if err != nil {
		return errors.Wrap("register job", err)
	}
	if err := j.manifest.AddMetadataFromFlags(ID, j.metadata); err != nil {
		return errors.Wrap("register job", err)
	}
	return errors.Wrap("register job", j.manifest.Flush())
}
