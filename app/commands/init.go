package commands

import (
	"fmt"

	"github.com/aenthill/aenthill/app/context"

	"github.com/aenthill/manifest"
	"github.com/apex/log"
	"github.com/spf13/cobra"
)

type manifestFileAlreadyExistingError struct{}

const manifestFileAlreadyExistingErrorMessage = "manifest %s already exists"

func (e *manifestFileAlreadyExistingError) Error() string {
	return fmt.Sprintf(manifestFileAlreadyExistingErrorMessage, manifest.DefaultManifestFileName)
}

// NewInitCmd creates a cobra.Command instance which will use the given
// Manifest and AppContext instances.
func NewInitCmd(m *manifest.Manifest, appCtx *context.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:           "init",
		Short:         fmt.Sprintf("Create the manifest %s in current directory", manifest.DefaultManifestFileName),
		Long:          fmt.Sprintf("Create the manifest %s in current directory", manifest.DefaultManifestFileName),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if m.Exist() {
				return &manifestFileAlreadyExistingError{}
			}

			err := m.Flush()
			if err == nil {
				log.Infof("%s created! May the swarm be with you", manifest.DefaultManifestFileName)
			}

			return err
		},
	}
}
