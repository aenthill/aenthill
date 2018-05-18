package app

import (
	"fmt"
	"os"

	"github.com/aenthill/aenthill/app/prompt"

	"github.com/aenthill/manifest"
	"github.com/spf13/cobra"
)

type manifestFileAlreadyExistingError struct{}

const manifestFileAlreadyExistingErrorMessage = "manifest %s already exists"

func (e *manifestFileAlreadyExistingError) Error() string {
	return fmt.Sprintf(manifestFileAlreadyExistingErrorMessage, manifest.DefaultManifestFileName)
}

var initCmd = &cobra.Command{
	Use:           "init",
	Short:         "TODO",
	Long:          "TODO.",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat(manifest.DefaultManifestFileName); err == nil {
			return &manifestFileAlreadyExistingError{}
		}

		m, err := prompt.AskManifestValues()
		if err != nil {
			return err
		}

		return manifest.Flush(manifest.DefaultManifestFileName, m)
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}
