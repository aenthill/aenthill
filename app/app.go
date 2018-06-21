/*
Package app is a wrapper around cobra library.

Its main goal is to initialize the application context and validate it.
*/
package app

import (
	"os"

	"github.com/aenthill/aenthill/app/commands"
	"github.com/aenthill/aenthill/app/context"

	"github.com/aenthill/log"
	"github.com/aenthill/manifest"
	"github.com/spf13/cobra"
)

// App is our working struct.
type App struct {
	name     string
	version  string
	manifest *manifest.Manifest
	ctx      *context.AppContext
}

// New creates an App instance with the given Manifest instance.
func New(version string, m *manifest.Manifest) *App {
	return &App{
		name:     "aenthill",
		version:  version,
		manifest: m,
		ctx:      &context.AppContext{},
	}
}

// Execute executes a command from CLI.
func (app *App) Execute() error {
	rootCmd := &cobra.Command{
		Use:     app.name,
		Version: app.version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return app.entrypoint(cmd, args)
		},
	}
	rootCmd.PersistentFlags().BoolVarP(&app.ctx.IsVerbose, "verbose", "v", false, "configures the log level to INFO")
	rootCmd.PersistentFlags().BoolVarP(&app.ctx.IsVeryVerbose, "debug", "d", false, "configures the log level to DEBUG")
	rootCmd.AddCommand(commands.NewInitCmd(app.manifest, app.ctx))
	rootCmd.AddCommand(commands.NewAddCmd(app.manifest, app.ctx))
	rootCmd.AddCommand(commands.NewRemoveCmd(app.manifest, app.ctx))
	rootCmd.AddCommand(commands.NewSelfUpdateCmd(app.version, app.ctx))

	return rootCmd.Execute()
}

func (app *App) entrypoint(cmd *cobra.Command, args []string) error {
	if cmd == nil || cmd.Use == "help [command]" {
		return nil
	}

	err := app.initialize()
	if err != nil {
		log.Error(app.ctx.EntryContext, err, "initialization failed")
	}

	return err
}

func (app *App) initialize() error {
	app.ctx.EntryContext = &log.EntryContext{}

	projectDir, err := os.Getwd()
	if err != nil {
		return err
	}
	app.ctx.ProjectDir = projectDir

	if app.ctx.IsVerbose {
		app.ctx.LogLevel = log.InfoLevel
	}

	if app.ctx.IsVeryVerbose {
		app.ctx.LogLevel = log.DebugLevel
	}

	if app.ctx.LogLevel == "" {
		app.ctx.LogLevel = log.ErrorLevel
	}

	return log.SetLevel(app.ctx.LogLevel)
}
