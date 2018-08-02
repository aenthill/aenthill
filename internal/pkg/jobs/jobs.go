// Package jobs provides core logic of the commands of the application.
package jobs

// Job is a job interface.
type Job interface {
	Execute() error
}
