package app

import (
	"fmt"
	"os"

	"github.com/aenthill/docker"
	"github.com/aenthill/manifest"
	"github.com/spf13/cobra"
)

type noImagesToRemoveError struct{}

const noImagesToRemoveErrorMessage = "usage: %s %s image [image...]"

func (e *noImagesToRemoveError) Error() string {
	return fmt.Sprintf(noImagesToRemoveErrorMessage, RootCmd.Use, removeCmd.Use)
}

var removeCmd = &cobra.Command{
	Use:           "rm",
	Short:         "remove one or more aents",
	Long:          "remove one or more aents",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat(manifest.DefaultManifestFileName); err != nil {
			return &manifestFileDoestNotExistError{}
		}

		if len(args) == 0 {
			return &noImagesToRemoveError{}
		}

		m, err := manifest.Parse(manifest.DefaultManifestFileName)
		if err != nil {
			return err
		}

		for _, image := range args {
			if err := manifest.RemoveAent(image, m); err != nil {
				return err
			}

			ctx := &docker.EventContext{
				Image:          image,
				Binary:         docker.DefaultBinary,
				HostProjectDir: projectDir,
			}

			if err := docker.Send("REMOVE", ctx); err != nil {
				return err
			}

			if err := manifest.Flush(manifest.DefaultManifestFileName, m); err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(removeCmd)
}
