package commands

import (
	"github.com/aenthill/aenthill/context"
	"github.com/aenthill/aenthill/errors"
	"github.com/aenthill/aenthill/jobs"

	"github.com/urfave/cli"
)

func NewReplyCommand(context *context.Context) cli.Command {
	return cli.Command{
		Name:      "reply",
		Usage:     "Replies to the aent which awakened current aent",
		UsageText: "aenthill reply event [payload]",
		Action: func(ctx *cli.Context) error {
			job, err := jobs.NewReplyJob(ctx.Args().Get(0), ctx.Args().Get(1), context)
			if err != nil {
				return errors.Wrap("reply command", err)
			}
			return errors.Wrap("reply command", job.Execute())
		},
	}
}
