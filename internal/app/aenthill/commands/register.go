package commands

import (
	"github.com/aenthill/aenthill/internal/pkg/context"
	"github.com/aenthill/aenthill/internal/pkg/errors"
	"github.com/aenthill/aenthill/internal/pkg/jobs"
	"github.com/aenthill/aenthill/internal/pkg/manifest"

	"github.com/urfave/cli"
)

// NewRegisterCommand creates a cli.Command instance.
func NewRegisterCommand(context *context.Context, m *manifest.Manifest) cli.Command {
	cmd := cli.Command{
		Name:      "register",
		Usage:     "Adds a dependency to current aent in the manifest",
		UsageText: "aenthill [global options] register image key [command options]",
		Action: func(ctx *cli.Context) error {
			if err := validateArgsLength(ctx, 2, 2); err != nil {
				return errors.Wrap("register command", err)
			}
			job, err := jobs.NewRegisterJob(ctx.Args().Get(0), ctx.Args().Get(1), ctx.StringSlice("metadata"), context, m)
			if err != nil {
				return errors.Wrap("register command", err)
			}
			return errors.Wrap("register command", job.Execute())
		},
	}
	cmd.Flags = []cli.Flag{
		cli.StringSliceFlag{Name: "metadata, m", Usage: "add one metadata (cumulative) - format: key=value"},
	}
	return cmd
}
