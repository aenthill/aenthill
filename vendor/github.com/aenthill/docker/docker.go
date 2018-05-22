// Package docker helps sending Aenthill events with the docker client binary.
package docker

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	isatty "github.com/mattn/go-isatty"
)

const (
	// HostProjectDirEnvVariable will be populated with the host project directory by Aenthill.
	// You may use this constant with os.Getenv to retrieve the current project directory.
	HostProjectDirEnvVariable = "AENTHILL_HOST_PROJECT_DIR"
	// InsideContainerProjectDir is the location the the container where the host project directory is mounted.
	InsideContainerProjectDir = "/aenthill"
	// DefaultBinary is the default binary to call in an image.
	DefaultBinary                  = "aent"
	dockerSocket                   = "/var/run/docker.sock"
	defaultWindowsShellEnvVariable = "COMSPEC"
	defaultPosixShellEnvVariable   = "SHELL"
)

// EventContext gathers all required data of an event.
type EventContext struct {
	// Image is the image which receives the event.
	Image string
	// Binary is the command which handles the event in the targeted image.
	Binary string
	// HostProjectDir is the project directory on the host.
	HostProjectDir string
	// Payload are custom event data.
	Payload string
}

/*
Send use the docker client binary to send an event.

It will in fact run a command in the targeted image, using the following template:

 docker run [opts] IMAGE BINARY EVENT PAYLOAD SENDER_IMAGE

Important: it relies on COMSPEC environment variable on Windows and SHELL
on posix system to know which interpreter to use for calling the docker client
binary.
*/
func Send(event string, context *EventContext) error {
	var dockerOpts []string

	// attaches Stdin if TTY is a terminal.
	if isatty.IsTerminal(os.Stdin.Fd()) {
		dockerOpts = append(dockerOpts, "-ti")
	}

	dockerOpts = append(dockerOpts, "--rm")
	dockerOpts = append(dockerOpts, fmt.Sprintf("-v \"%s:%s\"", dockerSocket, dockerSocket))
	dockerOpts = append(dockerOpts, fmt.Sprintf("-v \"%s:%s\"", context.HostProjectDir, InsideContainerProjectDir))
	dockerOpts = append(dockerOpts, fmt.Sprintf("-e \"%s=%s\"", HostProjectDirEnvVariable, context.HostProjectDir))

	var args []string
	args = append(args, []string{"docker", "run"}...)
	args = append(args, dockerOpts...)
	args = append(args, []string{context.Image, context.Binary, event, context.Payload}...)

	var e *exec.Cmd
	if runtime.GOOS == "windows" {
		e = exec.Command(os.Getenv(defaultWindowsShellEnvVariable), "/c", strings.Join(args, " "))
	} else {
		e = exec.Command(os.Getenv(defaultPosixShellEnvVariable), "-c", strings.Join(args, " "))
	}

	e.Stdout = os.Stdout
	e.Stderr = os.Stderr
	e.Stdin = os.Stdin

	return e.Run()
}
