// Package tests contains useful functions used across tests.
package tests

import (
	"os"
	"testing"

	"github.com/aenthill/aenthill/context"
)

// MakeTestContext creates a new context for tests.
func MakeTestContext(t *testing.T) *context.Context {
	ctx, err := context.New()
	if err != nil {
		t.Errorf(`An unexpected error occurred while creating the context: got "%s"`, err.Error())
	}
	ctx.HostProjectDir = os.Getenv("HOST_PROJECT_DIR")
	return ctx
}
