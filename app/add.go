package app

import (
	"fmt"
	"os"

	"github.com/aenthill/aenthill/app/manifest"

	"github.com/spf13/cobra"
)

type manifestFileDoestNotExistError struct{}

const manifestFileDoestNotExistErrorMessage = "manifest %s not found in current directory. Did you run %s %s?"

func (e *manifestFileDoestNotExistError) Error() string {
	return fmt.Sprintf(manifestFileDoestNotExistErrorMessage, manifest.DefaultManifestFileName, RootCmd.Use, initCmd.Use)
}

var addCmd = &cobra.Command{
	Use:           "add",
	Short:         "TODO",
	Long:          "TODO.",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat(manifest.DefaultManifestFileName); err != nil {
			return &manifestFileDoestNotExistError{}
		}

		if len(args) == 0 {
			return nil
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
}
