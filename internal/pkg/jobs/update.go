package jobs

import (
	"strings"

	"github.com/aenthill/aenthill/internal/pkg/context"
	"github.com/aenthill/aenthill/internal/pkg/errors"
	"github.com/aenthill/aenthill/internal/pkg/log"
	"github.com/aenthill/aenthill/internal/pkg/manifest"
)

type updateJob struct {
	metadata []string
	events   []string
	ctx      *context.Context
	manifest *manifest.Manifest
}

// NewUpdateJob creates a new Job instance.
func NewUpdateJob(metadata, events []string, ctx *context.Context, m *manifest.Manifest) (Job, error) {
	if err := m.Parse(); err != nil {
		return nil, errors.Wrap("update job", err)
	}
	if ctx.ID == "" {
		return nil, errors.Error("update job", "aent ID is missing")
	}
	return &updateJob{metadata, events, ctx, m}, nil
}

func (j *updateJob) Execute() error {
	log.Infof(`updating "%s" identified as "%s" in manifest`, j.ctx.Image, j.ctx.ID)
	log.Debugf(`events = "%s"`, strings.Join(j.events, ", "))
	log.Debugf(`metadata = "%s"`, strings.Join(j.metadata, ", "))
	if err := j.manifest.AddEvents(j.ctx.ID, j.events...); err != nil {
		return errors.Wrap("update job", err)
	}
	if err := j.manifest.AddMetadataFromFlags(j.ctx.ID, j.metadata); err != nil {
		return errors.Wrap("update job", err)
	}
	return errors.Wrap("update job", j.manifest.Flush())
}
