package jobs

import (
	"strings"

	"github.com/aenthill/aenthill/context"
	"github.com/aenthill/aenthill/errors"
	"github.com/aenthill/aenthill/log"
	"github.com/aenthill/aenthill/manifest"
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
	if ctx.Key == "" {
		return nil, errors.Error("update job", "aent key is missing")
	}
	return &updateJob{metadata, events, ctx, m}, nil
}

func (j *updateJob) Execute() error {
	log.Infof(`updating "%s" identified as "%s" in manifest`, j.ctx.Image, j.ctx.Key)
	log.Debugf(`events = "%s"`, strings.Join(j.events, ", "))
	log.Debugf(`metadata = "%s"`, strings.Join(j.metadata, ", "))
	if err := j.manifest.AddEvents(j.ctx.Key, j.events...); err != nil {
		return errors.Wrap("update job", err)
	}
	if err := j.manifest.AddMetadataFromFlags(j.ctx.Key, j.metadata); err != nil {
		return errors.Wrap("update job", err)
	}
	return errors.Wrap("update job", j.manifest.Flush())
}
