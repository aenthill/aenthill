package commands

import (
	"fmt"
	"os"

	"github.com/aenthill/docker"
	"github.com/aenthill/manifest"
	"github.com/apex/log"
	"github.com/spf13/cobra"
)

type noImagesToAddError struct{}

const noImagesToAddErrorMessage = "usage: %s %s image [image...]"

func (e *noImagesToAddError) Error() string {
	return fmt.Sprintf(noImagesToAddErrorMessage, RootCmd.Use, AddCmd.Use)
}

/*
AddCmd adds one or more aents.

It adds the given aents in manifest (if they do not already exist)
and sends an "ADD" event to them.

Usage:

 aenthill add image [image...] [flags]
*/
var AddCmd = &cobra.Command{
	Use:           "add",
	Short:         "Add one or more aents",
	Long:          "Add one or more aents",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat(manifest.DefaultManifestFileName); err != nil {
			return &manifestFileDoesNotExistError{}
		}

		if len(args) == 0 {
			return &noImagesToAddError{}
		}

		m, err := manifest.Parse(manifest.DefaultManifestFileName)
		if err != nil {
			return err
		}

		for _, image := range args {
			ctx := &docker.EventContext{
				WhoAmI:         RootCmd.Use,
				Image:          image,
				Binary:         docker.DefaultBinary,
				HostProjectDir: projectDir,
				LogLevel:       logLevel,
			}

			if err := docker.Send("ADD", "", ctx); err != nil {
				return err
			}

			if err := manifest.AddAent(image, m); err == nil {
				if err := manifest.Flush(manifest.DefaultManifestFileName, m); err != nil {
					return err
				}
				log.WithField("aent", fmt.Sprintf("%s", image)).Info("added new aent in manifest")
			} else {
				log.WithField("aent", fmt.Sprintf("%s", image)).Info("aent already in manifest")
			}
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(AddCmd)
}
