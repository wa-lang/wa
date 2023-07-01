// 版权 @2023 凹语言 作者。保留所有权利。

package wamime

import (
	"strings"

	"wa-lang.org/wa/internal/scanner"
	"wa-lang.org/wa/internal/token"
)

const mimePrefix0 = "#syntax="   // 旧语法, 先保持兼容
const mimePrefix = "#wa:syntax=" // #wa:syntax=wa, #wa:syntax=wz

func GetCodeMime(filename string, code []byte) string {
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile(filename, fset.Base(), len(code))
	s.Init(file, code, nil, scanner.ScanComments)

	// 解析 #wa:syntax=xx
	for {
		_, tok, lit := s.Scan()
		if tok != token.COMMENT {
			break
		}
		// 旧语法 #syntax=
		if strings.HasPrefix(lit, mimePrefix0) {
			if mime := lit[len(mimePrefix0):]; mime != "" {
				return mime
			}
			return ""
		}
		// 新语法 #wa:syntax=
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
