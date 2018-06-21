/*
Package context provides a struct gathering some data
used by our commands.
*/
package context

import "github.com/aenthill/log"

// AppContext is our working struct.
type AppContext struct {
	ProjectDir    string
	IsVerbose     bool
	IsVeryVerbose bool
	LogLevel      string
	EntryContext  *log.EntryContext
}
