package commands

import (
	"fmt"

	"github.com/aenthill/aenthill/app/context"

	"github.com/aenthill/docker"
	"github.com/aenthill/manifest"
	"github.com/apex/log"
	"github.com/spf13/cobra"
)

type noImagesToAddError struct{}

const noImagesToAddErrorMessage = "usage: aenthill add image [image...]"

func (e *noImagesToAddError) Error() string {
	return noImagesToAddErrorMessage
}

func NewAddCmd(m *manifest.Manifest, appCtx *context.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:           "add",
		Short:         "Add one or more aents",
		Long:          "Add one or more aents",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return &noImagesToAddError{}
			}

			if !m.Exist() {
				return &manifestFileDoesNotExistError{}
			}

			if err := m.Parse(); err != nil {
				return err
			}

			for _, image := range args {
				ctx := &docker.EventContext{
					WhoAmI:         "aenthill",
					Image:          image,
					Binary:         docker.DefaultBinary,
					HostProjectDir: appCtx.ProjectDir,
					LogLevel:       appCtx.LogLevel,
				}

				if err := docker.Send("ADD", "", ctx); err != nil {
					return err
				}

				if err := m.AddAent(image); err == nil {
					if err := m.Flush(); err != nil {
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
}
