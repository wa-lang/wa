// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wire

import (
	"wa-lang.org/wa/internal/logger"
)

/**************************************
本文件定义了 wire 中与值类型相关的各种对象
**************************************/

//-------------------------------------

/**************************************
Types : 类型库，用于管理值类型，类型的名字是其身份标志
类型不可以直接声明，必须通过 TypeLib.GenXXX 系列方法创建
对于大部分 TypeLib.GenXXX()：
  - 若指定名称的类型已存在于类型库中，则返回现存的类型
**************************************/
type Types struct {
	typs map[string]Type

	Void, Bool, U8, U16, U32, U64, Uint, I8, I16, I32, I64, Int, F32, F64, Complex64, Complex128, Rune, String Type
}

// Init 初始化类型库
func (tl *Types) Init() {
	tl.typs = make(map[string]Type)

	tl.Void = &Void{}
	tl.Bool = &Bool{}
	tl.U8 = &U8{}
	tl.U16 = &U16{}
	tl.U32 = &U32{}
	tl.U64 = &U64{}
	tl.Uint = &Uint{}
	tl.I8 = &I8{}
	tl.I16 = &I16{}
	tl.I32 = &I32{}
	tl.I64 = &I64{}
	tl.Int = &Int{}
	tl.F32 = &F32{}
	tl.F64 = &F64{}
	tl.Complex64 = &Complex64{}
	tl.Complex128 = &Complex128{}
	tl.Rune = &Rune{}
	tl.String = &String{}

	tl.typs[tl.Void.Name()] = tl.Void
	tl.typs[tl.Bool.Name()] = tl.Bool
	tl.typs[tl.U8.Name()] = tl.U8
	tl.typs[tl.U16.Name()] = tl.U16
	tl.typs[tl.U32.Name()] = tl.U32
	tl.typs[tl.U64.Name()] = tl.U64
	tl.typs[tl.Uint.Name()] = tl.Uint
	tl.typs[tl.I8.Name()] = tl.I8
	tl.typs[tl.I16.Name()] = tl.I16
	tl.typs[tl.I32.Name()] = tl.I32
	tl.typs[tl.I64.Name()] = tl.I64
	tl.typs[tl.Int.Name()] = tl.Int
	tl.typs[tl.F32.Name()] = tl.F32
	tl.typs[tl.F64.Name()] = tl.F64
	tl.typs[tl.Complex64.Name()] = tl.Complex64
	tl.typs[tl.Complex128.Name()] = tl.Complex128
	tl.typs[tl.Rune.Name()] = tl.Rune
	tl.typs[tl.String.Name()] = tl.String

}

// Lookup 根据给定的名字查找值类型
func (tl *Types) Lookup(name string) (t Type, ok bool) {
	t, ok = tl.typs[name]
	return
}

// add 向 TypeLib 中添加一个新类型，注意不可重复添加
func (tl *Types) add(t Type) {
	if _, ok := tl.Lookup(t.Name()); ok {
		logger.Fatalf("Type:%T already registered.", t)
	}
	tl.typs[t.Name()] = t
}

func (tl *Types) All() map[string]Type {
	return tl.typs
}

/**************************************
tCommon: 实现 Type 接口 Name、Method 相关方法
**************************************/
//type tCommon struct {
//	name string
//	//methods []Method
//}
//
//func (t *tCommon) Name() string { return t.name }
//func (t *tCommon) AddMethod(m Method) int {
//	t.methods = append(t.methods, m)
//	return len(t.methods) - 1
//}
//func (t *tCommon) NumMethods() int     { return len(t.methods) }
//func (t *tCommon) Method(i int) Method { return t.methods[i] }

/**************************************
Void: void，0字节
**************************************/
type Void struct{}

func (t *Void) Name() string      { return "void" }
func (t *Void) Kind() TypeKind    { return TypeKindVoid }
func (t *Void) Equal(u Type) bool { _, ok := u.(*Void); return ok }

/**************************************
Bool: 布尔，1字节
**************************************/
type Bool struct{}

func (t *Bool) Name() string      { return "bool" }
func (t *Bool) Kind() TypeKind    { return TypeKindBool }
func (t *Bool) Equal(u Type) bool { _, ok := u.(*Bool); return ok }

/**************************************
I8: 8位有符号整数，1字节
**************************************/
type I8 struct{}

func (t *I8) Name() string      { return "i8" }
func (t *I8) Kind() TypeKind    { return TypeKindI8 }
func (t *I8) Equal(u Type) bool { _, ok := u.(*I8); return ok }

/**************************************
U8: 8位无符号整数，1字节
**************************************/
type U8 struct{}

func (t *U8) Name() string      { return "u8" }
func (t *U8) Kind() TypeKind    { return TypeKindU8 }
func (t *U8) Equal(u Type) bool { _, ok := u.(*U8); return ok }

/**************************************
I16: 16位有符号整数，2字节
**************************************/
type I16 struct{}

func (t *I16) Name() string      { return "i16" }
func (t *I16) Kind() TypeKind    { return TypeKindI16 }
func (t *I16) Equal(u Type) bool { _, ok := u.(*I16); return ok }

/**************************************
tU16: 16位无符号整数，2字节
**************************************/
type U16 struct{}

func (t *U16) Name() string      { return "u16" }
func (t *U16) Kind() TypeKind    { return TypeKindU16 }
func (t *U16) Equal(u Type) bool { _, ok := u.(*U16); return ok }

/**************************************
I32: 32位有符号整数，4字节
**************************************/
type I32 struct{}

func (t *I32) Name() string      { return "i32" }
func (t *I32) Kind() TypeKind    { return TypeKindI32 }
func (t *I32) Equal(u Type) bool { _, ok := u.(*I32); return ok }

/**************************************
U32: 32位无符号整数，4字节
**************************************/
type U32 struct{}

func (t *U32) Name() string      { return "u32" }
func (t *U32) Kind() TypeKind    { return TypeKindU32 }
func (t *U32) Equal(u Type) bool { _, ok := u.(*U32); return ok }

/**************************************
I64: 64位有符号整数，8字节
**************************************/
type I64 struct{}

func (t *I64) Name() string      { return "i64" }
func (t *I64) Kind() TypeKind    { return TypeKindI64 }
func (t *I64) Equal(u Type) bool { _, ok := u.(*I64); return ok }

/**************************************
U64: 64位无符号整数，8字节
**************************************/
type U64 struct{}

func (t *U64) Name() string      { return "u64" }
func (t *U64) Kind() TypeKind    { return TypeKindU64 }
func (t *U64) Equal(u Type) bool { _, ok := u.(*U64); return ok }

/**************************************
Uint: 平台相关无符号整型
**************************************/
type Uint struct{}

func (t *Uint) Name() string      { return "uint" }
func (t *Uint) Kind() TypeKind    { return TypeKindUint }
func (t *Uint) Equal(u Type) bool { _, ok := u.(*Uint); return ok }

/**************************************
Int: 平台相关有符号整型
**************************************/
type Int struct{}

func (t *Int) Name() string      { return "int" }
func (t *Int) Kind() TypeKind    { return TypeKindInt }
func (t *Int) Equal(u Type) bool { _, ok := u.(*Int); return ok }

/**************************************
F32: 单精度浮点数，4字节
**************************************/
type F32 struct{}

func (t *F32) Name() string      { return "f32" }
func (t *F32) Kind() TypeKind    { return TypeKindF32 }
func (t *F32) Equal(u Type) bool { _, ok := u.(*F32); return ok }

/**************************************
F64: 双精度浮点数，8字节
**************************************/
type F64 struct{}

func (t *F64) Name() string      { return "f64" }
func (t *F64) Kind() TypeKind    { return TypeKindF64 }
func (t *F64) Equal(u Type) bool { _, ok := u.(*F64); return ok }

/**************************************
Complex64: 单精度复数，8字节
**************************************/
type Complex64 struct{}

func (t *Complex64) Name() string      { return "complex64" }
func (t *Complex64) Kind() TypeKind    { return TypeKindComplex64 }
func (t *Complex64) Equal(u Type) bool { _, ok := u.(*Complex64); return ok }

/**************************************
Complex128: 双精度复数，16字节
**************************************/
type Complex128 struct{}

func (t *Complex128) Name() string      { return "complex128" }
func (t *Complex128) Kind() TypeKind    { return TypeKindComplex128 }
func (t *Complex128) Equal(u Type) bool { _, ok := u.(*Complex128); return ok }

/**************************************
Rune: unicode字符，4字节
**************************************/
type Rune struct{}

func (t *Rune) Name() string      { return "rune" }
func (t *Rune) Kind() TypeKind    { return TypeKindRune }
func (t *Rune) Equal(u Type) bool { _, ok := u.(*Rune); return ok }

/**************************************
String: 字符串
**************************************/
type String struct{}

func (t *String) Name() string      { return "string" }
func (t *String) Kind() TypeKind    { return TypeKindString }
func (t *String) Equal(u Type) bool { _, ok := u.(*String); return ok }

/**************************************
Ptr: 指针，长度取决于目标平台
**************************************/
type Ptr struct {
	Base Type
}

func (t *Ptr) Name() string   { return t.Base.Name() + "$$ptr" }
func (t *Ptr) Kind() TypeKind { return TypeKindPtr }
func (t *Ptr) Equal(u Type) bool {
	if ut, ok := u.(*Ptr); ok {
		return t.Base.Equal(ut.Base)
	}
	return false
}

func (tl *Types) GenPtr(base Type) *Ptr {
	nt := &Ptr{Base: base}
	if t, ok := tl.Lookup(nt.Name()); ok {
		return t.(*Ptr)
	}

	tl.add(nt)
	return nt
}

/**************************************
Tuple: 元组
**************************************/
type Tuple struct {
	fields []Type
}

func (t *Tuple) Name() string   { panic("Todo") }
func (t *Tuple) Kind() TypeKind { return TypeKindTuple }
func (t *Tuple) Equal(u Type) bool {
	ut, ok := u.(*Tuple)
	if !ok {
		return false
	}

	if len(t.fields) != len(ut.fields) {
		return false
	}

	for i := range t.fields {
		if !t.fields[i].Equal(ut.fields[i]) {
			return false
		}
	}

	return true
}

func (tl *Types) GenTuple(fields []Type) *Tuple {
	panic("Todo") //name

	nt := &Tuple{fields: fields}
	tl.add(nt)
	return nt
}

/**************************************
Struct: 结构体
**************************************/
type Struct struct {
	fields []*StructField
}

func (t *Struct) Name() string   { panic("Todo") }
func (t *Struct) Kind() TypeKind { return TypeKindStruct }
func (t *Struct) Equal(u Type) bool {
	ut, ok := u.(*Struct)
	if !ok {
		return false
	}

	if len(t.fields) != len(ut.fields) {
		return false
	}

	for i := range t.fields {
		if !t.fields[i].Equal(ut.fields[i]) {
			return false
		}
	}

	return true
}

func (tl *Types) GenStruct() (*Struct, bool) {
	panic("Todo")
	//if t, ok := tl.Lookup(name); ok {
	//	return t.(*Struct), true
	//}
	//
	//var nt Struct
	//nt.name = name
	//tl.Add(&nt)
	//return &nt, false
}

func (t *Struct) AppendField(f *StructField) {
	f.id = len(t.fields)
	t.fields = append(t.fields, f)
}

type StructField struct {
	Name string
	Type Type
	id   int
}

func (i *StructField) Equal(u *StructField) bool { return i.Name == u.Name && i.Type.Equal(u.Type) }
