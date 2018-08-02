// Package commands gather all the commands the application.
package commands

import (
	"fmt"

	"github.com/urfave/cli"
)

func validateArgsLength(ctx *cli.Context, min int, max int) error {
	nargs := ctx.NArg()
	if min == 0 && max == 0 {
		if nargs > 0 {
			return fmt.Errorf(`command does not need any arguments: got "%d"`, nargs)
		}
		return nil
	}
	if nargs < min {
		return fmt.Errorf(`command requires at leat "%d" argument(s): got "%d"`, min, nargs)
	}
	if nargs > max {
		return fmt.Errorf(`command requires a maximum of "%d" argument(s): got "%d"`, max, nargs)
	}
	return nil
}
