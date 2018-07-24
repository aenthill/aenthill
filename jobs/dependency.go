package jobs

import (
	"fmt"

	"github.com/aenthill/aenthill/context"
	"github.com/aenthill/aenthill/errors"
	"github.com/aenthill/aenthill/manifest"
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
	if j.key == "" {
		j.printAll(dependencies)
		return nil
	}
	return j.print(dependencies)
}

func (j *dependencyJob) printAll(dependencies map[string]string) {
	for key, ID := range dependencies {
		fmt.Println(fmt.Sprintf("%s=%s", key, ID))
	}
}

func (j *dependencyJob) print(dependencies map[string]string) error {
	ID, ok := dependencies[j.key]
	if !ok {
		return errors.Errorf("dependency job", `"%s" does not exist in dependencies of aent "%s"`, j.key, j.ctx.ID)
	}
	fmt.Println(ID)
	return nil
}
