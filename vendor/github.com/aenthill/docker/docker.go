// Package docker helps sending Aenthill events with the docker client binary.
package docker

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/apex/log"
	isatty "github.com/mattn/go-isatty"
)

const (
	// HostProjectDirEnvVariable is the name of the environment variable which contains the host project directory.
	// The recipient image will have this variable populated with the HostProjectDir attribute of the EventContext.
	HostProjectDirEnvVariable = "AENTHILL_HOST_PROJECT_DIR"
	// SenderEnvVariable is the name of the environment variable which contains the image name of the sender.
	// The recipient image will have this variable populated with the WhoAmI attribute of the EventContext.
	SenderEnvVariable = "AENTHILL_SENDER"
	// LogLevelEnvVariable is the name of the environment variable which contains the log level.
	// The recipient image will have this variable populated with the LogLevel attribute of the EventContext.
	LogLevelEnvVariable = "AENTHILL_LOG_LEVEL"
	// WhoAmIEnvVariable is the name of the environment variable which contains the image name.
	// The recipient and sender images should both have this environment variable populated with their respective image name.
	WhoAmIEnvVariable = "AENTHILL_WHOAMI"
	// InsideContainerProjectDir is the location in the container where the host project directory is mounted.
	InsideContainerProjectDir = "/aenthill"
	// DefaultBinary is the default binary to call in an image.
	DefaultBinary = "aent"
)

// EventContext gathers all required data of an event.
type EventContext struct {
	// WhoAmI is the image which is sending the event.
	WhoAmI string
	// Image is the image which receives the event.
	Image string
	// Binary is the command which handles the event in the targeted image.
	Binary string
	// HostProjectDir is the project directory on the host.
	HostProjectDir string
	// LogLevel is the log level which should be used by the targeted image.
	// Accepted values for log level: DEBUG, INFO, WARN, ERROR.
	LogLevel string
}

/*
Send uses the docker client binary to send an event.

It will in fact run a command in the targeted image, using the following template:

 docker run [-ti] --rm
 -v "/var/run/docker.sock:/var/run/docker.sock"
 -v "HostProjectDir:/aenthill"
 -e "AENTHILL_SENDER=WhoAmI"
 -e "AENTHILL_HOST_PROJECT_DIR=HostProjectDir"
 -e "AENTHILL_LOG_LEVEL=LogLevel"
 Image Binary even payload

Important: it relies on COMSPEC environment variable on Windows and SHELL
on posix system to know which interpreter to use for calling the docker client
binary.
*/
func Send(event string, payload string, ctx *EventContext) error {
	if err := validate(event, payload, ctx); err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"from":     ctx.WhoAmI,
		"to":       ctx.Image,
		"event":    event,
		"paypload": payload,
	}).Info("sending event")

	var e *exec.Cmd
	if runtime.GOOS == "windows" {
		e = exec.Command(os.Getenv("COMSPEC"), "/c", buildArgs(event, payload, ctx))
	} else {
		e = exec.Command(os.Getenv("SHELL"), "-c", buildArgs(event, payload, ctx))
	}

	e.Stdout = os.Stdout
	e.Stderr = os.Stderr
	e.Stdin = os.Stdin

	log.WithField("from", ctx.WhoAmI).Debugf("executing command %s", e.Args)
	fmt.Println()
	defer fmt.Println()

	return e.Run()
}

type parameterIsEmptyError struct {
	parameterName string
}

const parameterIsEmptyErrorMessage = "parameter %s is empty"

func (e *parameterIsEmptyError) Error() string {
	return fmt.Sprintf(parameterIsEmptyErrorMessage, e.parameterName)
}

func validate(event string, payload string, ctx *EventContext) error {
	if event == "" {
		return &parameterIsEmptyError{"event"}
	}

	if err := validateEventContext(ctx); err != nil {
		return err
	}

	return validateLogLevel(ctx.LogLevel)
}

type attributeValueIsEmptyError struct {
	attributeName string
	reason        string
}

const attributeValueIsEmptyErrorMessage = "attribute %s is required: %s"

func (e *attributeValueIsEmptyError) Error() string {
	return fmt.Sprintf(attributeValueIsEmptyErrorMessage, e.attributeName, e.reason)
}

func validateEventContext(ctx *EventContext) error {
	if ctx.WhoAmI == "" {
		return &attributeValueIsEmptyError{"WhoAmI", "it is the image which is sending the event"}
	}

	if ctx.Image == "" {
		return &attributeValueIsEmptyError{"Image", "it is the image which receives the event"}
	}

	if ctx.Binary == "" {
		return &attributeValueIsEmptyError{"Binary", "it is the command which handles the event in the targeted image"}
	}

	if ctx.HostProjectDir == "" {
		return &attributeValueIsEmptyError{"HostProjectDir", "it is the project directory on the host"}
	}

	if ctx.LogLevel == "" {
		return &attributeValueIsEmptyError{"LogLevel", "it is the log level which should be used by the targeted image"}
	}

	return nil
}

// levels associates log levels as used with the --logLevel -l flag from aenthill
// with its counterpart from the github.com/apex/log library.
var levels = map[string]log.Level{
	"DEBUG": log.DebugLevel,
	"INFO":  log.InfoLevel,
	"WARN":  log.WarnLevel,
	"ERROR": log.ErrorLevel,
}

type wrongLogLevelError struct{}

const wrongLogLevelErrorMessage = "accepted values for log level: DEBUG, INFO, WARN, ERROR"

func (e *wrongLogLevelError) Error() string {
	return wrongLogLevelErrorMessage
}

func validateLogLevel(logLevel string) error {
	if _, ok := levels[logLevel]; !ok {
		return &wrongLogLevelError{}
	}

	return nil
}

func buildArgs(event string, payload string, ctx *EventContext) string {
	var dockerOpts []string

	// attaches Stdin if TTY.
	if isatty.IsTerminal(os.Stdin.Fd()) {
		dockerOpts = append(dockerOpts, "-ti")
	}

	dockerOpts = append(dockerOpts, "--rm")
	dockerOpts = append(dockerOpts, fmt.Sprintf("-v \"%s:%s\"", "/var/run/docker.sock", "/var/run/docker.sock"))
	dockerOpts = append(dockerOpts, fmt.Sprintf("-v \"%s:%s\"", ctx.HostProjectDir, InsideContainerProjectDir))
	dockerOpts = append(dockerOpts, fmt.Sprintf("-e \"%s=%s\"", SenderEnvVariable, ctx.WhoAmI))
	dockerOpts = append(dockerOpts, fmt.Sprintf("-e \"%s=%s\"", HostProjectDirEnvVariable, ctx.HostProjectDir))
	dockerOpts = append(dockerOpts, fmt.Sprintf("-e \"%s=%s\"", LogLevelEnvVariable, ctx.LogLevel))

	var args []string
	args = append(args, []string{"docker", "run"}...)
	args = append(args, dockerOpts...)
	args = append(args, []string{ctx.Image, ctx.Binary, event, payload}...)

	return strings.Join(args, " ")
}
