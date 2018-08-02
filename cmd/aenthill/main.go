/*
Aenthill is a solution for allowing Docker containers to communicate with each other
in order to fulfill a task in the current host directory.

Each of those Docker containers (called aent) is not necessarily aware of its peers.
It is however able to receive and communicate events to them.
In other words, an aent is not necessarily smart but
it shares a part of the task to handle like an ant within its colony.

Aenthill documentation is hosted at https://aenthill.github.io/.
*/
package main

import (
	"os"

	app "github.com/aenthill/aenthill/internal/app/aenthill"
	"github.com/aenthill/aenthill/internal/pkg/log"
	"github.com/aenthill/aenthill/internal/pkg/manifest"

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
	a, err := app.New(version, m)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	if err := a.Run(); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}
