package jobs

import (
	"strings"

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
func NewRegisterJob(image, dependencyKey string, metadata []string, ctx *context.Context, m *manifest.Manifest) Job {
	return &registerJob{image, dependencyKey, metadata, ctx, m}
}

func (j *registerJob) Execute() error {
	key, err := j.manifest.AddDependency(j.ctx.Key, j.image, j.dependencyKey)
	if err != nil {
		return err
	}
	if err := j.handleMetadata(key); err != nil {
		return err
	}
	if err := j.ctx.PopulateEnv(j.manifest); err != nil {
		return errors.Wrap("register job", err)
	}
	return errors.Wrap("register job", j.manifest.Flush())
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
