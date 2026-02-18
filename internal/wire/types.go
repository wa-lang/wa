// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wire

import "fmt"

/**************************************
本文件定义了 wire 中与值类型相关的各种对象
**************************************/

//-------------------------------------

/**************************************
Types : 类型库，用于管理值类型，类型的名字是其身份标志
类型不可以直接声明，必须通过 TypeLib.GenXXX 系列方法创建
对于大部分 TypeLib.GenXXX()：
  - 若指定的类型已存在于类型库中，则返回现存的类型
**************************************/
type Types struct {
	Void, Bool, U8, U16, U32, U64, Uint, I8, I16, I32, I64, Int, F32, F64, Rune, String Type

	ptrs map[Type]*Ptr // Base->Ptr
	//chunks map[Type]*chunk // Base->chunk
	refs map[Type]*Ref // Base->Ref

	tuples  []*Tuple
	structs []*Struct
	ifaces  []*Interface
	nameds  map[string]*Named

	intSize int
}

// Init 初始化类型库
func (tl *Types) Init(intsize int) {
	tl.Void = &Void{}
	tl.Bool = &Bool{}
	tl.U8 = &U8{}
	tl.U16 = &U16{}
	tl.U32 = &U32{}
	tl.U64 = &U64{}
	tl.I8 = &I8{}
	tl.I16 = &I16{}
	tl.I32 = &I32{}
	tl.I64 = &I64{}

	switch intsize {
	case 4:
		tl.Uint = tl.U32
		tl.Int = tl.I32
	case 8:
		tl.Uint = tl.U64
		tl.Int = tl.I64
	default:
		panic(fmt.Sprintf("Invalid intsize: %d", intsize))
	}
	tl.intSize = intsize

	tl.F32 = &F32{}
	tl.F64 = &F64{}
	//tl.Complex64 = &Complex64{}
	//tl.Complex128 = &Complex128{}
	tl.Rune = &Rune{}

	tl.String = tl.genString()

	tl.ptrs = make(map[Type]*Ptr)
	//tl.chunks = make(map[Type]*chunk)
	tl.refs = make(map[Type]*Ref)
	tl.nameds = make(map[string]*Named)
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
func (t *Void) hasRef() bool      { return false }
func (t *Void) Underlying() Type  { return t }

/**************************************
Bool: 布尔，1字节
**************************************/
type Bool struct{}

func (t *Bool) Name() string      { return "bool" }
func (t *Bool) Kind() TypeKind    { return TypeKindBool }
func (t *Bool) Equal(u Type) bool { _, ok := u.(*Bool); return ok }
func (t *Bool) hasRef() bool      { return false }
func (t *Bool) Underlying() Type  { return t }

/**************************************
I8: 8位有符号整数，1字节
**************************************/
type I8 struct{}

func (t *I8) Name() string      { return "i8" }
func (t *I8) Kind() TypeKind    { return TypeKindI8 }
func (t *I8) Equal(u Type) bool { _, ok := u.(*I8); return ok }
func (t *I8) hasRef() bool      { return false }
func (t *I8) Underlying() Type  { return t }

/**************************************
U8: 8位无符号整数，1字节
**************************************/
type U8 struct{}

func (t *U8) Name() string      { return "u8" }
func (t *U8) Kind() TypeKind    { return TypeKindU8 }
func (t *U8) Equal(u Type) bool { _, ok := u.(*U8); return ok }
func (t *U8) hasRef() bool      { return false }
func (t *U8) Underlying() Type  { return t }

/**************************************
I16: 16位有符号整数，2字节
**************************************/
type I16 struct{}

func (t *I16) Name() string      { return "i16" }
func (t *I16) Kind() TypeKind    { return TypeKindI16 }
func (t *I16) Equal(u Type) bool { _, ok := u.(*I16); return ok }
func (t *I16) hasRef() bool      { return false }
func (t *I16) Underlying() Type  { return t }

/**************************************
tU16: 16位无符号整数，2字节
**************************************/
type U16 struct{}

func (t *U16) Name() string      { return "u16" }
func (t *U16) Kind() TypeKind    { return TypeKindU16 }
func (t *U16) Equal(u Type) bool { _, ok := u.(*U16); return ok }
func (t *U16) hasRef() bool      { return false }
func (t *U16) Underlying() Type  { return t }

/**************************************
I32: 32位有符号整数，4字节
**************************************/
type I32 struct{}

func (t *I32) Name() string      { return "i32" }
func (t *I32) Kind() TypeKind    { return TypeKindI32 }
func (t *I32) Equal(u Type) bool { _, ok := u.(*I32); return ok }
func (t *I32) hasRef() bool      { return false }
func (t *I32) Underlying() Type  { return t }

/**************************************
U32: 32位无符号整数，4字节
**************************************/
type U32 struct{}

func (t *U32) Name() string      { return "u32" }
func (t *U32) Kind() TypeKind    { return TypeKindU32 }
func (t *U32) Equal(u Type) bool { _, ok := u.(*U32); return ok }
func (t *U32) hasRef() bool      { return false }
func (t *U32) Underlying() Type  { return t }

/**************************************
I64: 64位有符号整数，8字节
**************************************/
type I64 struct{}

func (t *I64) Name() string      { return "i64" }
func (t *I64) Kind() TypeKind    { return TypeKindI64 }
func (t *I64) Equal(u Type) bool { _, ok := u.(*I64); return ok }
func (t *I64) hasRef() bool      { return false }
func (t *I64) Underlying() Type  { return t }

/**************************************
U64: 64位无符号整数，8字节
**************************************/
type U64 struct{}

func (t *U64) Name() string      { return "u64" }
func (t *U64) Kind() TypeKind    { return TypeKindU64 }
func (t *U64) Equal(u Type) bool { _, ok := u.(*U64); return ok }
func (t *U64) hasRef() bool      { return false }
func (t *U64) Underlying() Type  { return t }

/**************************************
Uint: 平台相关无符号整型
**************************************/
type Uint struct{}

func (t *Uint) Name() string      { return "uint" }
func (t *Uint) Kind() TypeKind    { return TypeKindUint }
func (t *Uint) Equal(u Type) bool { _, ok := u.(*Uint); return ok }
func (t *Uint) hasRef() bool      { return false }
func (t *Uint) Underlying() Type  { return t }

/**************************************
Int: 平台相关有符号整型
**************************************/
type Int struct{}

func (t *Int) Name() string      { return "int" }
func (t *Int) Kind() TypeKind    { return TypeKindInt }
func (t *Int) Equal(u Type) bool { _, ok := u.(*Int); return ok }
func (t *Int) hasRef() bool      { return false }
func (t *Int) Underlying() Type  { return t }

/**************************************
F32: 单精度浮点数，4字节
**************************************/
type F32 struct{}

func (t *F32) Name() string      { return "f32" }
func (t *F32) Kind() TypeKind    { return TypeKindF32 }
func (t *F32) Equal(u Type) bool { _, ok := u.(*F32); return ok }
func (t *F32) hasRef() bool      { return false }
func (t *F32) Underlying() Type  { return t }

/**************************************
F64: 双精度浮点数，8字节
**************************************/
type F64 struct{}

func (t *F64) Name() string      { return "f64" }
func (t *F64) Kind() TypeKind    { return TypeKindF64 }
func (t *F64) Equal(u Type) bool { _, ok := u.(*F64); return ok }
func (t *F64) hasRef() bool      { return false }
func (t *F64) Underlying() Type  { return t }

/**************************************
Complex64: 单精度复数，8字节
**************************************/
//type Complex64 struct{}
//
//func (t *Complex64) Name() string      { return "complex64" }
//func (t *Complex64) Kind() TypeKind    { return TypeKindComplex64 }
//func (t *Complex64) Equal(u Type) bool { _, ok := u.(*Complex64); return ok }
//func (t *Complex64) hasRef() bool      { return false }
//func (t *Complex64) underlying() Type  { return t }

/**************************************
Complex128: 双精度复数，16字节
**************************************/
//ype Complex128 struct{}
//
//unc (t *Complex128) Name() string      { return "complex128" }
//unc (t *Complex128) Kind() TypeKind    { return TypeKindComplex128 }
//unc (t *Complex128) Equal(u Type) bool { _, ok := u.(*Complex128); return ok }
//unc (t *Complex128) hasRef() bool      { return false }
//unc (t *Complex128) underlying() Type  { return t }

/**************************************
Rune: unicode字符，4字节
**************************************/
type Rune struct{}

func (t *Rune) Name() string      { return "rune" }
func (t *Rune) Kind() TypeKind    { return TypeKindRune }
func (t *Rune) Equal(u Type) bool { _, ok := u.(*Rune); return ok }
func (t *Rune) hasRef() bool      { return false }
func (t *Rune) Underlying() Type  { return t }

/**************************************
String: 字符串
**************************************/
type String struct {
	underlying Struct
}

func (t *String) Name() string      { return "string" }
func (t *String) Kind() TypeKind    { return TypeKindString }
func (t *String) Equal(u Type) bool { _, ok := u.(*String); return ok }
func (t *String) hasRef() bool      { return true }
func (t *String) Underlying() Type  { return &t.underlying }

func (tl *Types) genString() *String {
	nt := &String{}
	c := StructMember{Name: "c", Type: &chunk{Base: tl.U8, size: tl.intSize}, id: 0}
	d := StructMember{Name: "d", Type: &Ptr{Base: tl.U8, size: tl.intSize}, id: 1}
	l := StructMember{Name: "l", Type: tl.Uint, id: 2}

	nt.underlying.member = []StructMember{c, d, l}
	return nt
}

/**************************************
Ptr: 指针，长度取决于目标平台
**************************************/
type Ptr struct {
	Base Type
	size int
}

func (t *Ptr) Name() string   { return "*" + t.Base.Name() }
func (t *Ptr) Kind() TypeKind { return TypeKindPtr }
func (t *Ptr) Equal(u Type) bool {
	if ut, ok := u.(*Ptr); ok {
		if t == ut {
			return true
		}
		return t.Base.Equal(ut.Base)
	}
	return false
}
func (t *Ptr) hasRef() bool     { return false }
func (t *Ptr) Underlying() Type { return t }

func (tl *Types) genPtr(base Type) *Ptr {
	if t, ok := tl.ptrs[base]; ok {
		return t
	}

	nt := &Ptr{Base: base, size: tl.intSize}
	tl.ptrs[base] = nt
	return nt
}

/**************************************
chunk: 数据块，本质是指针，runtime 自动管理内存时使用
**************************************/
type chunk struct {
	Base Type
	size int
}

func (t *chunk) Name() string   { return t.Base.Name() + "$$chunk" }
func (t *chunk) Kind() TypeKind { return TypeKindChunk }
func (t *chunk) Equal(u Type) bool {
	if ut, ok := u.(*chunk); ok {
		if t == ut {
			return true
		}
		return t.Base.Equal(ut.Base)
	}
	return false
}
func (t *chunk) hasRef() bool     { return true }
func (t *chunk) Underlying() Type { return t }

//func (tl *Types) genChunk(base Type) *chunk {
//	if t, ok := tl.chunks[base]; ok {
//		return t
//	}
//
//	nt := &chunk{Base: base}
//	tl.chunks[base] = nt
//	return nt
//}

/**************************************
Tuple: 元组
**************************************/
type Tuple struct {
	fields   []Type
	hasChunk bool
}

func (t *Tuple) Name() string {
	name := "tuple{"
	for i, f := range t.fields {
		if i > 0 {
			name += ", "
		}
		name += f.Name()
	}
	name += "}"
	return name
}
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
func (t *Tuple) hasRef() bool     { return t.hasChunk }
func (t *Tuple) Underlying() Type { return t }

func (t *Tuple) Len() int       { return len(t.fields) }
func (t *Tuple) At(id int) Type { return t.fields[id] }

func (tl *Types) GenTuple(fields []Type) *Tuple {
	nt := &Tuple{fields: fields}
	for _, f := range fields {
		if f.hasRef() {
			nt.hasChunk = true
			break
		}
	}

	for _, t := range tl.tuples {
		if nt.Equal(t) {
			return t
		}
	}

	tl.tuples = append(tl.tuples, nt)
	return nt
}

/**************************************
Struct: 结构体
**************************************/
type Struct struct {
	member   []StructMember
	hasChunk bool
}

func (t *Struct) Name() string {
	name := "struct{"
	for i, f := range t.member {
		if i > 0 {
			name += ", "
		}
		name += f.String()
	}

	name += "}"
	return name
}
func (t *Struct) Kind() TypeKind { return TypeKindStruct }
func (t *Struct) Equal(u Type) bool {
	ut, ok := u.(*Struct)
	if !ok {
		return false
	}

	if len(t.member) != len(ut.member) {
		return false
	}

	for i := range t.member {
		if !t.member[i].Equal(&ut.member[i]) {
			return false
		}
	}

	return true
}
func (t *Struct) hasRef() bool     { return t.hasChunk }
func (t *Struct) Underlying() Type { return t }

func (t *Struct) Len() int               { return len(t.member) }
func (t *Struct) At(id int) StructMember { return t.member[id] }
func (t *Struct) setFieldId() {
	for i := range t.member {
		t.member[i].id = i
	}
}

func (tl *Types) GenStruct(fields []StructMember) *Struct {
	nt := &Struct{member: fields}
	for _, f := range fields {
		if f.Type.hasRef() {
			nt.hasChunk = true
			break
		}
	}

	for _, t := range tl.structs {
		if nt.Equal(t) {
			return t
		}
	}

	nt.setFieldId()
	tl.structs = append(tl.structs, nt)
	return nt
}

type StructMember struct {
	Name string
	Type Type
	id   int
}

func (i *StructMember) Equal(u *StructMember) bool { return i.Name == u.Name && i.Type.Equal(u.Type) }
func (i *StructMember) String() string             { return i.Name + " " + i.Type.Name() }

/**************************************
Ref: 引用类型，runtime 自动内存管理时，Ptr 转化为 Ref
**************************************/
type Ref struct {
	Base       Type
	underlying Struct
}

func (t *Ref) Name() string   { return t.Base.Name() + "$$ref" }
func (t *Ref) Kind() TypeKind { return TypeKindRef }
func (t *Ref) Equal(u Type) bool {
	if ut, ok := u.(*Ref); ok {
		if t == ut {
			return true
		}
		return t.Base.Equal(ut.Base)
	}
	return false
}
func (t *Ref) hasRef() bool     { return true }
func (t *Ref) Underlying() Type { return &t.underlying }

func (tl *Types) GenRef(base Type) *Ref {
	if t, ok := tl.refs[base]; ok {
		return t
	}

	nt := &Ref{Base: base}
	b := StructMember{Name: "c", Type: &chunk{Base: base, size: tl.intSize}, id: 0}
	d := StructMember{Name: "d", Type: &Ptr{Base: base, size: tl.intSize}, id: 1}
	nt.underlying.member = []StructMember{b, d}
	return nt
}

/**************************************
Interface: 接口类型
**************************************/
type Interface struct {
	methods []Method
}

func (t *Interface) Name() string   { return "interface{...}" }
func (t *Interface) Kind() TypeKind { return TypeKindInterface }
func (t *Interface) Equal(u Type) bool {
	ut, ok := u.(*Interface)
	if !ok {
		return false
	}

	if len(t.methods) != len(ut.methods) {
		return false
	}

	for i := range t.methods {
		if !t.methods[i].Equal(&ut.methods[i]) {
			return false
		}
	}

	return true
}
func (t *Interface) hasRef() bool     { return true }
func (t *Interface) Underlying() Type { panic("todo") }

func (t *Interface) NumMethods() int     { return len(t.methods) }
func (t *Interface) Method(i int) Method { return t.methods[i] }

func (tl *Types) GenInterface(methods []Method) *Interface {
	nt := &Interface{methods: methods}

	for _, t := range tl.ifaces {
		if nt.Equal(t) {
			return t
		}
	}

	tl.ifaces = append(tl.ifaces, nt)
	return nt
}

/**************************************
Named: 具名类型，具名类型指向某个具体类型(underlying)
 - 当具名类型指向接口时，它的方法集和接口保持一致，既不可调用 Named.AddMethod() 添加新方法；
 - 具名类型指向其它类型时，它可以拥有自己的方法集
**************************************/
type Named struct {
	name       string
	underlying Type
	methods    []Method
}

func (t *Named) Name() string   { return t.name }
func (t *Named) Kind() TypeKind { return t.underlying.Kind() }
func (t *Named) Equal(u Type) bool {
	ut, ok := u.(*Named)
	if !ok {
		return false
	}

	return t.underlying.Equal(ut.underlying)
}
func (t *Named) hasRef() bool { return t.underlying.hasRef() }

func (t *Named) Underlying() Type         { return t.underlying }
func (t *Named) SetUnderlying(utype Type) { t.underlying = utype }

func (t *Named) AddMethod(m Method) int {
	if t.underlying.Kind() == TypeKindInterface {
		panic("Can't call AddMethod() on Named whose underlying-type is Interface")
	}
	t.methods = append(t.methods, m)
	return len(t.methods) - 1
}
func (t *Named) NumMethods() int {
	if t.underlying.Kind() == TypeKindInterface {
		return t.underlying.(*Interface).NumMethods()
	}

	return len(t.methods)
}
func (t *Named) Method(i int) Method {
	if t.underlying.Kind() == TypeKindInterface {
		return t.underlying.(*Interface).Method(i)
	}

	return t.methods[i]
}

// 创建一个具名类型，注意：新建的具名类型未指向任何具体类型，必须使用 Named.SetUnderlying() 设置
func (tl *Types) GenNamed(name string) *Named {
	if t, ok := tl.nameds[name]; ok {
		return t
	}

	nt := &Named{name: name}
	tl.nameds[name] = nt
	return nt
}
