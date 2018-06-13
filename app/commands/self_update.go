package commands

import (
	"github.com/aenthill/aenthill/app/context"
	"github.com/aenthill/aenthill/app/jobs"

	"github.com/spf13/cobra"
)

// NewSelfUpdateCmd creates a cobra.Command instance which will use the given
// version and AppContext instance.
func NewSelfUpdateCmd(version string, appCtx *context.AppContext) *cobra.Command {
	var target string

	cmd := &cobra.Command{
		Use:           "self:update",
		Short:         "Updates the current version of aenthill",
		Long:          "Updates the current version of aenthill",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			job := jobs.NewSelfUpdateJob(target, version, appCtx)
			return job.Run()
		},
	}
	cmd.Flags().StringVarP(&target, "target", "t", "", "specifies the target version")

	return cmd
}
