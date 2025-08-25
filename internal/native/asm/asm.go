// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package asm

import (
	"bytes"
	"encoding/binary"
	"log"

	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/riscv"
)

// 汇编代码翻译位机器码
func AssemblerRV64(fnBody *ast.FuncBody) []byte {
	var buf bytes.Buffer
	for i, inst := range fnBody.Insts {
		// TODO: 处理长地址跳转
		x, err := riscv.EncodeRV64(inst.As, inst.Arg)
		if err != nil {
			log.Fatal(i, err)
		}
		err = binary.Write(&buf, binary.LittleEndian, x)
		if err != nil {
			log.Fatal(i, err)
		}
	}
	return buf.Bytes()
}
