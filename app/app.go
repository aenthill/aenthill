/*
Package app is a wrapper around urfave/cli package.

Its main goal is to initialize the application context and validate it.
*/
package app

import (
	"fmt"
	"os"

	"github.com/aenthill/aenthill/app/commands"
	"github.com/aenthill/aenthill/context"
	"github.com/aenthill/aenthill/errors"
	"github.com/aenthill/aenthill/log"
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
	app.setup(version)
	return app, nil
}

func (app *App) setup(version string) {
	app.cli.Name, app.cli.Usage, app.cli.Version = "aenthill", "May the swarm be with you!", version
	log.SetLevel(app.ctx.LogLevel)
	if app.ctx.IsContainer() {
		app.manifest.SetPath(fmt.Sprintf("%s/%s", app.ctx.ProjectDir, manifest.DefaultManifestFileName))
	}
	app.registerFlags()
	app.registerCommands()
}

func (app *App) registerFlags() {
	cli.VersionFlag = cli.BoolFlag{Name: "version", Usage: "print the version"}
	app.cli.Flags = []cli.Flag{
		cli.BoolFlag{Name: "verbose, v", Usage: "print verbose output to the console"},
		cli.BoolFlag{Name: "debug, d", Usage: "print debug output to the console"},
	}
	app.cli.Before = func(ctx *cli.Context) error {
		if ctx.GlobalBool("verbose") {
			app.ctx.LogLevel = "INFO"
			log.SetLevel("INFO")
		}
		if ctx.GlobalBool("debug") {
			app.ctx.LogLevel = "DEBUG"
			log.SetLevel("DEBUG")
		}
		return nil
	}
}

func (app *App) registerCommands() {
	if app.ctx.IsContainer() {
		app.cli.Commands = append(app.cli.Commands, commands.NewUpdateCommand(app.ctx, app.manifest))
		app.cli.Commands = append(app.cli.Commands, commands.NewRegisterCommand(app.ctx, app.manifest))
		app.cli.Commands = append(app.cli.Commands, commands.NewRunCommand(app.ctx, app.manifest))
		app.cli.Commands = append(app.cli.Commands, commands.NewDispatchCommand(app.ctx, app.manifest))
		app.cli.Commands = append(app.cli.Commands, commands.NewReplyCommand(app.ctx, app.manifest))
	} else {
		app.cli.Commands = append(app.cli.Commands, commands.NewStartCommand(app.ctx, app.manifest))
		app.cli.Commands = append(app.cli.Commands, commands.NewAddCommand(app.ctx, app.manifest))
		app.cli.Commands = append(app.cli.Commands, commands.NewUpgradeCommand(app.cli.Version))
	}
}

// Run starts the application by reading CLI arguments.
func (app *App) Run() error {
	return errors.Wrap("app", app.cli.Run(os.Args))
}
