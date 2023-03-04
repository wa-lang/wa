// 版权 @2023 凹语言 作者。保留所有权利。

package app

import (
	"fmt"

	"wa-lang.org/wa/internal/scanner"
	"wa-lang.org/wa/internal/token"
)

func (p *App) Lex(filename string) error {
	src, err := p.readSource(filename, nil)
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
