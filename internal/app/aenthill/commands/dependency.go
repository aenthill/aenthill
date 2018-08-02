package commands

import (
	"github.com/aenthill/aenthill/internal/pkg/context"
	"github.com/aenthill/aenthill/internal/pkg/errors"
	"github.com/aenthill/aenthill/internal/pkg/jobs"
	"github.com/aenthill/aenthill/internal/pkg/manifest"

	"github.com/urfave/cli"
)

// NewDependencyCommand creates a cli.Command instance.
func NewDependencyCommand(context *context.Context, m *manifest.Manifest) cli.Command {
	cmd := cli.Command{
		Name:      "dependency",
		Usage:     "Prints a dependency ID of current aent",
		UsageText: "aenthill [global options] dependency key",
		Action: func(ctx *cli.Context) error {
			if err := validateArgsLength(ctx, 2, 2); err != nil {
				return errors.Wrap("dependency command", err)
			}
			job, err := jobs.NewDependencyJob(ctx.Args().Get(0), context, m)
			if err != nil {
				errors.Wrap("dependency command", err)
			}
			return errors.Wrap("dependency command", job.Execute())
		},
	}
	return cmd
}
