package commands

import (
	"github.com/aenthill/aenthill/context"
	"github.com/aenthill/aenthill/errors"
	"github.com/aenthill/aenthill/jobs"
	"github.com/aenthill/aenthill/manifest"

	"github.com/urfave/cli"
)

// NewDependencyCommand creates a cli.Command instance.
func NewDependencyCommand(context *context.Context, m *manifest.Manifest) cli.Command {
	cmd := cli.Command{
		Name:      "dependency",
		Usage:     "Prints one or all dependencies of current aent",
		UsageText: "aenthill [global options] dependency [key]",
		Action: func(ctx *cli.Context) error {
			job, err := jobs.NewDependencyJob(ctx.Args().Get(0), context, m)
			if err != nil {
				errors.Wrap("dependency command", err)
			}
			return errors.Wrap("dependency command", job.Execute())
		},
	}
	return cmd
}
