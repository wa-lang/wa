// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package xlang

import (
	"path/filepath"
	"strings"

	nav_scanner "wa-lang.org/wa/internal/native/scanner"
	nav_token "wa-lang.org/wa/internal/native/token"
	"wa-lang.org/wa/internal/scanner"
	"wa-lang.org/wa/internal/token"
	wat_scanner "wa-lang.org/wa/internal/wat/scanner"
	wat_token "wa-lang.org/wa/internal/wat/token"
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
	}

	// 后缀名有多个 '.'
	switch s := strings.ToLower(filename); true {
	case strings.HasSuffix(s, ".wa.s"):
		return token.LangType_Nasm
	case strings.HasSuffix(s, ".wz.s"):
		return token.LangType_Nasm
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

	// 判断 nasm-gas/nasm-zh
	{
		var s nav_scanner.Scanner
		fset := nav_token.NewFileSet()
		file := fset.AddFile(filename, fset.Base(), len(code))
		s.Init(file, code, nil, nav_scanner.ScanComments)

		for {
			_, tok, _ := s.Scan()
			if tok == nav_token.EOF || tok == nav_token.ILLEGAL {
				break
			}

			// 识别关键字
			if tok.IsGasKeyword() || tok.IsZhKeyword() {
				return token.LangType_Nasm
			}

			// 跳过普通注释
			if tok == nav_token.COMMENT {
				continue
			}
		}
	}

	// 识别 wat 类型
	{
		var s wat_scanner.Scanner
		file := wat_token.NewFile(filename, len(code))
		s.Init(file, code, nil, wat_scanner.ScanComments)

		for {
			_, tok, _ := s.Scan()
			if tok == wat_token.EOF || tok == wat_token.ILLEGAL {
				break
			}

			if tok.IsKeyword() {
				return token.LangType_Wat
			}

			// 跳过普通注释
			if tok == wat_token.COMMENT {
				continue
			}
		}
	}

	// 未知类型
	return token.LangType_Unknown
}
