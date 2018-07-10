package commands

import (
	"github.com/aenthill/aenthill/context"
	"github.com/aenthill/aenthill/errors"
	"github.com/aenthill/aenthill/jobs"
	"github.com/aenthill/aenthill/manifest"

	"github.com/urfave/cli"
)

func NewInstallCommand(context *context.Context, m *manifest.Manifest) cli.Command {
	cmd := cli.Command{
		Name:      "install",
		Aliases:   []string{"update"},
		Usage:     "Installs or updates current aent in the manifest",
		UsageText: "aenthill install [command options]",
		Action: func(ctx *cli.Context) error {
			job := jobs.NewInstallJob(ctx.StringSlice("events"), ctx.StringSlice("metadata"), context, m)
			return errors.Wrap("install command", job.Execute())
		},
	}
	cmd.Flags = []cli.Flag{
		cli.StringSliceFlag{Name: "events, e", Usage: "add one handled event (cumulative)"},
		cli.StringSliceFlag{Name: "metadata, m", Usage: "add one metadata (cumulative) - format: key=value"},
	}
	return cmd
}
