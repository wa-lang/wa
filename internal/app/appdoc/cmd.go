// 版权 @2023 凹语言 作者。保留所有权利。

package appdoc

import (
	"fmt"
	"os"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/app/appbase"
	"wa-lang.org/wa/internal/config"
	"wa-lang.org/wa/internal/loader"
)

var CmdDoc = &cli.Command{
	Hidden:    true,
	Name:      "doc",
	Usage:     "show documentation for package or symbol",
	ArgsUsage: "name...",
	Flags: []cli.Flag{
		appbase.MakeFlag_target(),
		appbase.MakeFlag_tags(),
	},
	Action: func(c *cli.Context) error {
		opt := appbase.BuildOptions(c)

		var pkgpath = "."
		var names []string
		if c.Args().Len() > 0 {
			pkgpath = c.Args().First()
		}
		if c.Args().Len() > 1 {
			names = c.Args().Slice()[1:]
		}

		RunDoc(opt.Config(), pkgpath, names...)
		return nil
	},
}

func RunDoc(cfg *config.Config, pkgpath string, names ...string) {
	prog, err := loader.LoadProgram(cfg, pkgpath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	doc := BuildPkgDoc(prog)
	doc.Show(names...)
}
