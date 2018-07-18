package commands

import (
	"github.com/aenthill/aenthill/context"
	"github.com/aenthill/aenthill/errors"
	"github.com/aenthill/aenthill/jobs"
	"github.com/aenthill/aenthill/manifest"

	"github.com/urfave/cli"
)

// NewDispatchCommand creates a cli.Command instance.
func NewDispatchCommand(context *context.Context, m *manifest.Manifest) cli.Command {
	return cli.Command{
		Name:      "dispatch",
		Usage:     "Dispatches an event to all aents from manifest which could handle it",
		UsageText: "aenthill dispatch event [payload]",
		Action: func(ctx *cli.Context) error {
			job, err := jobs.NewDispatchJob(ctx.Args().Get(0), ctx.Args().Get(1), context, m)
			if err != nil {
				return errors.Wrap("dispatch command", err)
			}
			return errors.Wrap("dispatch command", job.Execute())
		},
	}
}
