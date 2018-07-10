package jobs

import (
	"fmt"
	"os"
	"strings"

	"github.com/aenthill/aenthill/context"
	"github.com/aenthill/aenthill/errors"
	"github.com/aenthill/aenthill/manifest"
)

type registerJob struct {
	image    string
	envVar   string
	metadata []string
	events   []string
	ctx      *context.Context
	manifest *manifest.Manifest
}

func NewRegisterJob(image, envVar string, metadata, events []string, ctx *context.Context, m *manifest.Manifest) Job {
	return &registerJob{image, envVar, metadata, events, ctx, m}
}

func (j *registerJob) Execute() error {
	key := j.manifest.AddAent(j.image)
	if err := j.handleEvents(key); err != nil {
		return err
	}
	if err := j.handleMetadata(key); err != nil {
		return err
	}
	if j.envVar != "" {
		metadata := make(map[string]string)
		metadata[j.envVar] = key
		if err := j.manifest.AddMetadata(j.ctx.Key, metadata); err != nil {
			return errors.Wrap("register job", err)
		}
		if err := os.Setenv(fmt.Sprintf("PHEROMONE_METADATA_%s", strings.ToUpper(j.envVar)), key); err != nil {
			return errors.Wrap("register job", err)
		}
	}
	return errors.Wrap("register job", j.manifest.Flush())
}

func (j *registerJob) handleEvents(key string) error {
	return errors.Wrap("register job", j.manifest.AddEvents(key, j.events...))
}

func (j *registerJob) handleMetadata(key string) error {
	if j.metadata == nil {
		return nil
	}
	metadata := make(map[string]string)
	for _, data := range j.metadata {
		parts := strings.Split(data, "=")
		if len(parts) != 2 {
			return errors.Errorf("register job", `execute: wrong metadata format: got "%s" want "key=value"`, data)
		}
		metadata[parts[0]] = parts[1]
	}
	return errors.Wrap("register job", j.manifest.AddMetadata(key, metadata))
}
