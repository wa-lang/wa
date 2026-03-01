// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package abi

// 文件后缀名类型
// 测试文件名的前缀和后缀单独出理
type FileExtType int

const (
	FileExtType_Nil       FileExtType = iota
	FileExtType_gcc_txt               // *.gcc.txt, gcc 的编译+链接参数
	FileExtType_import_js             // *.import.js
	FileExtType_wa                    // *.wa
	FileExtType_wa_c                  // *.wa.c
	FileExtType_wa_s                  // *.wa.s
	FileExtType_wz                    // *.wz
	FileExtType_wz_c                  // *.wz.c
	FileExtType_wz_s                  // *.wz.s
	FileExtType_wat                   // *.wat
)

func ParseFileExtType(s string) FileExtType {
	switch {
	case strHasSuffix(s, ".gcc.txt"):
		return FileExtType_gcc_txt
	case strHasSuffix(s, ".import.js"):
		return FileExtType_import_js
	case strHasSuffix(s, ".wa"):
		return FileExtType_wa
	case strHasSuffix(s, ".wa.c"):
		return FileExtType_wa_c
	case strHasSuffix(s, ".wa.s"):
		return FileExtType_wa_s
	case strHasSuffix(s, ".wz"):
		return FileExtType_wz
	case strHasSuffix(s, ".wz.c"):
		return FileExtType_wz_c
	case strHasSuffix(s, ".wz.s"):
		return FileExtType_wz_s
	case strHasSuffix(s, ".wat"):
		return FileExtType_wat
	}
	return FileExtType_Nil
}

func FileExtTypeMatched(path string, target TargetType, ext FileExtType) bool {
	sepIdx := strLastIndexByte(path, '_')
	if sepIdx <= 0 {
		return false
	}

	suffixExpect := TargetTypeString(target) + FileExtTypeString(ext)
	suffixGot := path[sepIdx+1:]

	return suffixGot == suffixExpect
}

func FileExtTypeString(t FileExtType) string {
	switch t {
	case FileExtType_gcc_txt: // gcc 的编译+链接参数
		return ".gcc.txt"
	case FileExtType_import_js:
		return ".import.js"
	case FileExtType_wa: // *.wa
		return ".wa"
	case FileExtType_wa_c: // *.wa.c
		return ".wa.c"
	case FileExtType_wa_s: // *.wa.s
		return ".wa.s"
	case FileExtType_wz: // *.wz
		return ".wz"
	case FileExtType_wz_c: // *.wz.c
		return ".wz.c"
	case FileExtType_wz_s: // *.wz.s
		return ".wz.s"
	case FileExtType_wat: // *.wat
		return ".wat"
	}
	return "abi.FileExtType(" + int2str(int(t)) + ")"
}
