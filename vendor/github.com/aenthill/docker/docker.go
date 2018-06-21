// Package docker is a simple wrapper around the Docker client binary for sending Aenthill events.
package docker

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/aenthill/log"
	isatty "github.com/mattn/go-isatty"
)

const (
	// FromImageNameEnvVariable is the name of the environment variable which contains the sender image name which has started the recipient image.
	FromImageNameEnvVariable = "PHEROMONE_FROM_IMAGE_NAME"
	// FromContainerIDEnvVariable is the name of the environment variable which contains the sender container ID which has started the recipient image.
	FromContainerIDEnvVariable = "PHEROMONE_FROM_CONTAINER_ID"
	// WhoAmIEnvVariable is the name of the environment variable which contains the recipient image name.
	WhoAmIEnvVariable = "PHEROMONE_WHOAMI"
	// HostnameEnvVariable is the name of the environment variable which contains the recipient container id. It is populated by Docker.
	HostnameEnvVariable = "HOSTNAME"
	// HostProjectDirEnvVariable is the name of the environment variable which contains the host project directory path.
	HostProjectDirEnvVariable = "PHEROMONE_HOST_PROJECT_DIR"
	// ContainerProjectDirEnvVariable is the name of the environment variable which contains the mounted path of the host project directory.
	ContainerProjectDirEnvVariable = "PHEROMONE_CONTAINER_PROJECT_DIR"
	// LogLevelEnvVariable is the name of the environment variable which contains the log level.
	LogLevelEnvVariable = "PHEROMONE_LOG_LEVEL"
)

// EventContext gathers all required data of an event.
type EventContext struct {
	// fromImageName is the sender image name.
	fromImageName string
	// fromContainerID is the sender container ID.
	fromContainerID string
	// toImageName is the image which receives the event (when calling Run function).
	toImageName string
	// toContainerID is the container ID which receives the event (when calling Exec function).
	toContainerID string
	// hostProjectDir is the project directory on the host.
	// Only required when using the Run function.
	hostProjectDir string
	// logLevel is the log level which should be used by the targeted image.
	// Use one of the log level provided by the https://github.com/aenthill/log
	// library.
	// Only required when using the Run function.
	logLevel string
}

// NewRunEventContext creates an EventContext instance with correct attributes for
// the Run function.
func NewRunEventContext(fromImageName string, fromContainerID string, toImageName string, hostProjectDir string, logLevel string) *EventContext {
	return &EventContext{
		fromImageName:   fromImageName,
		fromContainerID: fromContainerID,
		toImageName:     toImageName,
		hostProjectDir:  hostProjectDir,
		logLevel:        logLevel,
	}
}

/*
Run uses the Docker client binary to start an image by sending an event and a payload to it.

It will in fact run a command in the targeted image, using the following template:

 docker run [-ti] --rm
 -v "/var/run/docker.sock:/var/run/docker.sock"
 -v "HostProjectDir:/aenthill"
 -e "PHEROMONE_FROM_CONTAINER_ID=EventContext.fromContainerID"
 -e "PHEROMONE_WHOAMI=EventContext.whoAmI"
 -e "PHEROMONE_HOST_PROJECT_DIR=EventContext.hostProjectDir"
 -e "PHEROMONE_CONTAINER_PROJECT_DIR=/aenthill"
 -e "PHEROMONE_LOG_LEVEL=EventContext.logLevel"
 EventContext.toImageName aent "event" "payload"
*/
func Run(event string, payload string, ctx *EventContext) error {
	cmd := makeDockerRunCommand(event, payload, ctx)
	e, err := newExecCmd(cmd)
	if err != nil {
		return err
	}

	entryCtx := &log.EntryContext{
		Event:         event,
		Payload:       payload,
		FromImageName: ctx.fromImageName,
		ToImageName:   ctx.toImageName,
	}
	log.Info(entryCtx, "awakening with event")
	log.Debugf(entryCtx, "executing command %s", e.Args)

	return e.Run()
}

// NewExecEventContext creates an EventContext instance with correct attributes for
// the Exec function.
func NewExecEventContext(fromImageName string, toImageName string, toContainerID string) *EventContext {
	return &EventContext{
		fromImageName: fromImageName,
		toImageName:   toImageName,
		toContainerID: toContainerID,
	}
}

/*
Exec uses the Docker client binary to reply to a container by sending an event and a payload to it.

It will in fact run a command in the targeted container, using the following template:

 docker exec [-ti]
 EventContext.to aent "event" "payload"
*/
func Exec(event string, payload string, ctx *EventContext) error {
	cmd := makeDockerExecCommand(event, payload, ctx)
	e, err := newExecCmd(cmd)
	if err != nil {
		return err
	}

	entryCtx := &log.EntryContext{
		Event:         event,
		Payload:       payload,
		FromImageName: ctx.fromImageName,
		ToImageName:   ctx.toImageName,
	}
	log.Info(entryCtx, "replying with event")
	log.Debugf(entryCtx, "executing command %s", e.Args)

	return e.Run()
}

type dockerBinaryNotFoundError struct{}

const dockerBinaryNotFoundErrorMessage = "Docker client binary was not found"

func (e *dockerBinaryNotFoundError) Error() string {
	return dockerBinaryNotFoundErrorMessage
}

type interpreterNotFoundError struct {
	envVar string
}

const interpreterNotFoundErrorMessage = "%s is a required environment variable: it allows to know which interpreter to use for executing the Docker command"

func (e *interpreterNotFoundError) Error() string {
	return fmt.Sprintf(interpreterNotFoundErrorMessage, e.envVar)
}

func newExecCmd(command string) (*exec.Cmd, error) {
	if _, err := exec.LookPath("docker"); err != nil {
		return nil, &dockerBinaryNotFoundError{}
	}

	var (
		envVar string
		flag   string
	)

	if runtime.GOOS == "windows" {
		envVar = "COMSPEC"
		flag = "/c"
	} else {
		envVar = "SHELL"
		flag = "-c"
	}

	interpreter := os.Getenv(envVar)
	if interpreter == "" {
		return nil, &interpreterNotFoundError{envVar}
	}

	e := exec.Command(interpreter, flag, command)
	e.Stdout = os.Stdout
	e.Stderr = os.Stderr
	e.Stdin = os.Stdin

	return e, nil
}

func makeDockerRunCommand(event string, payload string, ctx *EventContext) string {
	var flags []string
	flags = append(flags, makeTTYFlag())
	flags = append(flags, "--rm")
	flags = append(flags, makeVolumeFlags(ctx)...)
	flags = append(flags, makeEnvFlags(ctx)...)

	var command []string
	command = append(command, []string{"docker", "run"}...)
	command = append(command, flags...)
	command = append(command, []string{ctx.toImageName, imageOrContainerBinary, sanitize(event), sanitize(payload)}...)

	return strings.Join(command, " ")
}

func makeDockerExecCommand(event string, payload string, ctx *EventContext) string {
	var flags []string
	flags = append(flags, makeTTYFlag())

	var command []string
	command = append(command, []string{"docker", "exec"}...)
	command = append(command, flags...)
	command = append(command, []string{ctx.toContainerID, imageOrContainerBinary, sanitize(event), sanitize(payload)}...)

	return strings.Join(command, " ")
}

const (
	dockerSocket           = "/var/run/docker.sock"
	containerProjectDir    = "/aenthill"
	imageOrContainerBinary = "aent"
)

func makeTTYFlag() string {
	var flag string

	// attaches Stdin if TTY.
	if isatty.IsTerminal(os.Stdin.Fd()) {
		flag = "-ti"
	}

	return flag
}

func makeEnvFlags(ctx *EventContext) []string {
	var flags []string
	flags = append(flags, fmt.Sprintf(`-e "%s=%s"`, FromImageNameEnvVariable, ctx.fromImageName))
	flags = append(flags, fmt.Sprintf(`-e "%s=%s"`, FromContainerIDEnvVariable, ctx.fromContainerID))
	flags = append(flags, fmt.Sprintf(`-e "%s=%s"`, WhoAmIEnvVariable, ctx.toImageName))
	flags = append(flags, fmt.Sprintf(`-e "%s=%s"`, HostProjectDirEnvVariable, ctx.hostProjectDir))
	flags = append(flags, fmt.Sprintf(`-e "%s=%s"`, ContainerProjectDirEnvVariable, containerProjectDir))
	flags = append(flags, fmt.Sprintf(`-e "%s=%s"`, LogLevelEnvVariable, ctx.logLevel))

	return flags
}

func makeVolumeFlags(ctx *EventContext) []string {
	var flags []string
	flags = append(flags, fmt.Sprintf(`-v "%s:%s"`, dockerSocket, dockerSocket))
	flags = append(flags, fmt.Sprintf(`-v "%s:%s"`, ctx.hostProjectDir, containerProjectDir))

	return flags
}

func sanitize(str string) string {
	return fmt.Sprintf(`"%s"`, strings.Replace(str, `"`, `\"`, -1))
}
