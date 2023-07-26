package appdev

import (
	"fmt"
	"os"

	"wa-lang.org/wa/api"
	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/config"
)

var CmdDev = &cli.Command{
	Hidden: true,
	Name:   "debug",
	Usage:  "only for dev/debug",
	Action: func(c *cli.Context) error {
		wat, err := api.BuildFile(
			config.DefaultConfig(),
			"hello.wa", "func main() { println(123) }",
		)
		if err != nil {
			if len(wat) != 0 {
				fmt.Println(string(wat))
			}
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(wat))
		return nil
	},
}
