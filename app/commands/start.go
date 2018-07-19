package commands

import (
	"github.com/aenthill/aenthill/context"
	"github.com/aenthill/aenthill/errors"
	"github.com/aenthill/aenthill/jobs"
	"github.com/aenthill/aenthill/manifest"

	"github.com/urfave/cli"
)

// NewStartCommand creates a cli.Command instance.
func NewStartCommand(context *context.Context, m *manifest.Manifest) cli.Command {
	return cli.Command{
		Name:      "start",
		Aliases:   []string{"start"},
		Usage:     "Starts an aent",
		UsageText: "aenthill start image",
		Action: func(ctx *cli.Context) error {
			job, err := jobs.NewRunJob(ctx.Args().Get(0), "START", "", context, m)
			if err != nil {
				return errors.Wrap("start command", err)
			}
			return errors.Wrap("start command", job.Execute())
		},
	}
}
