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
	// FromEnvVariable is the name of the environment variable which contains the sender image name.
	// The recipient image will have this variable populated with the WhoAmI attribute of the EventContext.
	FromEnvVariable = "PHEROMONE_FROM"
	// WhoAmIEnvVariable is the name of the environment variable which contains the recipient image name.
	// The recipient image will have this variable populated with the To attribute of the EventContext.
	WhoAmIEnvVariable = "PHEROMONE_WHOAMI"
	// HostProjectDirEnvVariable is the name of the environment variable which contains the host project directory path.
	// The recipient image will have this variable populated with the HostProjectDir attribute of the EventContext.
	HostProjectDirEnvVariable = "PHEROMONE_HOST_PROJECT_DIR"
	// ContainerProjectDirEnvVariable is the name of the environment variable which contains the mounted path of the host project directory.
	// The recipient image will have this variable populated with "/aenthill".
	ContainerProjectDirEnvVariable = "PHEROMONE_CONTAINER_PROJECT_DIR"
	// LogLevelEnvVariable is the name of the environment variable which contains the log level.
	// The recipient image will have this variable populated with the LogLevel attribute of the EventContext.
	LogLevelEnvVariable = "PHEROMONE_LOG_LEVEL"
)

// EventContext gathers all required data of an event.
type EventContext struct {
	// Source is the binary which is sending the event.
	Source string
	// From is the image which is sending the event.
	From string
	// To is the image which receives the event.
	To string
	// HostProjectDir is the project directory on the host.
	HostProjectDir string
	// LogLevel is the log level which should be used by the targeted image.
	// Use one of the log level provided by the https://github.com/aenthill/log
	// library.
	LogLevel string
}

type dockerBinaryNotFoundError struct{}

const dockerBinaryNotFoundErrorMessage = "Docker client binary was not found"

func (e *dockerBinaryNotFoundError) Error() string {
	return dockerBinaryNotFoundErrorMessage
}

/*
Execute uses the Docker client binary to send an Aenthill event.

It will in fact run a command in the targeted image, using the following template:

 docker run [-ti] --rm
 -v "/var/run/docker.sock:/var/run/docker.sock"
 -v "HostProjectDir:/aenthill"
 -e "PHEROMONE_FROM=EventContext.From"
 -e "PHEROMONE_WHOAMI=EventContext.To"
 -e "PHEROMONE_HOST_PROJECT_DIR=EventContext.HostProjectDir"
 -e "PHEROMONE_CONTAINER_PROJECT_DIR=/aenthill"
 -e "PHEROMONE_LOG_LEVEL=EventContext.LogLevel"
 EventContext.To aent event payload
*/
func Execute(event string, payload string, ctx *EventContext) error {
	if _, err := exec.LookPath("docker"); err != nil {
		return &dockerBinaryNotFoundError{}
	}

	cmd := buildDockerCommand(event, payload, ctx)
	e, err := newExecCmd(cmd)
	if err != nil {
		return err
	}

	entryCtx := &log.EntryContext{
		Source:    ctx.Source,
		Event:     event,
		Payload:   payload,
		Image:     ctx.From,
		Recipient: ctx.To,
	}
	log.Info(entryCtx, "new event")
	log.Debugf(entryCtx, "executing command %s", e.Args)

	return e.Run()
}

type interpreterNotFoundError struct {
	envVar string
}

const interpreterNotFoundErrorMessage = "%s is a required environment variable: it allows to know which interpreter to use for executing the Docker command"

func (e *interpreterNotFoundError) Error() string {
	return fmt.Sprintf(interpreterNotFoundErrorMessage, e.envVar)
}

func newExecCmd(command string) (*exec.Cmd, error) {
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

func buildDockerCommand(event string, payload string, ctx *EventContext) string {
	var flags []string

	// attaches Stdin if TTY.
	if isatty.IsTerminal(os.Stdin.Fd()) {
		flags = append(flags, "-ti")
	}

	const (
		dockerSocket        = "/var/run/docker.sock"
		containerProjectDir = "/aenthill"
	)

	flags = append(flags, "--rm")
	flags = append(flags, fmt.Sprintf("-v \"%s:%s\"", dockerSocket, dockerSocket))
	flags = append(flags, fmt.Sprintf("-v \"%s:%s\"", ctx.HostProjectDir, containerProjectDir))
	flags = append(flags, fmt.Sprintf("-e \"%s=%s\"", FromEnvVariable, ctx.From))
	flags = append(flags, fmt.Sprintf("-e \"%s=%s\"", WhoAmIEnvVariable, ctx.To))
	flags = append(flags, fmt.Sprintf("-e \"%s=%s\"", HostProjectDirEnvVariable, ctx.HostProjectDir))
	flags = append(flags, fmt.Sprintf("-e \"%s=%s\"", ContainerProjectDirEnvVariable, containerProjectDir))
	flags = append(flags, fmt.Sprintf("-e \"%s=%s\"", LogLevelEnvVariable, ctx.LogLevel))

	var command []string
	command = append(command, []string{"docker", "run"}...)
	command = append(command, flags...)
	command = append(command, []string{ctx.To, "aent", event, payload}...)

	return strings.Join(command, " ")
}
