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
	typs    map[string]ValueType
	ptrSize int
}

// Init 初始化类型库
func (tl *Types) Init() {
	tl.typs = make(map[string]ValueType)
	tl.ptrSize = 4 //Todo 平台相关
}

// Lookup 根据给定的名字查找值类型
func (tl *Types) Lookup(name string) (t ValueType, ok bool) {
	t, ok = tl.typs[name]
	return
}

// Add 向 TypeLib 中添加一个新类型，注意不可重复添加
func (tl *Types) Add(t ValueType) {
	if _, ok := tl.Lookup(t.Name()); ok {
		logger.Fatalf("ValueType:%T already registered.", t)
	}
	tl.typs[t.Name()] = t
}

/**************************************
tCommon: 实现 ValueType 接口 Name、Method 相关方法
**************************************/
type tCommon struct {
	name    string
	methods []Method
}

func (t *tCommon) Name() string { return t.name }
func (t *tCommon) AddMethod(m Method) int {
	t.methods = append(t.methods, m)
	return len(t.methods) - 1
}
func (t *tCommon) NumMethods() int     { return len(t.methods) }
func (t *tCommon) Method(i int) Method { return t.methods[i] }

/**************************************
Void: void，0字节
**************************************/
type Void struct {
	tCommon
}

func (t *Void) Kind() TypeKind         { return TypeKindVoid }
func (t *Void) Equal(u ValueType) bool { _, ok := u.(*Void); return ok }

func (tl *Types) GenVoid(name string) *Void {
	if t, ok := tl.Lookup(name); ok {
		return t.(*Void)
	}

	nt := Void{}
	nt.name = name
	tl.Add(&nt)
	return &nt
}

/**************************************
Bool: 布尔，1字节
**************************************/
type Bool struct {
	tCommon
}

func (t *Bool) Kind() TypeKind         { return TypeKindBool }
func (t *Bool) Equal(u ValueType) bool { _, ok := u.(*Bool); return ok }

func (tl *Types) GenBool(name string) *Bool {
	if t, ok := tl.Lookup(name); ok {
		return t.(*Bool)
	}

	nt := Bool{}
	nt.name = name
	tl.Add(&nt)
	return &nt
}

/**************************************
I8: 8位有符号整数，1字节
**************************************/
type I8 struct {
	tCommon
}

func (t *I8) Kind() TypeKind         { return TypeKindI8 }
func (t *I8) Equal(u ValueType) bool { _, ok := u.(*I8); return ok }

func (tl *Types) GenI8(name string) *I8 {
	if t, ok := tl.Lookup(name); ok {
		return t.(*I8)
	}

	nt := I8{}
	nt.name = name
	tl.Add(&nt)
	return &nt
}

/**************************************
U8: 8位无符号整数，1字节
**************************************/
type U8 struct {
	tCommon
}

func (t *U8) Kind() TypeKind         { return TypeKindU8 }
func (t *U8) Equal(u ValueType) bool { _, ok := u.(*U8); return ok }

func (tl *Types) GenU8(name string) *U8 {
	if t, ok := tl.Lookup(name); ok {
		return t.(*U8)
	}

	nt := U8{}
	nt.name = name
	tl.Add(&nt)
	return &nt
}

/**************************************
I16: 16位有符号整数，2字节
**************************************/
type I16 struct {
	tCommon
}

func (t *I16) Kind() TypeKind         { return TypeKindI16 }
func (t *I16) Equal(u ValueType) bool { _, ok := u.(*I16); return ok }

func (tl *Types) GenI16(name string) *I16 {
	if t, ok := tl.Lookup(name); ok {
		return t.(*I16)
	}

	nt := I16{}
	nt.name = name
	tl.Add(&nt)
	return &nt
}

/**************************************
tU16: 16位无符号整数，2字节
**************************************/
type U16 struct {
	tCommon
}

func (t *U16) Kind() TypeKind         { return TypeKindU16 }
func (t *U16) Equal(u ValueType) bool { _, ok := u.(*U16); return ok }

func (tl *Types) GenU16(name string) *U16 {
	if t, ok := tl.Lookup(name); ok {
		return t.(*U16)
	}

	nt := U16{}
	nt.name = name
	tl.Add(&nt)
	return &nt
}

/**************************************
I32: 32位有符号整数，4字节
**************************************/
type I32 struct {
	tCommon
}

func (t *I32) Kind() TypeKind         { return TypeKindI32 }
func (t *I32) Equal(u ValueType) bool { _, ok := u.(*I32); return ok }

func (tl *Types) GenI32(name string) *I32 {
	if t, ok := tl.Lookup(name); ok {
		return t.(*I32)
	}

	nt := I32{}
	nt.name = name
	tl.Add(&nt)
	return &nt
}

/**************************************
U32: 32位无符号整数，4字节
**************************************/
type U32 struct {
	tCommon
}

func (t *U32) Kind() TypeKind         { return TypeKindU32 }
func (t *U32) Equal(u ValueType) bool { _, ok := u.(*U32); return ok }

func (tl *Types) GenU32(name string) *U32 {
	if t, ok := tl.Lookup(name); ok {
		return t.(*U32)
	}

	nt := U32{}
	nt.name = name
	tl.Add(&nt)
	return &nt
}

/**************************************
I64: 64位有符号整数，8字节
**************************************/
type I64 struct {
	tCommon
}

func (t *I64) Kind() TypeKind         { return TypeKindI64 }
func (t *I64) Equal(u ValueType) bool { _, ok := u.(*I64); return ok }

func (tl *Types) GenI64(name string) *I64 {
	if t, ok := tl.Lookup(name); ok {
		return t.(*I64)
	}

	nt := I64{}
	nt.name = name
	tl.Add(&nt)
	return &nt
}

/**************************************
U64: 64位无符号整数，8字节
**************************************/
type U64 struct {
	tCommon
}

func (t *U64) Kind() TypeKind         { return TypeKindU64 }
func (t *U64) Equal(u ValueType) bool { _, ok := u.(*U64); return ok }

func (tl *Types) GenU64(name string) *U64 {
	if t, ok := tl.Lookup(name); ok {
		return t.(*U64)
	}

	nt := U64{}
	nt.name = name
	tl.Add(&nt)
	return &nt
}

/**************************************
F32: 单精度浮点数，4字节
**************************************/
type F32 struct {
	tCommon
}

func (t *F32) Kind() TypeKind         { return TypeKindF32 }
func (t *F32) Equal(u ValueType) bool { _, ok := u.(*F32); return ok }

func (tl *Types) GenF32(name string) *F32 {
	if t, ok := tl.Lookup(name); ok {
		return t.(*F32)
	}

	nt := F32{}
	nt.name = name
	tl.Add(&nt)
	return &nt
}

/**************************************
F64: 双精度浮点数，8字节
**************************************/
type F64 struct {
	tCommon
}

func (t *F64) Kind() TypeKind         { return TypeKindF64 }
func (t *F64) Equal(u ValueType) bool { _, ok := u.(*F64); return ok }

func (tl *Types) GenF64(name string) *F64 {
	if t, ok := tl.Lookup(name); ok {
		return t.(*F64)
	}

	nt := F64{}
	nt.name = name
	tl.Add(&nt)
	return &nt
}

/**************************************
Complex64: 单精度复数，8字节
**************************************/
type Complex64 struct {
	tCommon
}

func (t *Complex64) Kind() TypeKind         { return TypeKindComplex64 }
func (t *Complex64) Equal(u ValueType) bool { _, ok := u.(*Complex64); return ok }

func (tl *Types) GenComplex64(name string) *Complex64 {
	if t, ok := tl.Lookup(name); ok {
		return t.(*Complex64)
	}

	nt := Complex64{}
	nt.name = name
	tl.Add(&nt)
	return &nt
}

/**************************************
Complex128: 双精度复数，16字节
**************************************/
type Complex128 struct {
	tCommon
}

func (t *Complex128) Kind() TypeKind         { return TypeKindComplex128 }
func (t *Complex128) Equal(u ValueType) bool { _, ok := u.(*Complex128); return ok }

func (tl *Types) GenComplex128(name string) *Complex128 {
	if t, ok := tl.Lookup(name); ok {
		return t.(*Complex128)
	}

	nt := Complex128{}
	nt.name = name
	tl.Add(&nt)
	return &nt
}

/**************************************
Rune: unicode字符，4字节
**************************************/
type Rune struct {
	tCommon
}

func (t *Rune) Kind() TypeKind         { return TypeKindRune }
func (t *Rune) Equal(u ValueType) bool { _, ok := u.(*Rune); return ok }

func (tl *Types) GenRune(name string) *Rune {
	if t, ok := tl.Lookup(name); ok {
		return t.(*Rune)
	}

	nt := Rune{}
	nt.name = name
	tl.Add(&nt)
	return &nt
}

/**************************************
String: 字符串
**************************************/
type String struct {
	tCommon
}

func (t *String) Kind() TypeKind         { return TypeKindString }
func (t *String) Equal(u ValueType) bool { _, ok := u.(*String); return ok }

func (tl *Types) GenString(name string) *String {
	if t, ok := tl.Lookup(name); ok {
		return t.(*String)
	}

	nt := String{}
	nt.name = name
	tl.Add(&nt)
	return &nt
}

/**************************************
Ptr: 指针，长度取决于目标平台
**************************************/
type Ptr struct {
	tCommon
	Base ValueType
}

func (t *Ptr) Kind() TypeKind { return TypeKindPtr }
func (t *Ptr) Equal(u ValueType) bool {
	if ut, ok := u.(*Ptr); ok {
		return t.Base.Equal(ut.Base)
	}
	return false
}

func (tl *Types) GenPtr(base ValueType) *Ptr {
	name := base.Name() + "$$ptr"
	if t, ok := tl.Lookup(name); ok {
		return t.(*Ptr)
	}

	nt := Ptr{Base: base}
	nt.name = name
	tl.Add(&nt)
	return &nt
}

/**************************************
Struct: 结构体
**************************************/
type Struct struct {
	tCommon
	fields []*StructField
}

func (t *Struct) Kind() TypeKind { return TypeKindStruct }
func (t *Struct) Equal(u ValueType) bool {
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

func (tl *Types) GenStruct(name string) (*Struct, bool) {
	if t, ok := tl.Lookup(name); ok {
		return t.(*Struct), true
	}

	var nt Struct
	nt.name = name
	tl.Add(&nt)
	return &nt, false
}

func (t *Struct) AppendField(f *StructField) {
	f.id = len(t.fields)
	t.fields = append(t.fields, f)
}

type StructField struct {
	Name string
	Type ValueType
	id   int
}

func (i *StructField) Equal(u *StructField) bool { return i.Name == u.Name && i.Type.Equal(u.Type) }
