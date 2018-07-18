package commands

import (
	"github.com/aenthill/aenthill/context"
	"github.com/aenthill/aenthill/errors"
	"github.com/aenthill/aenthill/jobs"
	"github.com/aenthill/aenthill/manifest"

	"github.com/urfave/cli"
)

// NewAddCommand creates a cli.Command instance.
func NewAddCommand(context *context.Context, m *manifest.Manifest) cli.Command {
	return cli.Command{
		Name:      "add",
		Aliases:   []string{"a"},
		Usage:     "Starts an aent",
		UsageText: "aenthill add image",
		Action: func(ctx *cli.Context) error {
			job, err := jobs.NewRunJob(ctx.Args().Get(0), "ADD", "", context, m)
			if err != nil {
				return errors.Wrap("add command", err)
			}
			return errors.Wrap("add command", job.Execute())
		},
	}
}
