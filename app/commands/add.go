package commands

import (
	"github.com/aenthill/aenthill/app/context"
	"github.com/aenthill/aenthill/app/jobs"

	"github.com/aenthill/log"
	"github.com/aenthill/manifest"
	"github.com/spf13/cobra"
)

// NewAddCmd creates a cobra.Command instance which will use the given
// Manifest and AppContext instances.
func NewAddCmd(m *manifest.Manifest, appCtx *context.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "add",
		Short:         "Adds one or more aents",
		Long:          "Adds one or more aents",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			job, err := jobs.NewAddJob(args, m, appCtx)
			if err != nil {
				log.Error(appCtx.EntryContext, err, "job initialization failed")
				return err
			}

			return job.Run()
		},
	}
	cmd.SetUsageTemplate(getUsageTemplate("aenthill add image [images...] [flags]"))

	return cmd
}
