package commands

import (
	"github.com/aenthill/aenthill/internal/pkg/context"
	"github.com/aenthill/aenthill/internal/pkg/errors"
	"github.com/aenthill/aenthill/internal/pkg/jobs"
	"github.com/aenthill/aenthill/internal/pkg/manifest"

	"github.com/urfave/cli"
)

// NewStartCommand creates a cli.Command instance.
func NewStartCommand(context *context.Context, m *manifest.Manifest) cli.Command {
	return cli.Command{
		Name:      "start",
		Aliases:   []string{"s"},
		Usage:     "Starts an aent",
		UsageText: "aenthill [global options] start image",
		Action: func(ctx *cli.Context) error {
			job, err := jobs.NewRunJob(ctx.Args().Get(0), "START", "", context, m)
			if err != nil {
				return errors.Wrap("start command", err)
			}
			return errors.Wrap("start command", job.Execute())
		},
	}
}
