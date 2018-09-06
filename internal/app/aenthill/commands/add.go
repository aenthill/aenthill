package commands

import (
	"github.com/aenthill/aenthill/internal/pkg/context"
	"github.com/aenthill/aenthill/internal/pkg/errors"
	"github.com/aenthill/aenthill/internal/pkg/jobs"
	"github.com/aenthill/aenthill/internal/pkg/manifest"

	"github.com/urfave/cli"
)

// NewAddCommand creates a cli.Command instance.
func NewAddCommand(context *context.Context, m *manifest.Manifest) cli.Command {
	return cli.Command{
		Name:      "add",
		Aliases:   []string{"a"},
		Usage:     "Adds an aent in the manifest",
		UsageText: "aenthill [global options] add image",
		Action: func(ctx *cli.Context) error {
			if err := validateArgsLength(ctx, 1, 1); err != nil {
				return errors.Wrap("add command", err)
			}
			job, err := jobs.NewAddJob(ctx.Args().Get(0), context, m)
			if err != nil {
				return errors.Wrap("add command", err)
			}
			return errors.Wrap("add command", job.Execute())
		},
	}
}
