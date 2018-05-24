// Package commands implements all commands of the application.
package commands

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
	return fmt.Sprintf(manifestFileDoestNotExistErrorMessage, manifest.DefaultManifestFileName, RootCmd.Use, InitCmd.Use)
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
	logLevel   string
	projectDir string

	// RootCmd is our commands entry point.
	RootCmd = &cobra.Command{
		Use:                "aenthill",
		SilenceErrors:      true,
		SilenceUsage:       true,
		DisableSuggestions: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if logLevel != "" {
				l, ok := levels[logLevel]
				if !ok {
					return &wrongLogLevelError{}
				}

				log.SetLevel(l)
			}

			wd, err := os.Getwd()
			if err != nil {
				return err
			}
			projectDir = wd
			log.WithField("path", projectDir).Debug("project directory found")

			return nil
		},
	}
)

func init() {
	RootCmd.PersistentFlags().StringVarP(&logLevel, "logLevel", "l", "", "configure the log level: DEBUG, INFO, WARN, ERROR, FATAL. Default is INFO")
}
