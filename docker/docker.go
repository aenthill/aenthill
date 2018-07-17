// Package docker is a simple wrapper around the Docker client binary for sending Aenthill events.
package docker

import (
	"fmt"
	"os/exec"

	"github.com/aenthill/aenthill/context"
	"github.com/aenthill/aenthill/errors"
)

type Docker struct {
	ctx *context.Context
}

func New(ctx *context.Context) (*Docker, error) {
	if _, err := exec.LookPath("docker"); err != nil {
		return nil, errors.Wrap("docker", err)
	}
	return &Docker{ctx}, nil
}

func (d *Docker) Run(key, image, event, payload string, metadata map[string]string) error {
	b := &builder{}
	b.run(image, event, payload)
	b.withEnv(context.KeyEnvVar, key)
	b.withEnv(context.ImageEnvVar, d.ctx.Image)
	b.withEnv(context.FromContainerIDEnvVar, d.ctx.FromContainerID)
	b.withEnv(context.HostProjectDirEnvVar, d.ctx.HostProjectDir)
	b.withEnv(context.ContainerProjectDirEnvVar, d.ctx.ProjectDir)
	b.withEnv(context.LogLevelEnvVar, d.ctx.LogLevel)
	if metadata != nil {
		for key, value := range metadata {
			b.withEnv(fmt.Sprintf("PHEROMONE_METADATA_%s", key), value)
		}
	}
	b.withMount("/var/run/docker.sock", "/var/run/docker.sock")
	b.withMount(d.ctx.HostProjectDir, d.ctx.ProjectDir)
	cmd := b.build()
	return errors.Wrapf("docker", cmd.Run(), "%s", cmd.Args)
}

func (d *Docker) Reply(event, payload string) error {
	b := &builder{}
	b.exec(d.ctx.FromContainerID, event, payload)
	cmd := b.build()
	return errors.Wrapf("docker", cmd.Run(), "%s", cmd.Args)
}
