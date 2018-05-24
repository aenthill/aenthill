/*
Package main handles the application startup.

TODO
*/
package main

import (
	"os"
	"time"

	"github.com/aenthill/aenthill/app"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
)

/*
version will be set by GoReleaser.
It will be the current Git tag (with v prefix stripped) or
the name of the snapshot if you're using the --snapshot flag.
*/
//var version = "master"

func init() {
	log.SetHandler(cli.Default)
}

func main() {
	start := time.Now()

	if err := app.RootCmd.Execute(); err != nil {
		log.WithError(err).Errorf("aenthill job failed after %0.2fs", time.Since(start).Seconds())
		os.Exit(1)
		return
	}

	log.Infof("aenthill job finished after %0.2fs", time.Since(start).Seconds())
}
