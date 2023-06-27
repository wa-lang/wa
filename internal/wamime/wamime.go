// 版权 @2023 凹语言 作者。保留所有权利。

package wamime

import (
	"strings"

	"wa-lang.org/wa/internal/scanner"
	"wa-lang.org/wa/internal/token"
)

const mimePrefix = "#syntax=" // #syntax=wa, #syntax=wz

func GetCodeMime(filename string, code []byte) string {
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile(filename, fset.Base(), len(code))
	s.Init(file, code, nil, scanner.ScanComments)

	// 解析 #syntax=xx
	for {
		_, tok, lit := s.Scan()
		if tok != token.COMMENT {
			break
		}
		if strings.HasPrefix(lit, mimePrefix) {
			if mime := lit[len(mimePrefix):]; mime != "" {
				return mime
			}
			return ""
		}
	}

	// 根据文件名后缀解析
	if i := strings.LastIndex(filename, "."); i > 0 {
		if s := filename[i+1:]; s != "" {
			return s
		}
	}

	// 未知类型
	return ""
}
