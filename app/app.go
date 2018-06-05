/*
Package app is a wrapper around cobra library.

Its main goal is to initialize the application context and validate it.
*/
package app

import (
	"os"

	"github.com/aenthill/aenthill/app/commands"
	"github.com/aenthill/aenthill/app/context"

	"github.com/aenthill/manifest"
	"github.com/apex/log"
	"github.com/spf13/cobra"
)

// App is our working struct.
type App struct {
	version  string
	manifest *manifest.Manifest
	ctx      *context.AppContext
}

// New creates an App instance with the given Manifest instance.
func New(m *manifest.Manifest, version string) *App {
	return &App{
		version:  version,
		manifest: m,
		ctx:      &context.AppContext{},
	}
}

// Execute executes a command from CLI.
func (app *App) Execute() error {
	rootCmd := &cobra.Command{
		Use:                "aenthill",
		Version:            app.version,
		SilenceErrors:      true,
		SilenceUsage:       true,
		DisableSuggestions: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return app.initialize()
		},
	}
	rootCmd.PersistentFlags().StringVarP(&app.ctx.LogLevel, "logLevel", "l", "", "configure the log level: DEBUG, INFO, WARN, ERROR. Default is INFO")
	rootCmd.AddCommand(commands.NewInitCmd(app.manifest, app.ctx))
	rootCmd.AddCommand(commands.NewAddCmd(app.manifest, app.ctx))
	rootCmd.AddCommand(commands.NewRemoveCmd(app.manifest, app.ctx))

	return rootCmd.Execute()
}

// levels associates log levels as used with the --logLevel -l flag from aenthill
// with its counterpart from the github.com/apex/log library.
var levels = map[string]log.Level{
	"DEBUG": log.DebugLevel,
	"INFO":  log.InfoLevel,
	"WARN":  log.WarnLevel,
	"ERROR": log.ErrorLevel,
}

type wrongLogLevelError struct{}

const wrongLogLevelErrorMessage = "accepted values for log level: DEBUG, INFO, WARN, ERROR"

func (e *wrongLogLevelError) Error() string {
	return wrongLogLevelErrorMessage
}

func (app *App) initialize() error {
	projectDir, err := os.Getwd()
	if err != nil {
		return err
	}
	app.ctx.ProjectDir = projectDir

	if app.ctx.LogLevel != "" {
		l, ok := levels[app.ctx.LogLevel]
		if !ok {
			return &wrongLogLevelError{}
		}
		log.SetLevel(l)
	} else {
		app.ctx.LogLevel = "INFO"
	}

	return nil
}
