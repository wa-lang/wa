// 版权 @2024 凹语言 作者。保留所有权利。

package ast

// 指令对应接口
type Instruction interface {
	aInstruction()
}

type Ins_Unreachable struct{}
type Ins_Nop struct{}

type Ins_Block struct{}
type Ins_Loop struct{}
type Ins_If struct{}
type Ins_Else struct{}
type Ins_End struct{}
type Ins_Br struct {
	X string
}
type Ins_BrIf struct{}
type Ins_BrTable struct{}
type Ins_Return struct{}

type Ins_Call struct {
	Name string
	Args []Instruction
}

type Ins_CallIndirect struct{}
type Ins_Drop struct{}
type Ins_Select struct{}
type Ins_TypedSelect struct{}
type Ins_LocalGet struct {
	X string
}
type Ins_LocalSet struct {
	X string
}
type Ins_LocalTee struct {
	X string
}
type Ins_GlobalGet struct {
	X string
}
type Ins_GlobalSet struct {
	X string
}
type Ins_TableGet struct{}
type Ins_TableSet struct{}
type Ins_I32Load struct{ Offset uint32 }
type Ins_I64Load struct{ Offset uint32 }
type Ins_F32Load struct{ Offset uint32 }
type Ins_F64Load struct{ Offset uint32 }
type Ins_I32Load8S struct{ Offset uint32 }
type Ins_I32Load8U struct{ Offset uint32 }
type Ins_I32Load16S struct{ Offset uint32 }
type Ins_I32Load16U struct{ Offset uint32 }
type Ins_I64Load8S struct{ Offset uint32 }
type Ins_I64Load8U struct{ Offset uint32 }
type Ins_I64Load16S struct{ Offset uint32 }
type Ins_I64Load16U struct{ Offset uint32 }
type Ins_I64Load32S struct{ Offset uint32 }
type Ins_I64Load32U struct{ Offset uint32 }

type Ins_I32Store struct {
	Offset uint32
	Value  int32
}

type Ins_I64Store struct {
}
type Ins_F32Store struct {
	Offset uint32
	Value  float32
}
type Ins_F64Store struct {
	Offset uint32
	Value  float64
}

type Ins_I32Store8 struct {
	Offset uint32
	Value  int32
}
type Ins_I32Store16 struct {
	Offset uint32
	Value  int32
}
type Ins_I64Store8 struct {
	Offset uint32
	Value  int64
}
type Ins_I64Store16 struct {
	Offset uint32
	Value  int64
}
type Ins_I64Store32 struct {
	Offset uint32
	Value  int64
}

type Ins_MemorySize struct{}
type Ins_MemoryGrow struct{}

type Ins_I32Const struct{ X int32 }
type Ins_I64Const struct{ X int64 }
type Ins_F32Const struct{ X float32 }
type Ins_F64Const struct{ X float64 }

type Ins_I32Eqz struct{}
type Ins_I32Eq struct{}
type Ins_I32Ne struct{}
type Ins_I32LtS struct{}
type Ins_I32LtU struct{}
type Ins_I32GtS struct{}
type Ins_I32GtU struct{}
type Ins_I32LeS struct{}
type Ins_I32LeU struct{}
type Ins_I32GeS struct{}
type Ins_I32GeU struct{}

type Ins_I64Eqz struct{}
type Ins_I64Eq struct{}
type Ins_I64Ne struct{}
type Ins_I64LtS struct{}
type Ins_I64LtU struct{}
type Ins_I64GtS struct{}
type Ins_I64GtU struct{}
type Ins_I64LeS struct{}
type Ins_I64LeU struct{}
type Ins_I64GeS struct{}
type Ins_I64GeU struct{}

type Ins_F32Eq struct{}
type Ins_F32Ne struct{}
type Ins_F32Lt struct{}
type Ins_F32Gt struct{}
type Ins_F32Le struct{}
type Ins_F32Ge struct{}

type Ins_F64Eq struct{}
type Ins_F64Ne struct{}
type Ins_F64Lt struct{}
type Ins_F64Gt struct{}
type Ins_F64Le struct{}
type Ins_F64Ge struct{}

type Ins_I32Clz struct{}
type Ins_I32Ctz struct{}
type Ins_I32Popcnt struct{}
type Ins_I32Add struct{}
type Ins_I32Sub struct{}
type Ins_I32Mul struct{}
type Ins_I32DivS struct{}
type Ins_I32DivU struct{}
type Ins_I32RemS struct{}
type Ins_I32RemU struct{}
type Ins_I32And struct{}
type Ins_I32Or struct{}
type Ins_I32Xor struct{}
type Ins_I32Shl struct{}
type Ins_I32ShrS struct{}
type Ins_I32ShrU struct{}
type Ins_I32Rotl struct{}
type Ins_I32Rotr struct{}

type Ins_I64Clz struct{}
type Ins_I64Ctz struct{}
type Ins_I64Popcnt struct{}
type Ins_I64Add struct{}
type Ins_I64Sub struct{}
type Ins_I64Mul struct{}
type Ins_I64DivS struct{}
type Ins_I64DivU struct{}
type Ins_I64RemS struct{}
type Ins_I64RemU struct{}
type Ins_I64And struct{}
type Ins_I64Or struct{}
type Ins_I64Xor struct{}
type Ins_I64Shl struct{}
type Ins_I64ShrS struct{}
type Ins_I64ShrU struct{}
type Ins_I64Rotl struct{}
type Ins_I64Rotr struct{}

type Ins_F32Abs struct{}
type Ins_F32Neg struct{}
type Ins_F32Ceil struct{}
type Ins_F32Floor struct{}
type Ins_F32Trunc struct{}
type Ins_F32Nearest struct{}
type Ins_F32Sqrt struct{}
type Ins_F32Add struct{}
type Ins_F32Sub struct{}
type Ins_F32Mul struct{}
type Ins_F32Div struct{}
type Ins_F32Min struct{}
type Ins_F32Max struct{}
type Ins_F32Copysign struct{}

type Ins_F64Abs struct{}
type Ins_F64Neg struct{}
type Ins_F64Ceil struct{}
type Ins_F64Floor struct{}
type Ins_F64Trunc struct{}
type Ins_F64Nearest struct{}
type Ins_F64Sqrt struct{}
type Ins_F64Add struct{}
type Ins_F64Sub struct{}
type Ins_F64Mul struct{}
type Ins_F64Div struct{}
type Ins_F64Min struct{}
type Ins_F64Max struct{}
type Ins_F64Copysign struct{}

type Ins_I32WrapI64 struct{}
type Ins_I32TruncF32S struct{}
type Ins_I32TruncF32U struct{}
type Ins_I32TruncF64S struct{}
type Ins_I32TruncF64U struct{}
type Ins_I64ExtendI32S struct{}
type Ins_I64ExtendI32U struct{}
type Ins_I64TruncF32S struct{}
type Ins_I64TruncF32U struct{}
type Ins_I64TruncF64S struct{}
type Ins_I64TruncF64U struct{}
type Ins_F32ConvertI32S struct{}
type Ins_F32ConvertI32U struct{}
type Ins_F32ConvertI64S struct{}
type Ins_F32ConvertI64U struct{}
type Ins_F32DemoteF64 struct{}
type Ins_F64ConvertI32S struct{}
type Ins_F64ConvertI32U struct{}
type Ins_F64ConvertI64S struct{}
type Ins_F64ConvertI64U struct{}
type Ins_F64DemoteF32 struct{}
type Ins_I32ReintepretF32 struct{}
type Ins_I64ReintepretF64 struct{}
type Ins_I32ReintepretI32 struct{}
type Ins_I64ReintepretI64 struct{}

func (*Ins_I32Const) aInstruction() {}
func (*Ins_I32Store) aInstruction() {}
func (*Ins_I32Load) aInstruction()  {}

func (*Ins_I64Const) aInstruction() {}
func (*Ins_I64Store) aInstruction() {}
func (*Ins_I64Load) aInstruction()  {}

func (*Ins_Br) aInstruction()           {}
func (*Ins_Call) aInstruction()         {}
func (*Ins_CallIndirect) aInstruction() {}
func (*Ins_Drop) aInstruction()         {}
func (*Ins_Return) aInstruction()       {}
func (*Ins_Unreachable) aInstruction()  {}
func (*Ins_If) aInstruction()           {}

func (*Ins_GlobalGet) aInstruction() {}
func (*Ins_GlobalSet) aInstruction() {}

func (*Ins_LocalGet) aInstruction() {}
func (*Ins_LocalSet) aInstruction() {}
func (*Ins_LocalTee) aInstruction() {}

func (*Ins_I32Add) aInstruction()  {}
func (*Ins_I32Sub) aInstruction()  {}
func (*Ins_I32Mul) aInstruction()  {}
func (*Ins_I32DivS) aInstruction() {}
func (*Ins_I32DivU) aInstruction() {}
func (*Ins_I32RemS) aInstruction() {}
func (*Ins_I32RemU) aInstruction() {}

func (*Ins_I32Eqz) aInstruction() {}
func (*Ins_I32Eq) aInstruction()  {}
func (*Ins_I32Ne) aInstruction()  {}
func (*Ins_I32LtS) aInstruction() {}
func (*Ins_I32LtU) aInstruction() {}
func (*Ins_I32GtS) aInstruction() {}
func (*Ins_I32GtU) aInstruction() {}
func (*Ins_I32LeS) aInstruction() {}
func (*Ins_I32LeU) aInstruction() {}
func (*Ins_I32GeS) aInstruction() {}
func (*Ins_I32GeU) aInstruction() {}
