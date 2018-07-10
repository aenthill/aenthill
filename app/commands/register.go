package commands

import (
	"github.com/aenthill/aenthill/errors"
	"github.com/aenthill/aenthill/jobs"
	"github.com/aenthill/aenthill/manifest"

	"github.com/urfave/cli"
)

func NewRegisterCommand(m *manifest.Manifest) cli.Command {
	cmd := cli.Command{
		Name:      "register",
		Usage:     "Registers an aent in the manifest",
		UsageText: "aenthill register image [command options]",
		Action: func(ctx *cli.Context) error {
			job := jobs.NewRegisterJob(ctx.Args().Get(0), ctx.String("register-key"), ctx.StringSlice("events"), ctx.StringSlice("metadata"), m)
			return errors.Wrap("register command", job.Execute())
		},
	}
	cmd.Flags = []cli.Flag{
		cli.StringFlag{Name: "register-key-as", Usage: `set the given environment variable prefixed with "PHEROMONE_METADATA_" with the key from the registred aent`},
		cli.StringSliceFlag{Name: "events, e", Usage: "add one handled event (cumulative)"},
		cli.StringSliceFlag{Name: "metadata, m", Usage: "add one metadata (cumulative) - format: key=value"},
	}
	return cmd
}
