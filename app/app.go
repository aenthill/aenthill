/*
Package app is a wrapper around urfave/cli package.

Its main goal is to initialize the application context and validate it.
*/
package app

import (
	"os"

	"github.com/aenthill/aenthill/app/commands"
	"github.com/aenthill/aenthill/context"
	"github.com/aenthill/aenthill/errors"
	"github.com/aenthill/aenthill/manifest"

	"github.com/urfave/cli"
)

// App is our working struct.
type App struct {
	cli      *cli.App
	ctx      *context.Context
	manifest *manifest.Manifest
}

// New creates an App instance with the given Manifest instance.
func New(version string, m *manifest.Manifest) (*App, error) {
	ctx, err := context.New()
	if err != nil {
		return nil, errors.Wrap("app", err)
	}
	app := &App{cli: cli.NewApp(), ctx: ctx, manifest: m}
	app.cli.Name, app.cli.Usage, app.cli.Version = "aenthill", "May the swarm be with you!", version
	app.registerCommands()
	return app, nil
}

func (app *App) registerCommands() {
	if app.ctx.IsContainer() {
		app.cli.Commands = append(app.cli.Commands, commands.NewRunCommand(app.ctx, app.manifest))
		app.cli.Commands = append(app.cli.Commands, commands.NewInstallCommand(app.ctx, app.manifest))
		app.cli.Commands = append(app.cli.Commands, commands.NewRegisterCommand(app.ctx, app.manifest))
		app.cli.Commands = append(app.cli.Commands, commands.NewDispatchCommand(app.ctx, app.manifest))
		app.cli.Commands = append(app.cli.Commands, commands.NewReplyCommand(app.ctx))
	} else {
		app.cli.Commands = append(app.cli.Commands, commands.NewAddCommand(app.ctx, app.manifest))
		app.cli.Commands = append(app.cli.Commands, commands.NewUpgradeCommand(app.cli.Version))
	}
}

// Run starts the application by reading CLI arguments.
func (app *App) Run() error {
	return errors.Wrap("app", app.cli.Run(os.Args))
}
