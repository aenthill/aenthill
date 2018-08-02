package jobs

import (
	"fmt"
	"os"

	"github.com/aenthill/aenthill/internal/pkg/context"
	"github.com/aenthill/aenthill/internal/pkg/errors"
	"github.com/aenthill/aenthill/internal/pkg/manifest"

	isatty "github.com/mattn/go-isatty"
)

type dependencyJob struct {
	key      string
	ctx      *context.Context
	manifest *manifest.Manifest
}

// NewDependencyJob creates a new Job instance.
func NewDependencyJob(key string, ctx *context.Context, m *manifest.Manifest) (Job, error) {
	if err := m.Parse(); err != nil {
		return nil, errors.Wrap("dependency job", err)
	}
	if ctx.ID == "" {
		return nil, errors.Error("dependency job", "aent ID is missing")
	}
	return &dependencyJob{key, ctx, m}, nil
}

func (j *dependencyJob) Execute() error {
	dependencies, err := j.manifest.Dependencies(j.ctx.ID)
	if err != nil {
		return errors.Wrap("dependency job", err)
	}
	value, ok := dependencies[j.key]
	if !ok {
		return errors.Errorf("dependency job", `"%s" does not exist in dependencies of aent "%s"`, j.key, j.ctx.ID)
	}
	if isatty.IsTerminal(os.Stdin.Fd()) {
		fmt.Println(value)
	} else {
		fmt.Print(value)
	}
	return nil
}
