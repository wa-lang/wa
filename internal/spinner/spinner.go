// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package spinner

import (
	"fmt"
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
	IntSize   IntSize
	info      *types.Info
	module    *wire.Module
	typeTable map[string]*wire.ValueType

	Void, Bool, U8, U16, U32, U64, I8, I16, I32, I64, F32, F64, Complex64, Complex128, Rune, String, Uint, Int wire.ValueType
}

// 初始化 Builder
func (b *Builder) Init(info *types.Info) {
	b.info = info

	b.module = &wire.Module{}
	b.module.Init()

	b.Void = b.module.Types.GenVoid("void")
	b.Bool = b.module.Types.GenBool("bool")
	b.U8 = b.module.Types.GenU8("u8")
	b.U16 = b.module.Types.GenU16("u16")
	b.U32 = b.module.Types.GenU32("u32")
	b.U64 = b.module.Types.GenU64("u64")
	b.I8 = b.module.Types.GenI8("i8")
	b.I16 = b.module.Types.GenI16("i16")
	b.I32 = b.module.Types.GenI32("i32")
	b.I64 = b.module.Types.GenI64("i64")
	b.F32 = b.module.Types.GenF32("f32")
	b.F64 = b.module.Types.GenF64("f64")
	b.Complex64 = b.module.Types.GenComplex64("complex64")
	b.Complex128 = b.module.Types.GenComplex128("complex128")
	b.Rune = b.module.Types.GenRune("rune")
	b.String = b.module.Types.GenString("string")
	switch b.IntSize {
	case IntSize32:
		b.Int = b.I32
		b.Uint = b.U32
	case IntSize64:
		b.Int = b.I64
		b.Uint = b.U64
	default:
		panic(fmt.Sprintf("Unknown TargetWidth: %v", b.IntSize))
	}

	b.typeTable = make(map[string]*wire.ValueType)
	b.typeTable["bool"] = &b.Bool
	b.typeTable["u8"] = &b.U8
	b.typeTable["u16"] = &b.U16
	b.typeTable["u32"] = &b.U32
	b.typeTable["u64"] = &b.U64
	b.typeTable["i8"] = &b.I8
	b.typeTable["i16"] = &b.I16
	b.typeTable["i32"] = &b.I32
	b.typeTable["i64"] = &b.I64
	b.typeTable["f32"] = &b.F32
	b.typeTable["f64"] = &b.F64
	b.typeTable["complex64"] = &b.Complex64
	b.typeTable["complex128"] = &b.Complex128
	b.typeTable["rune"] = &b.Rune
	b.typeTable["string"] = &b.String
	b.typeTable["int"] = &b.Int
	b.typeTable["uint"] = &b.Uint
}

func (b *Builder) constval(val constant.Value, typ types.Type, pos int) wire.Value {
	t := b.BuildType(typ)

	if t.Equal(b.Bool) {
		if constant.BoolVal(val) {
			return b.module.NewConst("1", t, pos)
		} else {
			return b.module.NewConst("0", t, pos)
		}
	} else if t.Equal(b.U8) || t.Equal(b.U16) || t.Equal(b.U32) || t.Equal(b.U64) || t.Equal(b.Uint) {
		val, _ := constant.Uint64Val(val)
		return b.module.NewConst(strconv.Itoa(int(val)), t, pos)
	} else if t.Equal(b.F32) {
		val, _ := constant.Float64Val(val)
		return b.module.NewConst(strconv.FormatFloat(val, 'f', -1, 32), t, pos)
	} else if t.Equal(b.F64) {
		val, _ := constant.Float64Val(val)
		return b.module.NewConst(strconv.FormatFloat(val, 'f', -1, 64), t, pos)
	} else if t.Equal(b.Complex64) || t.Equal(b.Complex128) {
		re, _ := constant.Float64Val(constant.Real(val))
		im, _ := constant.Float64Val(constant.Imag(val))
		res := strconv.FormatFloat(re, 'f', -1, 64)
		ims := strconv.FormatFloat(im, 'f', -1, 64)
		return b.module.NewConst(res+" "+ims, t, pos)
	} else if t.Equal(b.I8) || t.Equal(b.I16) || t.Equal(b.I32) || t.Equal(b.I64) || t.Equal(b.Int) || t.Equal(b.Rune) {
		val, _ := constant.Int64Val(val)
		return b.module.NewConst(strconv.Itoa(int(val)), t, pos)
	}

	panic("Todo")
}

func (b *Builder) nilConst(typ types.Type, pos int) wire.Value {
	t := b.BuildType(typ)
	return b.module.NewConst("0", t, pos)
}
