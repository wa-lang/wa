// 版权 @2023 凹语言 作者。保留所有权利。

package appast

import (
	"fmt"
	"os"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/app/appbase"
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/parser"
	"wa-lang.org/wa/internal/token"
	wat_parser "wa-lang.org/wa/internal/wat/parser"
)

var CmdAst = &cli.Command{
	Hidden: true,
	Name:   "ast",
	Usage:  "parse Wa/wat source code and print ast",
	Action: CmdAstAction,
}

func CmdAstAction(c *cli.Context) error {
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
}

func PrintAST(filename string) error {
	if appbase.IsNativeFile(filename, ".wat") {
		return PrintWatAST(filename)
	}
	if !appbase.HasExt(filename, ".wa", ".wz") {
		return fmt.Errorf("%q is not Wa file", filename)
	}

	if !appbase.PathExists(filename) {
		return fmt.Errorf("%q not found", filename)
	}
	if !appbase.IsNativeFile(filename) {
		return fmt.Errorf("%q must be file", filename)
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(nil, fset, filename, nil, 0)
	if err != nil {
		return err
	}

	return ast.Print(fset, f)
}

func PrintWatAST(filename string) error {
	m, err := wat_parser.ParseModule(filename, nil)
	if err != nil {
		return err
	}

	fmt.Println(m)
	return nil
}
