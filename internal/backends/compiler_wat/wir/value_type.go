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
	case *tI32, *tRune: //Todo: *tI8, *tI16*
		return wat.I32{}

	case *tU32, *tU8, *tU16, *tBool:
		return wat.U32{}

	case *tI64:
		return wat.I64{}

	case *tU64:
		return wat.U64{}

	case *tF32:
		return wat.F32{}

	case *tF64:
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
tVoid:
**************************************/
type tVoid struct {
	tCommon
}

func (m *Module) GenValueType_void(name string) *tVoid {
	nt := tVoid{}
	if len(name) > 0 {
		nt.name = name
	} else {
		nt.name = "void"
	}
	t, ok := m.findValueType(nt.name)
	if ok {
		return t.(*tVoid)
	}

	m.addValueType(&nt)
	return &nt
}
func (t *tVoid) Size() int              { return 0 }
func (t *tVoid) align() int             { return 0 }
func (t *tVoid) Kind() TypeKind         { return kVoid }
func (t *tVoid) onFree() int            { return 0 }
func (t *tVoid) Raw() []wat.ValueType   { return []wat.ValueType{} }
func (t *tVoid) Equal(u ValueType) bool { _, ok := u.(*tVoid); return ok }
func (t *tVoid) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	logger.Fatal("Unreachable")
	return nil
}

/**************************************
tBool:
**************************************/
type tBool struct {
	tCommon
}

func (m *Module) GenValueType_bool(name string) *tBool {
	nt := tBool{}
	if len(name) > 0 {
		nt.name = name
	} else {
		nt.name = "bool"
	}
	t, ok := m.findValueType(nt.name)
	if ok {
		return t.(*tBool)
	}

	m.addValueType(&nt)
	return &nt
}
func (t *tBool) Size() int              { return 1 }
func (t *tBool) align() int             { return 1 }
func (t *tBool) Kind() TypeKind         { return kBool }
func (t *tBool) onFree() int            { return 0 }
func (t *tBool) Raw() []wat.ValueType   { return []wat.ValueType{wat.I32{}} }
func (t *tBool) Equal(u ValueType) bool { _, ok := u.(*tBool); return ok }
func (t *tBool) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	//if !addr.Type().(*Ptr).Base.Equal(t) {
	//	logger.Fatal("Type not match")
	//	return nil
	//}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 1))
	return insts
}

/**************************************
tRune:
**************************************/
type tRune struct {
	tCommon
}

func (m *Module) GenValueType_rune(name string) *tRune {
	nt := tRune{}
	if len(name) > 0 {
		nt.name = name
	} else {
		nt.name = "rune"
	}
	t, ok := m.findValueType(nt.name)
	if ok {
		return t.(*tRune)
	}

	m.addValueType(&nt)
	return &nt
}
func (t *tRune) Size() int              { return 4 }
func (t *tRune) align() int             { return 4 }
func (t *tRune) Kind() TypeKind         { return kRune }
func (t *tRune) onFree() int            { return 0 }
func (t *tRune) Raw() []wat.ValueType   { return []wat.ValueType{wat.I32{}} }
func (t *tRune) Equal(u ValueType) bool { _, ok := u.(*tRune); return ok }
func (t *tRune) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	//if !addr.Type().(*Ptr).Base.Equal(t) {
	//	logger.Fatal("Type not match")
	//	return nil
	//}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 1))
	return insts
}

/**************************************
tI8:
**************************************/
type tI8 struct {
	tCommon
}

func (m *Module) GenValueType_i8(name string) *tI8 {
	nt := tI8{}
	if len(name) > 0 {
		nt.name = name
	} else {
		nt.name = "i8"
	}
	t, ok := m.findValueType(nt.name)
	if ok {
		return t.(*tI8)
	}

	m.addValueType(&nt)
	return &nt
}
func (t *tI8) Size() int              { return 1 }
func (t *tI8) align() int             { return 1 }
func (t *tI8) Kind() TypeKind         { return kI8 }
func (t *tI8) onFree() int            { return 0 }
func (t *tI8) Raw() []wat.ValueType   { return []wat.ValueType{wat.I32{}} }
func (t *tI8) Equal(u ValueType) bool { _, ok := u.(*tI8); return ok }
func (t *tI8) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	//if !addr.Type().(*Ptr).Base.Equal(t) {
	//	logger.Fatal("Type not match")
	//	return nil
	//}
	if _, ok := addr.(*aPtr); !ok {
		logger.Fatal("addr should be `*aPtr`")
	}

	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad8s(offset, 1))
	return insts
}

/**************************************
tU8:
**************************************/
type tU8 struct {
	tCommon
}

func (m *Module) GenValueType_u8(name string) *tU8 {
	nt := tU8{}
	if len(name) > 0 {
		nt.name = name
	} else {
		nt.name = "u8"
	}
	t, ok := m.findValueType(nt.name)
	if ok {
		return t.(*tU8)
	}

	m.addValueType(&nt)
	return &nt
}
func (t *tU8) Size() int              { return 1 }
func (t *tU8) align() int             { return 1 }
func (t *tU8) Kind() TypeKind         { return kU8 }
func (t *tU8) onFree() int            { return 0 }
func (t *tU8) Raw() []wat.ValueType   { return []wat.ValueType{wat.I32{}} }
func (t *tU8) Equal(u ValueType) bool { _, ok := u.(*tU8); return ok }
func (t *tU8) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	//if !addr.Type().(*Ptr).Base.Equal(t) {
	//	logger.Fatal("Type not match")
	//	return nil
	//}
	if _, ok := addr.(*aPtr); !ok {
		logger.Fatal("addr should be `*aPtr`")
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad8u(offset, 1))
	return insts
}

/**************************************
tI16:
**************************************/
type tI16 struct {
	tCommon
}

func (m *Module) GenValueType_i16(name string) *tI16 {
	nt := tI16{}
	if len(name) > 0 {
		nt.name = name
	} else {
		nt.name = "i16"
	}
	t, ok := m.findValueType(nt.name)
	if ok {
		return t.(*tI16)
	}

	m.addValueType(&nt)
	return &nt
}
func (t *tI16) Size() int              { return 2 }
func (t *tI16) align() int             { return 2 }
func (t *tI16) Kind() TypeKind         { return kI16 }
func (t *tI16) onFree() int            { return 0 }
func (t *tI16) Raw() []wat.ValueType   { return []wat.ValueType{wat.I32{}} }
func (t *tI16) Equal(u ValueType) bool { _, ok := u.(*tI16); return ok }
func (t *tI16) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	//if !addr.Type().(*Ptr).Base.Equal(t) {
	//	logger.Fatal("Type not match")
	//	return nil
	//}
	if _, ok := addr.(*aPtr); !ok {
		logger.Fatal("addr should be `*aPtr`")
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad16s(offset, 1))
	return insts
}

/**************************************
tU16:
**************************************/
type tU16 struct {
	tCommon
}

func (m *Module) GenValueType_u16(name string) *tU16 {
	nt := tU16{}
	if len(name) > 0 {
		nt.name = name
	} else {
		nt.name = "u16"
	}
	t, ok := m.findValueType(nt.name)
	if ok {
		return t.(*tU16)
	}

	m.addValueType(&nt)
	return &nt
}
func (t *tU16) Size() int              { return 2 }
func (t *tU16) align() int             { return 2 }
func (t *tU16) Kind() TypeKind         { return kU16 }
func (t *tU16) onFree() int            { return 0 }
func (t *tU16) Raw() []wat.ValueType   { return []wat.ValueType{wat.I32{}} }
func (t *tU16) Equal(u ValueType) bool { _, ok := u.(*tU16); return ok }
func (t *tU16) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	//if !addr.Type().(*Ptr).Base.Equal(t) {
	//	logger.Fatal("Type not match")
	//	return nil
	//}
	if _, ok := addr.(*aPtr); !ok {
		logger.Fatal("addr should be `*aPtr`")
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad16u(offset, 1))
	return insts
}

/**************************************
tI32:
**************************************/
type tI32 struct {
	tCommon
}

func (m *Module) GenValueType_int(name string) ValueType {
	return m.GenValueType_i32(name)
}

func (m *Module) GenValueType_i32(name string) *tI32 {
	nt := tI32{}
	if len(name) > 0 {
		nt.name = name
	} else {
		nt.name = "i32"
	}
	t, ok := m.findValueType(nt.name)
	if ok {
		return t.(*tI32)
	}

	m.addValueType(&nt)
	return &nt
}
func (t *tI32) Size() int              { return 4 }
func (t *tI32) align() int             { return 4 }
func (t *tI32) Kind() TypeKind         { return kI32 }
func (t *tI32) onFree() int            { return 0 }
func (t *tI32) Raw() []wat.ValueType   { return []wat.ValueType{wat.I32{}} }
func (t *tI32) Equal(u ValueType) bool { _, ok := u.(*tI32); return ok }
func (t *tI32) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	//if !addr.Type().(*Ptr).Base.Equal(t) {
	//	logger.Fatal("Type not match")
	//	return nil
	//}
	if _, ok := addr.(*aPtr); !ok {
		logger.Fatal("addr should be `*aPtr`")
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 1))
	return insts
}

/**************************************
tU32:
**************************************/
type tU32 struct {
	tCommon
}

func (m *Module) GenValueType_uint(name string) ValueType {
	return m.GenValueType_u32(name)
}

func (m *Module) GenValueType_u32(name string) *tU32 {
	nt := tU32{}
	if len(name) > 0 {
		nt.name = name
	} else {
		nt.name = "u32"
	}
	t, ok := m.findValueType(nt.name)
	if ok {
		return t.(*tU32)
	}

	m.addValueType(&nt)
	return &nt
}
func (t *tU32) Size() int              { return 4 }
func (t *tU32) align() int             { return 4 }
func (t *tU32) Kind() TypeKind         { return kU32 }
func (t *tU32) onFree() int            { return 0 }
func (t *tU32) Raw() []wat.ValueType   { return []wat.ValueType{wat.U32{}} }
func (t *tU32) Equal(u ValueType) bool { _, ok := u.(*tU32); return ok }
func (t *tU32) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	//if !addr.Type().(*Ptr).Base.Equal(t) {
	//	logger.Fatal("Type not match")
	//	return nil
	//}
	if _, ok := addr.(*aPtr); !ok {
		logger.Fatal("addr should be `*aPtr`")
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 1))
	return insts
}

/**************************************
tI64:
**************************************/
type tI64 struct {
	tCommon
}

func (m *Module) GenValueType_i64(name string) *tI64 {
	nt := tI64{}
	if len(name) > 0 {
		nt.name = name
	} else {
		nt.name = "i64"
	}
	t, ok := m.findValueType(nt.name)
	if ok {
		return t.(*tI64)
	}

	m.addValueType(&nt)
	return &nt
}
func (t *tI64) Size() int              { return 8 }
func (t *tI64) align() int             { return 8 }
func (t *tI64) Kind() TypeKind         { return kI64 }
func (t *tI64) onFree() int            { return 0 }
func (t *tI64) Raw() []wat.ValueType   { return []wat.ValueType{wat.I64{}} }
func (t *tI64) Equal(u ValueType) bool { _, ok := u.(*tI64); return ok }
func (t *tI64) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	//if !addr.Type().(*Ptr).Base.Equal(t) {
	//	logger.Fatal("Type not match")
	//	return nil
	//}
	if _, ok := addr.(*aPtr); !ok {
		logger.Fatal("addr should be `*aPtr`")
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 1))
	return insts
}

/**************************************
tU64:
**************************************/
type tU64 struct {
	tCommon
}

func (m *Module) GenValueType_u64(name string) *tU64 {
	nt := tU64{}
	if len(name) > 0 {
		nt.name = name
	} else {
		nt.name = "u64"
	}
	t, ok := m.findValueType(nt.name)
	if ok {
		return t.(*tU64)
	}

	m.addValueType(&nt)
	return &nt
}
func (t *tU64) Size() int              { return 8 }
func (t *tU64) align() int             { return 8 }
func (t *tU64) Kind() TypeKind         { return kU64 }
func (t *tU64) onFree() int            { return 0 }
func (t *tU64) Raw() []wat.ValueType   { return []wat.ValueType{wat.U64{}} }
func (t *tU64) Equal(u ValueType) bool { _, ok := u.(*tU64); return ok }
func (t *tU64) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	//if !addr.Type().(*Ptr).Base.Equal(t) {
	//	logger.Fatal("Type not match")
	//	return nil
	//}
	if _, ok := addr.(*aPtr); !ok {
		logger.Fatal("addr should be `*aPtr`")
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 1))
	return insts
}

/**************************************
tF32:
**************************************/
type tF32 struct {
	tCommon
}

func (m *Module) GenValueType_f32(name string) *tF32 {
	nt := tF32{}
	if len(name) > 0 {
		nt.name = name
	} else {
		nt.name = "f32"
	}
	t, ok := m.findValueType(nt.name)
	if ok {
		return t.(*tF32)
	}

	m.addValueType(&nt)
	return &nt
}
func (t *tF32) Size() int              { return 4 }
func (t *tF32) align() int             { return 4 }
func (t *tF32) Kind() TypeKind         { return kF32 }
func (t *tF32) onFree() int            { return 0 }
func (t *tF32) Raw() []wat.ValueType   { return []wat.ValueType{wat.F32{}} }
func (t *tF32) Equal(u ValueType) bool { _, ok := u.(*tF32); return ok }
func (t *tF32) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	//if !addr.Type().(*Ptr).Base.Equal(t) {
	//	logger.Fatal("Type not match")
	//	return nil
	//}
	if _, ok := addr.(*aPtr); !ok {
		logger.Fatal("addr should be `*aPtr`")
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 1))
	return insts
}

/**************************************
tF64:
**************************************/
type tF64 struct {
	tCommon
}

func (m *Module) GenValueType_f64(name string) *tF64 {
	nt := tF64{}
	if len(name) > 0 {
		nt.name = name
	} else {
		nt.name = "f64"
	}
	t, ok := m.findValueType(nt.name)
	if ok {
		return t.(*tF64)
	}

	m.addValueType(&nt)
	return &nt
}
func (t *tF64) Size() int              { return 8 }
func (t *tF64) align() int             { return 8 }
func (t *tF64) Kind() TypeKind         { return kF64 }
func (t *tF64) onFree() int            { return 0 }
func (t *tF64) Raw() []wat.ValueType   { return []wat.ValueType{wat.F64{}} }
func (t *tF64) Equal(u ValueType) bool { _, ok := u.(*tF64); return ok }
func (t *tF64) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	//if !addr.Type().(*Ptr).Base.Equal(t) {
	//	logger.Fatal("Type not match")
	//	return nil
	//}
	if _, ok := addr.(*aPtr); !ok {
		logger.Fatal("addr should be `*aPtr`")
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 1))
	return insts
}
