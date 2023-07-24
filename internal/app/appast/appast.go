// 版权 @2023 凹语言 作者。保留所有权利。

package appast

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/parser"
	"wa-lang.org/wa/internal/token"
)

var CmdAst = &cli.Command{
	Hidden: true,
	Name:   "ast",
	Usage:  "parse Wa source code and print ast",
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
	ext := strings.ToLower(filepath.Ext(filename))
	if ext != ".wa" && ext != ".wz" {
		return fmt.Errorf("%q is not Wa file", filename)
	}
	if fi, _ := os.Lstat(filename); fi == nil || fi.IsDir() {
		return fmt.Errorf("%q not found", filename)
	}

	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(nil, fset, filename, nil, 0)
	if err != nil {
		return err
	}

	return ast.Print(fset, f)
}
