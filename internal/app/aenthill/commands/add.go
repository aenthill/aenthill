package commands

import (
	"github.com/aenthill/aenthill/internal/pkg/context"
	"github.com/aenthill/aenthill/internal/pkg/errors"
	"github.com/aenthill/aenthill/internal/pkg/jobs"
	"github.com/aenthill/aenthill/internal/pkg/manifest"

	"github.com/urfave/cli"
)

// NewAddCommand creates a cli.Command instance.
func NewAddCommand(context *context.Context, m *manifest.Manifest) cli.Command {
	return cli.Command{
		Name:      "add",
		Aliases:   []string{"a"},
		Usage:     "Adds a service in your Docker project",
		UsageText: "aenthill [global options] add image",
		Action: func(ctx *cli.Context) error {
			if err := validateArgsLength(ctx, 1, 1); err != nil {
				return errors.Wrap("add command", err)
			}
			job, err := jobs.NewRunJob(ctx.Args().Get(0), "ADD", "", context, m)
			if err != nil {
				return errors.Wrap("add command", err)
			}
			return errors.Wrap("add command", job.Execute())
		},
	}
}
