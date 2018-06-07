package commands

import (
	"github.com/aenthill/aenthill/app/context"
	"github.com/aenthill/aenthill/app/jobs"

	"github.com/aenthill/manifest"
	"github.com/spf13/cobra"
)

// NewRemoveCmd creates a cobra.Command instance which will use the given
// Manifest and AppContext instances.
func NewRemoveCmd(m *manifest.Manifest, appCtx *context.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:           "rm",
		Short:         "Removes one or more aents",
		Long:          "Removes one or more aents",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			job, err := jobs.NewRemoveJob(args, m, appCtx)
			if err != nil {
				return err
			}

			return job.Run()
		},
	}
}
