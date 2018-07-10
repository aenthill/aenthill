package context

import (
	"os"

	"github.com/aenthill/aenthill/errors"
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

type Context struct {
	isContainer     bool
	Key             string
	Image           string
	FromContainerID string
	Hostname        string
	AentKey         string
	HostProjectDir  string
	ProjectDir      string
	LogLevel        string
}

func New() (*Context, error) {
	if _, err := lookupEnv(ImageEnvVar); err != nil {
		return makeFromHost()
	}
	return makeFromEnv()
}

func (ctx *Context) IsContainer() bool {
	return ctx.isContainer
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

func makeFromEnv() (*Context, error) {
	var (
		ctx = &Context{}
		v   string
		err error
	)
	ctx.isContainer = true
	v, err = lookupEnv(KeyEnvVar)
	if err != nil {
		return nil, err
	}
	v, err = lookupEnv(ImageEnvVar)
	if err != nil {
		return nil, err
	}
	v, err = lookupEnv(FromContainerIDEnvVar)
	if err != nil {
		return nil, err
	}
	ctx.FromContainerID = v
	v, err = lookupEnv(HostnameEnvVar)
	if err != nil {
		return nil, err
	}
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
