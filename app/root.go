// Package app implements all commands of the application.
package app

import (
	"fmt"
	"os"

	"github.com/aenthill/log"
	"github.com/aenthill/manifest"

	"github.com/spf13/cobra"
)

type manifestFileDoestNotExistError struct{}

const manifestFileDoestNotExistErrorMessage = "manifest %s not found in current directory. Did you run %s %s?"

func (e *manifestFileDoestNotExistError) Error() string {
	return fmt.Sprintf(manifestFileDoestNotExistErrorMessage, manifest.DefaultManifestFileName, RootCmd.Use, initCmd.Use)
}

var (
	projectDir string

	logLevel string

	// RootCmd is the instance of the root of all commands.
	RootCmd = &cobra.Command{
		Use:           "aenthill",
		Short:         "TODO",
		Long:          "TODO.",
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			projectDir = wd

			if logLevel != "" {
				return log.SetLevel(logLevel)
			}

			return nil
		},
	}
)

func init() {
	RootCmd.PersistentFlags().StringVarP(&logLevel, "logLevel", "l", "", "configure the log level")
}
