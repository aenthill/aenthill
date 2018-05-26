package commands

import (
	"fmt"
	"os"

	"github.com/aenthill/docker"
	"github.com/aenthill/manifest"
	"github.com/apex/log"
	"github.com/spf13/cobra"
)

type noImagesToRemoveError struct{}

const noImagesToRemoveErrorMessage = "usage: %s %s image [image...]"

func (e *noImagesToRemoveError) Error() string {
	return fmt.Sprintf(noImagesToRemoveErrorMessage, RootCmd.Use, RemoveCmd.Use)
}

/*
RemoveCmd removes one or more aents.

It removes the given aents from manifest and sends a "REMOVE" event to them.

Usage:

 aenthill rm image [image...] [flags]
*/
var RemoveCmd = &cobra.Command{
	Use:           "rm",
	Short:         "Remove one or more aents",
	Long:          "Remove one or more aents",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat(manifest.DefaultManifestFileName); err != nil {
			return &manifestFileDoesNotExistError{}
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
			log.WithField("aent", fmt.Sprintf("%s", image)).Info("removed aent from manifest")

			ctx := &docker.EventContext{
				WhoAmI:         RootCmd.Use,
				Image:          image,
				Binary:         docker.DefaultBinary,
				HostProjectDir: projectDir,
				LogLevel:       logLevel,
			}

			if err := docker.Send("REMOVE", "", ctx); err != nil {
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
	RootCmd.AddCommand(RemoveCmd)
}
