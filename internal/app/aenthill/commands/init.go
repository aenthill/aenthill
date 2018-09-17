package commands

import (
	"github.com/aenthill/aenthill/internal/pkg/context"
	"github.com/aenthill/aenthill/internal/pkg/errors"
	"github.com/aenthill/aenthill/internal/pkg/jobs"
	"github.com/aenthill/aenthill/internal/pkg/manifest"

	"github.com/urfave/cli"
)

// NewInitCommand creates a cli.Command instance.
func NewInitCommand(context *context.Context, m *manifest.Manifest) cli.Command {
	return cli.Command{
		Name:      "init",
		Aliases:   []string{"i"},
		Usage:     "Initializes a Docker project for your web application",
		UsageText: "aenthill [global options] init",
		Action: func(ctx *cli.Context) error {
			if err := validateArgsLength(ctx, 0, 0); err != nil {
				return errors.Wrap("init command", err)
			}
			job, err := jobs.NewInitJob(context, m)
			if err != nil {
				return errors.Wrap("init command", err)
			}
			return errors.Wrap("init command", job.Execute())
		},
	}
}
