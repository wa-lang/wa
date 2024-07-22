// 版权 @2023 凹语言 作者。保留所有权利。

package applex

import (
	"fmt"
	"os"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/app/appbase"
	"wa-lang.org/wa/internal/scanner"
	"wa-lang.org/wa/internal/token"
	wat_scanner "wa-lang.org/wa/internal/wat/scanner"
	wat_token "wa-lang.org/wa/internal/wat/token"
)

var CmdLex = &cli.Command{
	Hidden: true,
	Name:   "lex",
	Usage:  "lex Wa/wat source code and print token list",
	Action: func(c *cli.Context) error {
		if c.NArg() == 0 {
			fmt.Fprintf(os.Stderr, "no input file")
			os.Exit(1)
		}

		err := Lex(c.Args().First())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return nil
	},
}

func Lex(filename string) error {
	if appbase.IsNativeFile(filename, ".wa", ".wz") {
		return lexWa(filename)
	}
	if appbase.IsNativeFile(filename, ".wat") {
		return lexWat(filename)
	}

	return fmt.Errorf("%q is not valid path", filename)
}

func lexWa(filename string) error {
	src, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile(filename, fset.Base(), len(src))
	s.Init(file, src, nil, scanner.ScanComments)

	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		fmt.Printf("%s\t%s\t%q\n", fset.Position(pos), tok, lit)
	}

	return nil
}

func lexWat(filename string) error {
	src, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var s wat_scanner.Scanner
	file := wat_token.NewFile(filename, len(src))
	s.Init(file, src, nil, wat_scanner.ScanComments)

	for {
		pos, tok, lit := s.Scan()
		if tok == wat_token.EOF {
			break
		}
		fmt.Printf("%s\t%s\t%q\n", file.Position(pos), tok, lit)
	}

	return nil
}
