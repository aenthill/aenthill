package commands

import (
	"github.com/aenthill/aenthill/errors"
	"github.com/aenthill/aenthill/jobs"

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
			job := jobs.NewUpgradeJob(ctx.String("target"), version)
			return errors.Wrap("upgrade command", job.Execute())
		},
	}
	cmd.Flags = []cli.Flag{
		cli.StringFlag{Name: "target, t", Usage: "specify the targeted version"},
	}
	return cmd
}
