package commands

import (
	"github.com/aenthill/aenthill/context"
	"github.com/aenthill/aenthill/errors"
	"github.com/aenthill/aenthill/jobs"
	"github.com/aenthill/aenthill/manifest"

	"github.com/urfave/cli"
)

// NewRunCommand creates a cli.Command instance.
func NewRunCommand(context *context.Context, m *manifest.Manifest) cli.Command {
	return cli.Command{
		Name:      "run",
		Usage:     "Starts an aent",
		UsageText: "aenthill [global options] run image|ID event [payload]",
		Action: func(ctx *cli.Context) error {
			job, err := jobs.NewRunJob(ctx.Args().Get(0), ctx.Args().Get(1), ctx.Args().Get(2), context, m)
			if err != nil {
				return errors.Wrap("run command", err)
			}
			return errors.Wrap("run command", job.Execute())
		},
	}
}
