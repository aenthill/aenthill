package docker

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	isatty "github.com/mattn/go-isatty"
)

type builder struct {
	command string
	target  string
	event   string
	payload string
	env     map[string]string
	mount   map[string]string
}

func (b *builder) run(image, event, payload string) {
	b.command = "run"
	b.target = image
	b.event = event
	b.payload = payload
}

func (b *builder) exec(containerID, event, payload string) {
	b.command = "exec"
	b.target = containerID
	b.event = event
	b.payload = payload
}

func (b *builder) withEnv(key, value string) {
	if b.env == nil {
		b.env = make(map[string]string)
	}
	b.env[key] = value
}

func (b *builder) withMount(source, target string) {
	if b.mount == nil {
		b.mount = make(map[string]string)
	}
	b.mount[source] = target
}

func (b *builder) build() *exec.Cmd {
	var args []string
	args = append(args, b.command)
	// attaches Stdin if TTY.
	if isatty.IsTerminal(os.Stdin.Fd()) {
		args = append(args, "-ti")
	}
	if b.command == "run" {
		args = append(args, "--rm")
	}
	for key, value := range b.env {
		args = append(args, "-e", fmt.Sprintf("%s=%s", key, value))
	}
	for source, target := range b.mount {
		args = append(args, "--mount", fmt.Sprintf("type=bind,source=%s,target=%s", source, target))
	}
	args = append(args, b.target, "aent", b.event, sanitize(b.payload))
	cmd := exec.Command("docker", args...)
	cmd.Stdout, cmd.Stderr, cmd.Stdin = os.Stdout, os.Stderr, os.Stdin
	return cmd
}

func sanitize(payload string) string {
	return fmt.Sprintf(`"%s"`, strings.Replace(payload, `"`, `\"`, -1))
}
