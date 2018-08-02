package commands

import (
	"github.com/aenthill/aenthill/internal/pkg/context"
	"github.com/aenthill/aenthill/internal/pkg/errors"
	"github.com/aenthill/aenthill/internal/pkg/jobs"
	"github.com/aenthill/aenthill/internal/pkg/manifest"

	"github.com/urfave/cli"
)

// NewUpdateCommand creates a cli.Command instance.
func NewUpdateCommand(context *context.Context, m *manifest.Manifest) cli.Command {
	cmd := cli.Command{
		Name:      "update",
		Usage:     "Updates current aent in the manifest",
		UsageText: "aenthill [global options] update [command options]",
		Action: func(ctx *cli.Context) error {
			job, err := jobs.NewUpdateJob(ctx.StringSlice("metadata"), ctx.StringSlice("events"), context, m)
			if err != nil {
				return errors.Wrap("update command", err)
			}
			return errors.Wrap("update command", job.Execute())
		},
	}
	cmd.Flags = []cli.Flag{
		cli.StringSliceFlag{Name: "metadata, m", Usage: "add one metadata (cumulative) - format: key=value"},
		cli.StringSliceFlag{Name: "events, e", Usage: "add one handled event (cumulative)"},
	}
	return cmd
}
