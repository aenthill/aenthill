/*
Package context provides a struct gathering some data
used by our commands.
*/
package context

// AppContext is our working struct.
type AppContext struct {
	LogLevel   string
	ProjectDir string
}
