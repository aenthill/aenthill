package commands

import (
	"github.com/aenthill/aenthill/app/context"

	"github.com/aenthill/docker"
	"github.com/aenthill/manifest"
	"github.com/apex/log"
	"github.com/spf13/cobra"
)

type noImagesToRemoveError struct{}

const noImagesToRemoveErrorMessage = "usage: aenthill rm image [image...]"

func (e *noImagesToRemoveError) Error() string {
	return noImagesToRemoveErrorMessage
}

// NewRemoveCmd creates a cobra.Command instance which will use the given
// Manifest and AppContext instances.
func NewRemoveCmd(m *manifest.Manifest, appCtx *context.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:           "rm",
		Short:         "Remove one or more aents",
		Long:          "Remove one or more aents",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return &noImagesToRemoveError{}
			}

			if !m.Exist() {
				return &manifestFileDoesNotExistError{}
			}

			if err := m.Parse(); err != nil {
				return err
			}

			for _, image := range args {
				ctx := &docker.EventContext{
					From:           "aenthill",
					To:             image,
					HostProjectDir: appCtx.ProjectDir,
					LogLevel:       appCtx.LogLevel,
				}

				if err := docker.Execute("REMOVE", "", ctx); err != nil {
					return err
				}

				if err := m.RemoveAent(image); err != nil {
					return err
				}
				log.WithField("aent", image).Info("removed aent from manifest")

				if err := m.Flush(); err != nil {
					return err
				}
			}

			return nil
		},
	}
}
