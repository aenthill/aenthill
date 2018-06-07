/*
Package main is the entry point of the application.

TODO
*/
package main

import (
	"os"
	"time"

	"github.com/aenthill/aenthill/app"

	"github.com/aenthill/log"
	"github.com/aenthill/manifest"
	"github.com/spf13/afero"
)

/*
version will be set by GoReleaser.
It will be the current Git tag (with v prefix stripped) or
the name of the snapshot if you're using the --snapshot flag.
*/
var version = "master"

func main() {
	m := manifest.New(manifest.DefaultManifestFileName, afero.NewOsFs())
	a := app.New(m, version)
	entryCtx := &log.EntryContext{Source: "aenthill"}
	start := time.Now()

	if err := a.Execute(); err != nil {
		log.Errorf(entryCtx, err, "job has failed after %0.2fs", time.Since(start).Seconds())
		os.Exit(1)
	}

	if shouldDisplayTime() {
		log.Infof(entryCtx, "job has successfully finished after %0.2fs", time.Since(start).Seconds())
	}
}

func shouldDisplayTime() bool {
	hasCommand, hasHelpFlag := false, false
	for _, arg := range os.Args {
		// we ignore init command as it is no relevant.
		if arg == "add" || arg == "rm" {
			hasCommand = true
		}

		if arg == "-h" || arg == "--help" {
			hasHelpFlag = true
		}
	}

	return hasCommand && !hasHelpFlag
}
