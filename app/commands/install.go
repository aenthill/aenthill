package commands

import (
	"github.com/aenthill/aenthill/context"
	"github.com/aenthill/aenthill/errors"
	"github.com/aenthill/aenthill/jobs"
	"github.com/aenthill/aenthill/manifest"

	"github.com/urfave/cli"
)

// NewInstallCommand creates a cli.Command instance.
func NewInstallCommand(context *context.Context, m *manifest.Manifest) cli.Command {
	cmd := cli.Command{
		Name:      "install",
		Aliases:   []string{"update"},
		Usage:     "Installs or updates current aent in the manifest",
		UsageText: "aenthill install [command options]",
		Action: func(ctx *cli.Context) error {
			job := jobs.NewInstallJob(ctx.StringSlice("metadata"), ctx.StringSlice("events"), context, m)
			return errors.Wrap("install command", job.Execute())
		},
	}
	cmd.Flags = []cli.Flag{
		cli.StringSliceFlag{Name: "metadata, m", Usage: "add one metadata (cumulative) - format: key=value"},
		cli.StringSliceFlag{Name: "events, e", Usage: "add one handled event (cumulative)"},
	}
	return cmd
}
