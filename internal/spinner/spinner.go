// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package spinner

import (
	"strconv"

	"wa-lang.org/wa/internal/constant"
	"wa-lang.org/wa/internal/types"
	"wa-lang.org/wa/internal/wire"
)

/**************************************
本包用于将 AST 转为 wire
**************************************/

//-------------------------------------

/**************************************
IntSize: 目标平台整型数的宽度，不定宽整形、指针宽度均由其决定
**************************************/
type IntSize int

const (
	IntSize32 IntSize = iota
	IntSize64
)

/**************************************
Builder:
**************************************/
type Builder struct {
	IntSize IntSize
	info    *types.Info
	module  *wire.Module

	typeTable map[string]*wire.Type
}

// 初始化 Builder
func (b *Builder) Init(info *types.Info) {
	b.info = info
	b.module = &wire.Module{}
	b.module.Init()

	b.typeTable = make(map[string]*wire.Type)
	b.typeTable[b.module.Types.Void.Name()] = &b.module.Types.Void
	b.typeTable[b.module.Types.Bool.Name()] = &b.module.Types.Bool
	b.typeTable[b.module.Types.U8.Name()] = &b.module.Types.U8
	b.typeTable[b.module.Types.U16.Name()] = &b.module.Types.U16
	b.typeTable[b.module.Types.U32.Name()] = &b.module.Types.U32
	b.typeTable[b.module.Types.U64.Name()] = &b.module.Types.U64
	b.typeTable[b.module.Types.Uint.Name()] = &b.module.Types.Uint
	b.typeTable[b.module.Types.I8.Name()] = &b.module.Types.I8
	b.typeTable[b.module.Types.I16.Name()] = &b.module.Types.I16
	b.typeTable[b.module.Types.I32.Name()] = &b.module.Types.I32
	b.typeTable[b.module.Types.I64.Name()] = &b.module.Types.I64
	b.typeTable[b.module.Types.Int.Name()] = &b.module.Types.Int
	b.typeTable[b.module.Types.F32.Name()] = &b.module.Types.F32
	b.typeTable[b.module.Types.F64.Name()] = &b.module.Types.F64
	b.typeTable[b.module.Types.Complex64.Name()] = &b.module.Types.Complex64
	b.typeTable[b.module.Types.Complex128.Name()] = &b.module.Types.Complex128
	b.typeTable[b.module.Types.Rune.Name()] = &b.module.Types.Rune
	b.typeTable[b.module.Types.String.Name()] = &b.module.Types.String
}

func (b *Builder) constval(val constant.Value, typ types.Type, pos int) wire.Value {
	t := b.BuildType(typ)

	if t.Equal(b.module.Types.Bool) {
		if constant.BoolVal(val) {
			return b.module.NewConst("1", t, pos)
		} else {
			return b.module.NewConst("0", t, pos)
		}
	} else if t.Equal(b.module.Types.U8) || t.Equal(b.module.Types.U16) || t.Equal(b.module.Types.U32) || t.Equal(b.module.Types.U64) || t.Equal(b.module.Types.Uint) {
		val, _ := constant.Uint64Val(val)
		return b.module.NewConst(strconv.Itoa(int(val)), t, pos)
	} else if t.Equal(b.module.Types.F32) {
		val, _ := constant.Float64Val(val)
		return b.module.NewConst(strconv.FormatFloat(val, 'f', -1, 32), t, pos)
	} else if t.Equal(b.module.Types.F64) {
		val, _ := constant.Float64Val(val)
		return b.module.NewConst(strconv.FormatFloat(val, 'f', -1, 64), t, pos)
	} else if t.Equal(b.module.Types.Complex64) || t.Equal(b.module.Types.Complex128) {
		re, _ := constant.Float64Val(constant.Real(val))
		im, _ := constant.Float64Val(constant.Imag(val))
		res := strconv.FormatFloat(re, 'f', -1, 64)
		ims := strconv.FormatFloat(im, 'f', -1, 64)
		return b.module.NewConst(res+" "+ims, t, pos)
	} else if t.Equal(b.module.Types.I8) || t.Equal(b.module.Types.I16) || t.Equal(b.module.Types.I32) || t.Equal(b.module.Types.I64) || t.Equal(b.module.Types.Int) || t.Equal(b.module.Types.Rune) {
		val, _ := constant.Int64Val(val)
		return b.module.NewConst(strconv.Itoa(int(val)), t, pos)
	}

	panic("Todo")
}

func (b *Builder) nilConst(typ types.Type, pos int) wire.Value {
	t := b.BuildType(typ)
	return b.module.NewConst("0", t, pos)
}
