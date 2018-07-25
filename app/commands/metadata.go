package commands

import (
	"github.com/aenthill/aenthill/context"
	"github.com/aenthill/aenthill/errors"
	"github.com/aenthill/aenthill/jobs"
	"github.com/aenthill/aenthill/manifest"

	"github.com/urfave/cli"
)

// NewMetadataCommand creates a cli.Command instance.
func NewMetadataCommand(context *context.Context, m *manifest.Manifest) cli.Command {
	cmd := cli.Command{
		Name:      "metadata",
		Usage:     "Prints an entry of metadata of current aent",
		UsageText: "aenthill [global options] metadata key",
		Action: func(ctx *cli.Context) error {
			job, err := jobs.NewMetadataJob(ctx.Args().Get(0), context, m)
			if err != nil {
				errors.Wrap("metadata command", err)
			}
			return errors.Wrap("metadata command", job.Execute())
		},
	}
	return cmd
}
