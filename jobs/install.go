package jobs

import (
	"fmt"
	"os"
	"strings"

	"github.com/aenthill/aenthill/context"
	"github.com/aenthill/aenthill/errors"
	"github.com/aenthill/aenthill/manifest"
)

type installJob struct {
	metadata []string
	events   []string
	ctx      *context.Context
	manifest *manifest.Manifest
}

func NewInstallJob(metadata, events []string, ctx *context.Context, m *manifest.Manifest) Job {
	return &installJob{metadata, events, ctx, m}
}

func (j *installJob) Execute() error {
	if j.ctx.Key == "" {
		j.ctx.Key = j.manifest.AddAent(j.ctx.Image)
	}
	if err := j.handleEvents(); err != nil {
		return err
	}
	if err := j.handleMetadata(); err != nil {
		return err
	}
	if err := os.Setenv(context.KeyEnvVar, j.ctx.Key); err != nil {
		return errors.Wrap("install job", err)
	}
	return errors.Wrap("install job", j.manifest.Flush())
}

func (j *installJob) handleEvents() error {
	return errors.Wrap("install job", j.manifest.AddEvents(j.ctx.Key, j.events...))
}

func (j *installJob) handleMetadata() error {
	if j.metadata == nil {
		return nil
	}
	metadata := make(map[string]string)
	for _, data := range j.metadata {
		parts := strings.Split(data, "=")
		if len(parts) != 2 {
			return errors.Errorf("install job", `execute: wrong metadata format: got "%s" want "key=value"`, data)
		}
		metadata[parts[0]] = parts[1]
		if err := os.Setenv(fmt.Sprintf("PHEROMONE_DATA_%s", strings.ToUpper(parts[0])), parts[1]); err != nil {
			return errors.Wrap("install job", err)
		}
	}
	return errors.Wrap("install job", j.manifest.AddMetadata(j.ctx.Key, metadata))
}
