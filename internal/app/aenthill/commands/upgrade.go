package commands

import (
	"github.com/aenthill/aenthill/internal/pkg/errors"
	"github.com/aenthill/aenthill/internal/pkg/jobs"

	"github.com/urfave/cli"
)

// NewUpgradeCommand creates a cli.Command instance.
func NewUpgradeCommand(version string) cli.Command {
	cmd := cli.Command{
		Name:      "upgrade",
		Aliases:   []string{"u"},
		Usage:     "Upgrades Aenthill",
		UsageText: "aenthill [global options] upgrade [command options]",
		Action: func(ctx *cli.Context) error {
			if err := validateArgsLength(ctx, 0, 0); err != nil {
				return errors.Wrap("upgrade command", err)
			}
			job := jobs.NewUpgradeJob(ctx.String("target"), version)
			return errors.Wrap("upgrade command", job.Execute())
		},
	}
	cmd.Flags = []cli.Flag{
		cli.StringFlag{Name: "target, t", Usage: "specify the targeted version"},
	}
	return cmd
}
