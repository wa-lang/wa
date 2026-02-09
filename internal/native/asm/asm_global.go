// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package asm

import (
	"encoding/binary"
	"fmt"
	"math"

	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/token"
)

func (p *_Assembler) asmGlobal(g *ast.Global) (err error) {
	// g.LinkInfo.Data 空间需要提前初始化
	if g.Init.Symbal != "" {
		v, ok := p.symbolAddress(g.Init.Symbal)
		if !ok {
			panic(fmt.Errorf("symbol %q not found", g.Init.Symbal))
		}
		if p.opt.CPU == abi.RISCV32 {
			binary.LittleEndian.PutUint32(g.LinkInfo.Data, uint32(v))
		} else {
			binary.LittleEndian.PutUint64(g.LinkInfo.Data, uint64(v))
		}
		return nil
	}

	// 常量面值初始化
	switch g.Type {
	case token.Bin:
		copy(g.LinkInfo.Data, []byte(g.Init.Lit.ConstV.(string)))
	case token.I8:
		v := g.Init.Lit.ConstV.(int64)
		switch {
		case v >= 0 && v <= math.MaxUint8:
			g.LinkInfo.Data = append(g.LinkInfo.Data, uint8(v))
		case v >= math.MinInt8 && v <= math.MaxInt8:
			g.LinkInfo.Data = append(g.LinkInfo.Data, uint8(int8(v)))
		default:
			panic(fmt.Errorf("global %q init value overflow: %v", g.Name, g.Init.Lit))
		}
	case token.I16:
		v := g.Init.Lit.ConstV.(int64)
		switch {
		case v >= 0 && v <= math.MaxUint16:
			binary.LittleEndian.PutUint16(g.LinkInfo.Data, uint16(v))
		case v >= math.MinInt16 && v <= math.MaxInt16:
			binary.LittleEndian.PutUint16(g.LinkInfo.Data, uint16(int16(v)))
		default:
			panic(fmt.Errorf("global %q init value overflow: %v", g.Name, g.Init.Lit))
		}
	case token.I32:
		v := g.Init.Lit.ConstV.(int64)
		switch {
		case v >= 0 && v <= math.MaxUint32:
			binary.LittleEndian.PutUint32(g.LinkInfo.Data, uint32(v))
		case v >= math.MinInt32 && v <= math.MaxInt32:
			binary.LittleEndian.PutUint32(g.LinkInfo.Data, uint32(int32(v)))
		default:
			panic(fmt.Errorf("global %q init value overflow: %v", g.Name, g.Init.Lit))
		}
	case token.I64:
		v := g.Init.Lit.ConstV.(int64)
		binary.LittleEndian.PutUint64(g.LinkInfo.Data, uint64(v))
	case token.F32:
		v := g.Init.Lit.ConstV.(float64)
		binary.LittleEndian.PutUint32(g.LinkInfo.Data, math.Float32bits(float32(v)))
	case token.F64:
		v := g.Init.Lit.ConstV.(float64)
		binary.LittleEndian.PutUint64(g.LinkInfo.Data, math.Float64bits(float64(v)))
	default:
		panic("unreahable")
	}

	return nil
}
