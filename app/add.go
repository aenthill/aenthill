package app

import (
	"fmt"
	"os"

	"github.com/aenthill/docker"
	"github.com/aenthill/manifest"
	"github.com/spf13/cobra"
)

type manifestFileDoestNotExistError struct{}

const manifestFileDoestNotExistErrorMessage = "manifest %s not found in current directory. Did you run %s %s?"

func (e *manifestFileDoestNotExistError) Error() string {
	return fmt.Sprintf(manifestFileDoestNotExistErrorMessage, manifest.DefaultManifestFileName, RootCmd.Use, initCmd.Use)
}

type noImagesToAddError struct{}

const noImagesToAddErrorMessage = "usage: %s %s image [image...]"

func (e *noImagesToAddError) Error() string {
	return fmt.Sprintf(noImagesToAddErrorMessage, RootCmd.Use, addCmd.Use)
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
			return &noImagesToAddError{}
		}

		m, err := manifest.Parse(manifest.DefaultManifestFileName)
		if err != nil {
			return err
		}

		for _, image := range args {
			if err := manifest.AddAent(image, m); err != nil {
				return err
			}

			ctx := &docker.EventContext{
				Image:          image,
				Binary:         docker.DefaultBinary,
				HostProjectDir: projectDir,
			}

			if err := docker.Send("ADD", ctx); err != nil {
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
	RootCmd.AddCommand(addCmd)
}
