// 版权 @2023 凹语言 作者。保留所有权利。

package appast

import (
	"fmt"
	"os"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/parser"
	"wa-lang.org/wa/internal/token"
)

var CmdAst = &cli.Command{
	Hidden: true,
	Name:   "ast",
	Usage:  "parse Wa source code and print ast",
	Action: func(c *cli.Context) error {
		if c.NArg() == 0 {
			fmt.Fprintf(os.Stderr, "no input file")
			os.Exit(1)
		}

		err := PrintAST(c.Args().First())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return nil
	},
}

func PrintAST(filename string) error {
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(nil, fset, filename, nil, 0)
	if err != nil {
		return err
	}

	return ast.Print(fset, f)
}
