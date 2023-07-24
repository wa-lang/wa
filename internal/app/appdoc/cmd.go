package appdoc

import (
	"fmt"

	"wa-lang.org/wa/internal/3rdparty/cli"
)

var CmdDoc = &cli.Command{
	Hidden: true,
	Name:   "doc",
	Usage:  "show documentation for package or symbol",
	Action: func(c *cli.Context) error {
		fmt.Println("TODO")
		return nil
	},
}
