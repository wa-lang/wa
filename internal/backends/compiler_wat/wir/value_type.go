// 版权 @2022 凹语言 作者。保留所有权利。

package wir

import (
	"wa-lang.org/wa/internal/backends/compiler_wat/wir/wat"
	"wa-lang.org/wa/internal/logger"
)

type TypeKind uint8

const (
	kUnknown TypeKind = iota
	kVoid
	kBool
	kU8
	kU16
	kU32
	kU64
	kI8
	kI16
	kI32
	kI64
	kF32
	kF64
	kComplex64
	kComplex128
	kRune
	kPtr
	kBlock
	kStruct
	kTuple
	kRef
	kString
	kSlice
	kArray
	kMap
	kInterface
	kDup
)

func toWatType(t ValueType) wat.ValueType {
	switch t.(type) {
	case *I32, *Rune: //Todo: *tI8, *tI16*
		return wat.I32{}

	case *U32, *U8, *U16, *Bool:
		return wat.U32{}

	case *I64:
		return wat.I64{}

	case *U64:
		return wat.U64{}

	case *F32:
		return wat.F32{}

	case *F64:
		return wat.F64{}

	case *Ptr:
		return wat.U32{}

	case *Block:
		return wat.I32{}

	default:
		logger.Fatalf("Todo:%v\n", t)
	}

	return nil
}

/**************************************
tCommon:
**************************************/
type tCommon struct {
	name    string
	addr    int
	hash    int
	comp    int
	methods []Method
}

func (t *tCommon) Named() string       { return t.name }
func (t *tCommon) Hash() int           { return t.hash }
func (t *tCommon) SetHash(h int)       { t.hash = h }
func (t *tCommon) AddMethod(m Method)  { t.methods = append(t.methods, m) }
func (t *tCommon) NumMethods() int     { return len(t.methods) }
func (t *tCommon) Method(i int) Method { return t.methods[i] }
func (t *tCommon) typeInfoAddr() int   { return t.addr }
func (t *tCommon) OnComp() int         { return t.comp }
func (t *tCommon) setOnComp(c int)     { t.comp = c }

/**************************************
Void:
**************************************/
type Void struct {
	tCommon
}

func (m *Module) GenValueType_void(name string) *Void {
	nt := Void{}
	if len(name) > 0 {
		nt.name = name
	} else {
		nt.name = "void"
	}
	t, ok := m.findValueType(nt.name)
	if ok {
		return t.(*Void)
	}

	m.addValueType(&nt)
	return &nt
}
func (t *Void) Size() int              { return 0 }
func (t *Void) align() int             { return 0 }
func (t *Void) Kind() TypeKind         { return kVoid }
func (t *Void) OnFree() int            { return 0 }
func (t *Void) Raw() []wat.ValueType   { return []wat.ValueType{} }
func (t *Void) Equal(u ValueType) bool { _, ok := u.(*Void); return ok }
func (t *Void) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	logger.Fatal("Unreachable")
	return nil
}
func (t *Void) EmitLoadFromAddrNoRetain(addr Value, offset int) []wat.Inst {
	logger.Fatal("Unreachable")
	return nil
}

/**************************************
Bool:
**************************************/
type Bool struct {
	tCommon
}

func (m *Module) GenValueType_bool(name string) *Bool {
	nt := Bool{}
	if len(name) > 0 {
		nt.name = name
	} else {
		nt.name = "bool"
	}
	t, ok := m.findValueType(nt.name)
	if ok {
		return t.(*Bool)
	}

	m.addValueType(&nt)
	return &nt
}
func (t *Bool) Size() int              { return 1 }
func (t *Bool) align() int             { return 1 }
func (t *Bool) Kind() TypeKind         { return kBool }
func (t *Bool) OnFree() int            { return 0 }
func (t *Bool) Raw() []wat.ValueType   { return []wat.ValueType{wat.I32{}} }
func (t *Bool) Equal(u ValueType) bool { _, ok := u.(*Bool); return ok }
func (t *Bool) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad8u(offset, 1))
	return insts
}
func (t *Bool) EmitLoadFromAddrNoRetain(addr Value, offset int) []wat.Inst {
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad8u(offset, 1))
	return insts
}

/**************************************
tRune:
**************************************/
type Rune struct {
	tCommon
}

func (m *Module) GenValueType_rune(name string) *Rune {
	nt := Rune{}
	if len(name) > 0 {
		nt.name = name
	} else {
		nt.name = "rune"
	}
	t, ok := m.findValueType(nt.name)
	if ok {
		return t.(*Rune)
	}

	m.addValueType(&nt)
	return &nt
}
func (t *Rune) Size() int              { return 4 }
func (t *Rune) align() int             { return 4 }
func (t *Rune) Kind() TypeKind         { return kRune }
func (t *Rune) OnFree() int            { return 0 }
func (t *Rune) Raw() []wat.ValueType   { return []wat.ValueType{wat.I32{}} }
func (t *Rune) Equal(u ValueType) bool { _, ok := u.(*Rune); return ok }
func (t *Rune) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 4))
	return insts
}
func (t *Rune) EmitLoadFromAddrNoRetain(addr Value, offset int) []wat.Inst {
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 4))
	return insts
}

/**************************************
I8:
**************************************/
type I8 struct {
	tCommon
}

func (m *Module) GenValueType_i8(name string) *I8 {
	nt := I8{}
	if len(name) > 0 {
		nt.name = name
	} else {
		nt.name = "i8"
	}
	t, ok := m.findValueType(nt.name)
	if ok {
		return t.(*I8)
	}

	m.addValueType(&nt)
	return &nt
}
func (t *I8) Size() int              { return 1 }
func (t *I8) align() int             { return 1 }
func (t *I8) Kind() TypeKind         { return kI8 }
func (t *I8) OnFree() int            { return 0 }
func (t *I8) Raw() []wat.ValueType   { return []wat.ValueType{wat.I32{}} }
func (t *I8) Equal(u ValueType) bool { _, ok := u.(*I8); return ok }
func (t *I8) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	if _, ok := addr.(*aPtr); !ok {
		logger.Fatal("addr should be `*aPtr`")
	}

	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad8s(offset, 1))
	return insts
}
func (t *I8) EmitLoadFromAddrNoRetain(addr Value, offset int) []wat.Inst {
	if _, ok := addr.(*aPtr); !ok {
		logger.Fatal("addr should be `*aPtr`")
	}

	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad8s(offset, 1))
	return insts
}

/**************************************
U8:
**************************************/
type U8 struct {
	tCommon
}

func (m *Module) GenValueType_u8(name string) *U8 {
	nt := U8{}
	if len(name) > 0 {
		nt.name = name
	} else {
		nt.name = "u8"
	}
	t, ok := m.findValueType(nt.name)
	if ok {
		return t.(*U8)
	}

	m.addValueType(&nt)
	return &nt
}
func (t *U8) Size() int              { return 1 }
func (t *U8) align() int             { return 1 }
func (t *U8) Kind() TypeKind         { return kU8 }
func (t *U8) OnFree() int            { return 0 }
func (t *U8) Raw() []wat.ValueType   { return []wat.ValueType{wat.I32{}} }
func (t *U8) Equal(u ValueType) bool { _, ok := u.(*U8); return ok }
func (t *U8) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	if _, ok := addr.(*aPtr); !ok {
		logger.Fatal("addr should be `*aPtr`")
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad8u(offset, 1))
	return insts
}
func (t *U8) EmitLoadFromAddrNoRetain(addr Value, offset int) []wat.Inst {
	if _, ok := addr.(*aPtr); !ok {
		logger.Fatal("addr should be `*aPtr`")
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad8u(offset, 1))
	return insts
}

/**************************************
I16:
**************************************/
type I16 struct {
	tCommon
}

func (m *Module) GenValueType_i16(name string) *I16 {
	nt := I16{}
	if len(name) > 0 {
		nt.name = name
	} else {
		nt.name = "i16"
	}
	t, ok := m.findValueType(nt.name)
	if ok {
		return t.(*I16)
	}

	m.addValueType(&nt)
	return &nt
}
func (t *I16) Size() int              { return 2 }
func (t *I16) align() int             { return 2 }
func (t *I16) Kind() TypeKind         { return kI16 }
func (t *I16) OnFree() int            { return 0 }
func (t *I16) Raw() []wat.ValueType   { return []wat.ValueType{wat.I32{}} }
func (t *I16) Equal(u ValueType) bool { _, ok := u.(*I16); return ok }
func (t *I16) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	if _, ok := addr.(*aPtr); !ok {
		logger.Fatal("addr should be `*aPtr`")
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad16s(offset, 2))
	return insts
}
func (t *I16) EmitLoadFromAddrNoRetain(addr Value, offset int) []wat.Inst {
	if _, ok := addr.(*aPtr); !ok {
		logger.Fatal("addr should be `*aPtr`")
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad16s(offset, 2))
	return insts
}

/**************************************
tU16:
**************************************/
type U16 struct {
	tCommon
}

func (m *Module) GenValueType_u16(name string) *U16 {
	nt := U16{}
	if len(name) > 0 {
		nt.name = name
	} else {
		nt.name = "u16"
	}
	t, ok := m.findValueType(nt.name)
	if ok {
		return t.(*U16)
	}

	m.addValueType(&nt)
	return &nt
}
func (t *U16) Size() int              { return 2 }
func (t *U16) align() int             { return 2 }
func (t *U16) Kind() TypeKind         { return kU16 }
func (t *U16) OnFree() int            { return 0 }
func (t *U16) Raw() []wat.ValueType   { return []wat.ValueType{wat.I32{}} }
func (t *U16) Equal(u ValueType) bool { _, ok := u.(*U16); return ok }
func (t *U16) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	if _, ok := addr.(*aPtr); !ok {
		logger.Fatal("addr should be `*aPtr`")
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad16u(offset, 2))
	return insts
}
func (t *U16) EmitLoadFromAddrNoRetain(addr Value, offset int) []wat.Inst {
	if _, ok := addr.(*aPtr); !ok {
		logger.Fatal("addr should be `*aPtr`")
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad16u(offset, 2))
	return insts
}

/**************************************
I32:
**************************************/
type I32 struct {
	tCommon
}

func (m *Module) GenValueType_int(name string) ValueType {
	return m.GenValueType_i32(name)
}

func (m *Module) GenValueType_i32(name string) *I32 {
	nt := I32{}
	if len(name) > 0 {
		nt.name = name
	} else {
		nt.name = "i32"
	}
	t, ok := m.findValueType(nt.name)
	if ok {
		return t.(*I32)
	}

	m.addValueType(&nt)
	return &nt
}
func (t *I32) Size() int              { return 4 }
func (t *I32) align() int             { return 4 }
func (t *I32) Kind() TypeKind         { return kI32 }
func (t *I32) OnFree() int            { return 0 }
func (t *I32) Raw() []wat.ValueType   { return []wat.ValueType{wat.I32{}} }
func (t *I32) Equal(u ValueType) bool { _, ok := u.(*I32); return ok }
func (t *I32) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	if _, ok := addr.(*aPtr); !ok {
		logger.Fatal("addr should be `*aPtr`")
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 4))
	return insts
}
func (t *I32) EmitLoadFromAddrNoRetain(addr Value, offset int) []wat.Inst {
	if _, ok := addr.(*aPtr); !ok {
		logger.Fatal("addr should be `*aPtr`")
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 4))
	return insts
}

/**************************************
U32:
**************************************/
type U32 struct {
	tCommon
}

func (m *Module) GenValueType_uint(name string) ValueType {
	return m.GenValueType_u32(name)
}

func (m *Module) GenValueType_u32(name string) *U32 {
	nt := U32{}
	if len(name) > 0 {
		nt.name = name
	} else {
		nt.name = "u32"
	}
	t, ok := m.findValueType(nt.name)
	if ok {
		return t.(*U32)
	}

	m.addValueType(&nt)
	return &nt
}
func (t *U32) Size() int              { return 4 }
func (t *U32) align() int             { return 4 }
func (t *U32) Kind() TypeKind         { return kU32 }
func (t *U32) OnFree() int            { return 0 }
func (t *U32) Raw() []wat.ValueType   { return []wat.ValueType{wat.U32{}} }
func (t *U32) Equal(u ValueType) bool { _, ok := u.(*U32); return ok }
func (t *U32) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	if _, ok := addr.(*aPtr); !ok {
		logger.Fatal("addr should be `*aPtr`")
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 4))
	return insts
}
func (t *U32) EmitLoadFromAddrNoRetain(addr Value, offset int) []wat.Inst {
	if _, ok := addr.(*aPtr); !ok {
		logger.Fatal("addr should be `*aPtr`")
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 4))
	return insts
}

/**************************************
I64:
**************************************/
type I64 struct {
	tCommon
}

func (m *Module) GenValueType_i64(name string) *I64 {
	nt := I64{}
	if len(name) > 0 {
		nt.name = name
	} else {
		nt.name = "i64"
	}
	t, ok := m.findValueType(nt.name)
	if ok {
		return t.(*I64)
	}

	m.addValueType(&nt)
	return &nt
}
func (t *I64) Size() int              { return 8 }
func (t *I64) align() int             { return 8 }
func (t *I64) Kind() TypeKind         { return kI64 }
func (t *I64) OnFree() int            { return 0 }
func (t *I64) Raw() []wat.ValueType   { return []wat.ValueType{wat.I64{}} }
func (t *I64) Equal(u ValueType) bool { _, ok := u.(*I64); return ok }
func (t *I64) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	if _, ok := addr.(*aPtr); !ok {
		logger.Fatal("addr should be `*aPtr`")
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 8))
	return insts
}
func (t *I64) EmitLoadFromAddrNoRetain(addr Value, offset int) []wat.Inst {
	if _, ok := addr.(*aPtr); !ok {
		logger.Fatal("addr should be `*aPtr`")
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 8))
	return insts
}

/**************************************
U64:
**************************************/
type U64 struct {
	tCommon
}

func (m *Module) GenValueType_u64(name string) *U64 {
	nt := U64{}
	if len(name) > 0 {
		nt.name = name
	} else {
		nt.name = "u64"
	}
	t, ok := m.findValueType(nt.name)
	if ok {
		return t.(*U64)
	}

	m.addValueType(&nt)
	return &nt
}
func (t *U64) Size() int              { return 8 }
func (t *U64) align() int             { return 8 }
func (t *U64) Kind() TypeKind         { return kU64 }
func (t *U64) OnFree() int            { return 0 }
func (t *U64) Raw() []wat.ValueType   { return []wat.ValueType{wat.U64{}} }
func (t *U64) Equal(u ValueType) bool { _, ok := u.(*U64); return ok }
func (t *U64) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	if _, ok := addr.(*aPtr); !ok {
		logger.Fatal("addr should be `*aPtr`")
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 8))
	return insts
}
func (t *U64) EmitLoadFromAddrNoRetain(addr Value, offset int) []wat.Inst {
	if _, ok := addr.(*aPtr); !ok {
		logger.Fatal("addr should be `*aPtr`")
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 8))
	return insts
}

/**************************************
F32:
**************************************/
type F32 struct {
	tCommon
}

func (m *Module) GenValueType_f32(name string) *F32 {
	nt := F32{}
	if len(name) > 0 {
		nt.name = name
	} else {
		nt.name = "f32"
	}
	t, ok := m.findValueType(nt.name)
	if ok {
		return t.(*F32)
	}

	m.addValueType(&nt)
	return &nt
}
func (t *F32) Size() int              { return 4 }
func (t *F32) align() int             { return 4 }
func (t *F32) Kind() TypeKind         { return kF32 }
func (t *F32) OnFree() int            { return 0 }
func (t *F32) Raw() []wat.ValueType   { return []wat.ValueType{wat.F32{}} }
func (t *F32) Equal(u ValueType) bool { _, ok := u.(*F32); return ok }
func (t *F32) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	if _, ok := addr.(*aPtr); !ok {
		logger.Fatal("addr should be `*aPtr`")
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 4))
	return insts
}
func (t *F32) EmitLoadFromAddrNoRetain(addr Value, offset int) []wat.Inst {
	if _, ok := addr.(*aPtr); !ok {
		logger.Fatal("addr should be `*aPtr`")
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 4))
	return insts
}

/**************************************
F64:
**************************************/
type F64 struct {
	tCommon
}

func (m *Module) GenValueType_f64(name string) *F64 {
	nt := F64{}
	if len(name) > 0 {
		nt.name = name
	} else {
		nt.name = "f64"
	}
	t, ok := m.findValueType(nt.name)
	if ok {
		return t.(*F64)
	}

	m.addValueType(&nt)
	return &nt
}
func (t *F64) Size() int              { return 8 }
func (t *F64) align() int             { return 8 }
func (t *F64) Kind() TypeKind         { return kF64 }
func (t *F64) OnFree() int            { return 0 }
func (t *F64) Raw() []wat.ValueType   { return []wat.ValueType{wat.F64{}} }
func (t *F64) Equal(u ValueType) bool { _, ok := u.(*F64); return ok }
func (t *F64) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	if _, ok := addr.(*aPtr); !ok {
		logger.Fatal("addr should be `*aPtr`")
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 8))
	return insts
}
func (t *F64) EmitLoadFromAddrNoRetain(addr Value, offset int) []wat.Inst {
	if _, ok := addr.(*aPtr); !ok {
		logger.Fatal("addr should be `*aPtr`")
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 8))
	return insts
}
