package jobs

import (
	"fmt"
	"os"

	"github.com/aenthill/aenthill/context"
	"github.com/aenthill/aenthill/errors"
	"github.com/aenthill/aenthill/manifest"

	isatty "github.com/mattn/go-isatty"
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
	value, ok := metadata[j.key]
	if !ok {
		return errors.Errorf("metadata job", `"%s" does not exist in metadata of aent "%s"`, j.key, j.ctx.ID)
	}
	if isatty.IsTerminal(os.Stdin.Fd()) {
		fmt.Println(value)
	} else {
		fmt.Print(value)
	}
	return nil
}
