// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package arch

// 汇编器和链接器的参数
type Flags struct {
	Debug      bool   // 调试模式 dump instructions as they are parsed
	OutputFile string // 输出文件 output file; default foo.6 for /a/b/c/foo.s on amd64
	PrintOut   bool   // 打印汇编和机器码 print assembly and machine code
	Shared     bool   // 生成动态库 generate code that can be linked into a shared library
	Dynlink    bool   // 生成可重定位的符号 support references to symbols defined in other shared libraries

	Defines     []string // 预定义的值 predefined symbol with optional simple value -D=identifer=value; can be set multiple times
	IncludeDirs []string // 头文件路径 include directory; can be set multiple times
}
