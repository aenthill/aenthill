// Package context is a solution for gathering all required data of the application.
package context

import (
	"fmt"
	"os"

	"github.com/aenthill/aenthill/errors"
	"github.com/aenthill/aenthill/manifest"
)

const (
	// KeyEnvVar is the name of the environment variable which contains the current aent key from the manifest if the aent is a registred instance.
	KeyEnvVar = "PHEROMONE_KEY"
	// ImageEnvVar is the name of the environment variable which contains the current aent image name.
	ImageEnvVar = "PHEROMONE_IMAGE_NAME"
	// FromContainerIDEnvVar is the name of the environment variable which contains the sender container ID which has started the recipient image.
	FromContainerIDEnvVar = "PHEROMONE_FROM_CONTAINER_ID"
	// HostnameEnvVar is the name of the environment variable which contains the recipient container id. It is populated by Docker.
	HostnameEnvVar = "HOSTNAME"
	// HostProjectDirEnvVar is the name of the environment variable which contains the host project directory path.
	HostProjectDirEnvVar = "PHEROMONE_HOST_PROJECT_DIR"
	// ContainerProjectDirEnvVar is the name of the environment variable which contains the mounted path of the host project directory.
	ContainerProjectDirEnvVar = "PHEROMONE_CONTAINER_PROJECT_DIR"
	// LogLevelEnvVar is the name of the environment variable which contains the log level.
	LogLevelEnvVar = "PHEROMONE_LOG_LEVEL"
)

// Context is our working struct.
type Context struct {
	isContainer     bool
	Image           string
	Key             string
	FromContainerID string
	Hostname        string
	HostProjectDir  string
	ProjectDir      string
	LogLevel        string
}

// New creates a Context instance according to where is launched the application
// (form a container or from the host).
func New() (*Context, error) {
	v, err := lookupEnv(ImageEnvVar)
	if err != nil {
		return makeFromHost()
	}
	return makeFromEnv(v)
}

// IsContainer returns true if the application is launched from a container, false otherwise.
func (ctx *Context) IsContainer() bool {
	return ctx.isContainer
}

// PopulateEnv tries to populate the PHEROMONE_KEY, PHEROMONE_METADATA_* and PHEROMONE_DEPENDENCY_*
// environment variables.
func (ctx *Context) PopulateEnv(m *manifest.Manifest) error {
	if ctx.Key == "" {
		return nil
	}
	if err := os.Setenv(KeyEnvVar, ctx.Key); err != nil {
		return errors.Wrap("context", err)
	}
	aent, err := m.Aent(ctx.Key)
	if err != nil {
		return errors.Wrap("context", err)
	}
	if aent.Metadata != nil {
		for key, value := range aent.Metadata {
			if err := os.Setenv(fmt.Sprintf("PHEROMONE_METADATA_%s", key), value); err != nil {
				return errors.Wrap("context", err)
			}
		}
	}
	if aent.Dependencies != nil {
		for key, value := range aent.Dependencies {
			if err := os.Setenv(fmt.Sprintf("PHEROMONE_DEPENDENCY_%s", key), value); err != nil {
				return errors.Wrap("context", err)
			}
		}
	}
	return nil
}

func makeFromHost() (*Context, error) {
	ctx := &Context{}
	ctx.isContainer = false
	ctx.LogLevel = "ERROR"
	projectDir, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrap("context", err)
	}
	ctx.HostProjectDir, ctx.ProjectDir = projectDir, "/aenthill"
	return ctx, nil
}

func makeFromEnv(image string) (*Context, error) {
	var (
		ctx = &Context{isContainer: true, Image: image}
		v   string
		err error
	)
	v, err = lookupEnv(KeyEnvVar)
	if err != nil {
		return nil, err
	}
	ctx.Key = v
	v, err = lookupEnv(FromContainerIDEnvVar)
	if err != nil {
		return nil, err
	}
	ctx.FromContainerID = v
	v, err = lookupEnv(HostnameEnvVar)
	if err != nil {
		return nil, err
	}
	ctx.Hostname = v
	v, err = lookupEnv(HostProjectDirEnvVar)
	if err != nil {
		return nil, err
	}
	ctx.HostProjectDir = v
	v, err = lookupEnv(ContainerProjectDirEnvVar)
	if err != nil {
		return nil, err
	}
	ctx.ProjectDir = v
	v, err = lookupEnv(LogLevelEnvVar)
	if err != nil {
		return nil, err
	}
	ctx.LogLevel = v
	return ctx, nil
}

func lookupEnv(key string) (string, error) {
	v, ok := os.LookupEnv(key)
	if !ok {
		return "", errors.Errorf("context", `env key "%s" does not exist`, key)
	}
	if key != FromContainerIDEnvVar && key != KeyEnvVar && v == "" {
		return "", errors.Errorf("context", `env key "%s" has an empty value`, key)
	}
	return v, nil
}
