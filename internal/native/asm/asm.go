// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package asm

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/riscv"
	"wa-lang.org/wa/internal/native/token"
)

// 汇编代码翻译位机器码
func AssemblerRV64(fnBody *ast.FuncBody) []byte {
	var buf bytes.Buffer
	for i, inst := range fnBody.Insts {
		// TODO: 处理长地址跳转
		x, err := riscv.EncodeRV64(inst.As, inst.Arg)
		if err != nil {
			panic(fmt.Errorf("%d: %w", i, err))
		}
		err = binary.Write(&buf, binary.LittleEndian, x)
		if err != nil {
			panic(fmt.Errorf("%d: %w", i, err))
		}
	}
	return buf.Bytes()
}

// 将汇编语法树转为固定位置的机器码
func AssembleFile(fset *token.FileSet, file *ast.File, opt *abi.LinkOptions) (prog *abi.LinkedProgram, err error) {
	prog = &abi.LinkedProgram{
		CPU:      file.CPU,
		TextAddr: int64(opt.Tdata),
		DataAddr: int64(opt.Tdata),
	}

	// 编译全局变量
	var bufData bytes.Buffer
	for _, g := range file.Globals {
		_ = g // todo
	}

	// 编译函数
	var bufText bytes.Buffer
	for _, fn := range file.Funcs {
		for _, inst := range fn.Body.Insts {
			if inst.Label != "" {
				panic("TODO: gen Label address")
			}
			if inst.As == 0 {
				continue
			}

			// TODO: 处理长地址跳转
			x, err := riscv.EncodeRV64(inst.As, inst.Arg)
			if err != nil {
				panic(fmt.Errorf("%v: %w", inst, err))
			}
			err = binary.Write(&bufText, binary.LittleEndian, x)
			if err != nil {
				panic(fmt.Errorf("%v: %w", inst, err))
			}
		}
	}

	prog.TextData = bufText.Bytes()
	prog.DataData = bufData.Bytes()

	return
}
