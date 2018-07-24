package jobs

import (
	"fmt"

	"github.com/aenthill/aenthill/context"
	"github.com/aenthill/aenthill/errors"
	"github.com/aenthill/aenthill/manifest"
)

type metadataJob struct {
	key      string
	ctx      *context.Context
	manifest *manifest.Manifest
}

// NewMetadataJob creates a new Job instance.
func NewMetadataJob(key string, ctx *context.Context, m *manifest.Manifest) (Job, error) {
	if err := m.Parse(); err != nil {
		return nil, errors.Wrap("metatada job", err)
	}
	if ctx.ID == "" {
		return nil, errors.Error("metatada job", "aent ID is missing")
	}
	return &metadataJob{key, ctx, m}, nil
}

func (j *metadataJob) Execute() error {
	metadata, err := j.manifest.Metadata(j.ctx.ID)
	if err != nil {
		return errors.Wrap("metadata job", err)
	}
	if j.key == "" {
		j.printAll(metadata)
		return nil
	}
	return j.print(metadata)
}

func (j *metadataJob) printAll(metadata map[string]string) {
	for key, value := range metadata {
		fmt.Println(fmt.Sprintf("%s=%s", key, value))
	}
}

func (j *metadataJob) print(metadata map[string]string) error {
	value, ok := metadata[j.key]
	if !ok {
		return errors.Errorf("metadata job", `"%s" does not exist in metadata of aent "%s"`, j.key, j.ctx.ID)
	}
	fmt.Println(value)
	return nil
}
