// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package xlang

import (
	"path/filepath"
	"strings"

	"wa-lang.org/wa/internal/scanner"
	"wa-lang.org/wa/internal/token"
)

// 判断代码的类型
func DetectLang(filename string, code []byte) token.LangType {
	// 线根据文件名判断
	switch strings.ToLower(filepath.Ext(filename)) {
	case ".wa":
		return token.LangType_Wa
	case ".wz":
		return token.LangType_Wz
	case ".wat":
		return token.LangType_Wat
	case ".ws":
		return token.LangType_Native
	}

	// 判断 wa/wz
	{
		var s scanner.Scanner
		fset := token.NewFileSet()
		file := fset.AddFile(filename, fset.Base(), len(code))
		s.Init(file, code, nil, scanner.ScanComments)

		for {
			_, tok, lit := s.Scan()
			if tok == token.EOF || tok == token.ILLEGAL {
				break
			}

			// 英文的关键字
			if tok.IsKeyword() {
				return token.LangType_Wa
			}

			// 中文的关键字
			if tok == token.IDENT {
				if token.LookupEx(lit, true) != token.IDENT {
					return token.LangType_Wz
				}
			}

			// 出现中文注释
			if tok.IsWzKeyword() || tok.IsWzComment(lit) {
				return token.LangType_Wz
			}

			// 跳过普通注释
			if tok == token.COMMENT {
				continue
			}

			// 如果不能识别的其他token, 就结束
			return token.LangType_Unknown
		}
	}

	// wat/ws 不做深度识别
	return token.LangType_Unknown
}
