// 版权 @2023 凹语言 作者。保留所有权利。

package appgo2wa

import (
	"fmt"
	"os"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/app/appbase"
	"wa-lang.org/wa/internal/format"
	"wa-lang.org/wa/internal/parser"
	"wa-lang.org/wa/internal/token"
)

var CmdGo2wa = &cli.Command{
	Hidden: true,
	Name:   "go2wa",
	Usage:  "convert wago to wa",
	Action: CmdGo2waAction,
}

func CmdGo2waAction(c *cli.Context) error {
	if c.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "no input file")
		os.Exit(1)
	}

	code, err := go2wa(c.Args().First())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Print(string(code))
	return nil
}

func go2wa(filename string) ([]byte, error) {
	if !appbase.HasExt(filename, ".go") {
		return nil, fmt.Errorf("%q is not Go file", filename)
	}

	if !appbase.PathExists(filename) {
		return nil, fmt.Errorf("%q not found", filename)
	}
	if !appbase.IsNativeFile(filename) {
		return nil, fmt.Errorf("%q must be file", filename)
	}

	src, _ := os.ReadFile(filename)

	fset := token.NewFileSet()
	f, err := parser.ParseFile(nil, fset, filename, nil, 0)
	if err != nil {
		return nil, err
	}

	f.Name.Name = ""
	return format.DevFormat(fset, f, src)
}
