package app

import (
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:           "add",
	Short:         "TODO",
	Long:          "TODO.",
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
}
