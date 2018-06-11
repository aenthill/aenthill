package commands

import (
	"fmt"

	"github.com/aenthill/aenthill/app/context"
	"github.com/aenthill/aenthill/app/jobs"

	"github.com/aenthill/log"
	"github.com/aenthill/manifest"
	"github.com/spf13/cobra"
)

// NewInitCmd creates a cobra.Command instance which will use the given
// Manifest and AppContext instances.
func NewInitCmd(m *manifest.Manifest, appCtx *context.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:           "init",
		Short:         fmt.Sprintf("Creates the manifest %s in current directory", manifest.DefaultManifestFileName),
		Long:          fmt.Sprintf("Creates the manifest %s in current directory", manifest.DefaultManifestFileName),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			job, err := jobs.NewInitJob(m, appCtx)
			if err != nil {
				log.Error(appCtx.EntryContext, err, "job initialization failed")
				return err
			}

			return job.Run()
		},
	}
}
