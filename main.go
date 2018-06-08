/*
Package main is the entry point of the application.

TODO
*/
package main

import (
	"os"

	"github.com/aenthill/aenthill/app"

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
	a := app.New(version, m)
	if err := a.Execute(); err != nil {
		os.Exit(1)
	}
}
