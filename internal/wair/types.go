// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wair

import (
	"wa-lang.org/wa/internal/logger"
)

/**************************************
本文件定义了 wair 中与值类型相关的各种对象
**************************************/

//-------------------------------------

/**************************************
Types : 类型库，用于管理值类型，类型的名字是其身份标志
类型不可以直接声明，必须通过 TypeLib.GenXXX 系列方法创建
对于大部分 TypeLib.GenXXX()：
  - name 参数可为 nil 或 ""，此时将使用默认名称
  - 若指定名称的类型已存在于类型库中，则返回现存的类型
**************************************/
type Types struct {
	typs map[string]ValueType
}

// Init 初始化类型库
func (t *Types) Init() {
	t.typs = make(map[string]ValueType)
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
//tCommon: 实现 ValueType 接口 Name、Method 相关方法
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
Void: void 类型
**************************************/
type Void struct {
	tCommon
}

func (tl *Types) GenVoid(name string) *Void {
	nt := Void{}
	if len(name) > 0 {
		nt.name = name
	} else {
		nt.name = "void"
	}
	t, ok := tl.Lookup(nt.name)
	if ok {
		return t.(*Void)
	}

	tl.Add(&nt)
	return &nt
}
func (t *Void) Size() int              { return 0 }
func (t *Void) Align() int             { return 0 }
func (t *Void) Kind() TypeKind         { return TypeKindVoid }
func (t *Void) Equal(u ValueType) bool { _, ok := u.(*Void); return ok }

/**************************************
Bool:
**************************************/
type Bool struct {
	tCommon
}

func (tl *Types) GenBool(name string) *Bool {
	nt := Bool{}
	if len(name) > 0 {
		nt.name = name
	} else {
		nt.name = "bool"
	}
	t, ok := tl.Lookup(nt.name)
	if ok {
		return t.(*Bool)
	}

	tl.Add(&nt)
	return &nt
}
func (t *Bool) Size() int              { return 1 }
func (t *Bool) Align() int             { return 1 }
func (t *Bool) Kind() TypeKind         { return TypeKindBool }
func (t *Bool) Equal(u ValueType) bool { _, ok := u.(*Bool); return ok }

/**************************************
tRune:
**************************************/
type Rune struct {
	tCommon
}

func (tl *Types) GenRune(name string) *Rune {
	nt := Rune{}
	if len(name) > 0 {
		nt.name = name
	} else {
		nt.name = "rune"
	}
	t, ok := tl.Lookup(nt.name)
	if ok {
		return t.(*Rune)
	}

	tl.Add(&nt)
	return &nt
}
func (t *Rune) Size() int              { return 4 }
func (t *Rune) Align() int             { return 4 }
func (t *Rune) Kind() TypeKind         { return TypeKindRune }
func (t *Rune) Equal(u ValueType) bool { _, ok := u.(*Rune); return ok }

/**************************************
I8:
**************************************/
type I8 struct {
	tCommon
}

func (tl *Types) GenI8(name string) *I8 {
	nt := I8{}
	if len(name) > 0 {
		nt.name = name
	} else {
		nt.name = "i8"
	}
	t, ok := tl.Lookup(nt.name)
	if ok {
		return t.(*I8)
	}

	tl.Add(&nt)
	return &nt
}
func (t *I8) Size() int              { return 1 }
func (t *I8) Align() int             { return 1 }
func (t *I8) Kind() TypeKind         { return TypeKindI8 }
func (t *I8) Equal(u ValueType) bool { _, ok := u.(*I8); return ok }

/**************************************
U8:
**************************************/
type U8 struct {
	tCommon
}

func (tl *Types) GenU8(name string) *U8 {
	nt := U8{}
	if len(name) > 0 {
		nt.name = name
	} else {
		nt.name = "u8"
	}
	t, ok := tl.Lookup(nt.name)
	if ok {
		return t.(*U8)
	}

	tl.Add(&nt)
	return &nt
}
func (t *U8) Size() int              { return 1 }
func (t *U8) Align() int             { return 1 }
func (t *U8) Kind() TypeKind         { return TypeKindU8 }
func (t *U8) Equal(u ValueType) bool { _, ok := u.(*U8); return ok }

/**************************************
I16:
**************************************/
type I16 struct {
	tCommon
}

func (tl *Types) GenI16(name string) *I16 {
	nt := I16{}
	if len(name) > 0 {
		nt.name = name
	} else {
		nt.name = "i16"
	}
	t, ok := tl.Lookup(nt.name)
	if ok {
		return t.(*I16)
	}

	tl.Add(&nt)
	return &nt
}
func (t *I16) Size() int              { return 2 }
func (t *I16) Align() int             { return 2 }
func (t *I16) Kind() TypeKind         { return TypeKindI16 }
func (t *I16) Equal(u ValueType) bool { _, ok := u.(*I16); return ok }

/**************************************
tU16:
**************************************/
type U16 struct {
	tCommon
}

func (tl *Types) GenU16(name string) *U16 {
	nt := U16{}
	if len(name) > 0 {
		nt.name = name
	} else {
		nt.name = "u16"
	}
	t, ok := tl.Lookup(nt.name)
	if ok {
		return t.(*U16)
	}

	tl.Add(&nt)
	return &nt
}
func (t *U16) Size() int              { return 2 }
func (t *U16) Align() int             { return 2 }
func (t *U16) Kind() TypeKind         { return TypeKindU16 }
func (t *U16) Equal(u ValueType) bool { _, ok := u.(*U16); return ok }

/**************************************
I32:
**************************************/
type I32 struct {
	tCommon
}

func (tl *Types) GenValueI32(name string) *I32 {
	nt := I32{}
	if len(name) > 0 {
		nt.name = name
	} else {
		nt.name = "i32"
	}
	t, ok := tl.Lookup(nt.name)
	if ok {
		return t.(*I32)
	}

	tl.Add(&nt)
	return &nt
}
func (t *I32) Size() int              { return 4 }
func (t *I32) Align() int             { return 4 }
func (t *I32) Kind() TypeKind         { return TypeKindI32 }
func (t *I32) Equal(u ValueType) bool { _, ok := u.(*I32); return ok }

/**************************************
U32:
**************************************/
type U32 struct {
	tCommon
}

func (tl *Types) GenU32(name string) *U32 {
	nt := U32{}
	if len(name) > 0 {
		nt.name = name
	} else {
		nt.name = "u32"
	}
	t, ok := tl.Lookup(nt.name)
	if ok {
		return t.(*U32)
	}

	tl.Add(&nt)
	return &nt
}
func (t *U32) Size() int              { return 4 }
func (t *U32) Align() int             { return 4 }
func (t *U32) Kind() TypeKind         { return TypeKindU32 }
func (t *U32) Equal(u ValueType) bool { _, ok := u.(*U32); return ok }

/**************************************
I64:
**************************************/
type I64 struct {
	tCommon
}

func (tl *Types) GenValueI64(name string) *I64 {
	nt := I64{}
	if len(name) > 0 {
		nt.name = name
	} else {
		nt.name = "i64"
	}
	t, ok := tl.Lookup(nt.name)
	if ok {
		return t.(*I64)
	}

	tl.Add(&nt)
	return &nt
}
func (t *I64) Size() int              { return 8 }
func (t *I64) Align() int             { return 8 }
func (t *I64) Kind() TypeKind         { return TypeKindI64 }
func (t *I64) Equal(u ValueType) bool { _, ok := u.(*I64); return ok }

/**************************************
U64:
**************************************/
type U64 struct {
	tCommon
}

func (tl *Types) GenU64(name string) *U64 {
	nt := U64{}
	if len(name) > 0 {
		nt.name = name
	} else {
		nt.name = "u64"
	}
	t, ok := tl.Lookup(nt.name)
	if ok {
		return t.(*U64)
	}

	tl.Add(&nt)
	return &nt
}
func (t *U64) Size() int              { return 8 }
func (t *U64) Align() int             { return 8 }
func (t *U64) Kind() TypeKind         { return TypeKindU64 }
func (t *U64) Equal(u ValueType) bool { _, ok := u.(*U64); return ok }

/**************************************
F32:
**************************************/
type F32 struct {
	tCommon
}

func (tl *Types) GenF32(name string) *F32 {
	nt := F32{}
	if len(name) > 0 {
		nt.name = name
	} else {
		nt.name = "f32"
	}
	t, ok := tl.Lookup(nt.name)
	if ok {
		return t.(*F32)
	}

	tl.Add(&nt)
	return &nt
}
func (t *F32) Size() int              { return 4 }
func (t *F32) Align() int             { return 4 }
func (t *F32) Kind() TypeKind         { return TypeKindF32 }
func (t *F32) Equal(u ValueType) bool { _, ok := u.(*F32); return ok }

/**************************************
F64:
**************************************/
type F64 struct {
	tCommon
}

func (tl *Types) GenF64(name string) *F64 {
	nt := F64{}
	if len(name) > 0 {
		nt.name = name
	} else {
		nt.name = "f64"
	}
	t, ok := tl.Lookup(nt.name)
	if ok {
		return t.(*F64)
	}

	tl.Add(&nt)
	return &nt
}
func (t *F64) Size() int              { return 8 }
func (t *F64) Align() int             { return 8 }
func (t *F64) Kind() TypeKind         { return TypeKindF64 }
func (t *F64) Equal(u ValueType) bool { _, ok := u.(*F64); return ok }
