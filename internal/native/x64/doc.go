// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package x64

// 指令有3种等价的表示:
// 1. `native/parser` 汇编解析器从文本解析得到的 `ast.Instructiton`
// 2. `native/x64/x86asm` 从二进制汇编指令解析得到的 `x86asm.Inst`
// 3. `native/x64/p9x86` 对应的 `p9x86.Prog` 定义的, 用于编码到机器指令
//
// 几种指令的使用场景:
// 1. `wa objdump` 从二进制解码并以 intel 语法显示, 使用 x86asm.Inst
// 2. `wa asm2elf` 解析汇编代码得到 `ast.Instructiton`, 然后通过 `x64.BuildProg()` 转化为 `p9x86.Prog` 后进行指令编码
//
// 补充数明:
// - `ast.Instructiton` 支持格式化打印
// - `x86asm.Inst` 支持格式化打印
// - `p9x86.Prog` 不支持格式化打印
