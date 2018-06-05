/*
Package context provides a struct gathering some data
used by our commands.
*/
package context

// AppContext is our working struct.
type AppContext struct {
	ProjectDir string
	LogLevel   string
}
