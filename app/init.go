package app

import (
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:           "init",
	Short:         "TODO",
	Long:          "TODO.",
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}
