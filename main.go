/*
Package main handles the application startup.

TODO
*/
package main

import (
	"os"

	"github.com/anthill-docker/anthill/app"
	"github.com/anthill-docker/anthill/app/log"
)

/*
version will be set by GoReleaser.
It will be the current Git tag (with v prefix stripped) or
the name of the snapshot if you're using the --snapshot flag.
*/
//var version = "master"

// main initializes the application and starts it.
func main() {
	if err := app.RootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
