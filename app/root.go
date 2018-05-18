// Package app implements all commands of the application.
package app

import (
	"github.com/aenthill/aenthill/app/log"

	"github.com/spf13/cobra"
)

var (
	logLevel string

	// RootCmd is the instance of the root of all commands.
	RootCmd = &cobra.Command{
		Use:           "aenthill",
		Short:         "TODO",
		Long:          "TODO.",
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
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
