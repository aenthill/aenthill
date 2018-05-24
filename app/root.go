// Package app implements all commands of the application.
package app

import (
	"fmt"
	"os"

	"github.com/aenthill/manifest"

	"github.com/apex/log"
	"github.com/spf13/cobra"
)

type manifestFileDoestNotExistError struct{}

const manifestFileDoestNotExistErrorMessage = "manifest %s not found in current directory. Did you run %s %s?"

func (e *manifestFileDoestNotExistError) Error() string {
	return fmt.Sprintf(manifestFileDoestNotExistErrorMessage, manifest.DefaultManifestFileName, RootCmd.Use, initCmd.Use)
}

// levels associates log levels as used with the --logLevel -l flag
// with its counterpart from the github.com/apex/log library.
var levels = map[string]log.Level{
	"DEBUG": log.DebugLevel,
	"INFO":  log.InfoLevel,
	"WARN":  log.WarnLevel,
	"ERROR": log.ErrorLevel,
	"FATAL": log.FatalLevel,
}

type wrongLogLevelError struct{}

const wrongLogLevelErrorMessage = "accepted values for log level: DEBUG, INFO, WARN, ERROR, FATAL"

func (e *wrongLogLevelError) Error() string {
	return wrongLogLevelErrorMessage
}

var (
	projectDir string

	logLevel string

	// RootCmd is the instance of the root of all commands.
	RootCmd = &cobra.Command{
		Use:           "aenthill",
		Short:         "aenthill is a CLI tool for managing aents",
		Long:          "aenthill is a CLI tool for managing aents",
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			wd, err := os.Getwd()
			if err != nil {
				return err
			}
			projectDir = wd

			if logLevel != "" {
				l, ok := levels[logLevel]
				if !ok {
					return &wrongLogLevelError{}
				}

				log.SetLevel(l)
			}

			return nil
		},
	}
)

func init() {
	RootCmd.PersistentFlags().StringVarP(&logLevel, "logLevel", "l", "", "configure the log level")
}
