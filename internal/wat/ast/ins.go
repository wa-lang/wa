// 版权 @2024 凹语言 作者。保留所有权利。

package ast

import "wa-lang.org/wa/internal/wat/token"

// 指令对应接口
type Instruction interface {
	Token() token.Token

	aInstruction()
}

type OpToken token.Token

func (OpToken) aInstruction() {}

func (tok OpToken) Token() token.Token {
	return token.Token(tok)
}
func (tok OpToken) Valid() bool {
	return token.Token(tok).IsIsntruction()
}
func (tok OpToken) String() string {
	return token.Token(tok).String()
}

type Ins_Unreachable struct{ OpToken }
type Ins_Nop struct{ OpToken }
type Ins_Block struct {
	OpToken
	Label   string
	Results []token.Token // 返回值类型: I32, I64, F32, F64
	List    []Instruction
}
type Ins_Loop struct {
	OpToken
	Label   string
	Results []token.Token // 返回值类型: I32, I64, F32, F64
	List    []Instruction
}
type Ins_If struct {
	OpToken
	Label   string
	Results []token.Token // 返回值类型: I32, I64, F32, F64
	Body    []Instruction
	Else    []Instruction
}
type Ins_Else struct{ OpToken }
type Ins_End struct{ OpToken }
type Ins_Br struct {
	OpToken
	X string // todo: IdxOrLabel
}
type Ins_BrIf struct {
	OpToken
	X string
}
type Ins_BrTable struct {
	OpToken
	XList []string
}
type Ins_Return struct{ OpToken }

type Ins_Call struct {
	OpToken
	X string
}

type Ins_CallIndirect struct {
	OpToken
	TableIdx string
	TypeIdx  string
}
type Ins_Drop struct{ OpToken }
type Ins_Select struct {
	OpToken
	ResultTyp token.Token // I32, I64, F32, F64
}

type Ins_LocalGet struct {
	OpToken
	X string
}
type Ins_LocalSet struct {
	OpToken
	X string
}
type Ins_LocalTee struct {
	OpToken
	X string
}
type Ins_GlobalGet struct {
	OpToken
	X string
}
type Ins_GlobalSet struct {
	OpToken
	X string
}
type Ins_TableGet struct {
	OpToken
	TableIdx string
}
type Ins_TableSet struct {
	OpToken
	TableIdx string
}
type Ins_I32Load struct {
	OpToken
	Align  uint
	Offset uint
}
type Ins_I64Load struct {
	OpToken
	Align  uint
	Offset uint
}
type Ins_F32Load struct {
	OpToken
	Align  uint
	Offset uint
}
type Ins_F64Load struct {
	OpToken
	Align  uint
	Offset uint
}
type Ins_I32Load8S struct {
	OpToken
	Align  uint
	Offset uint
}
type Ins_I32Load8U struct {
	OpToken
	Align  uint
	Offset uint
}
type Ins_I32Load16S struct {
	OpToken
	Align  uint
	Offset uint
}
type Ins_I32Load16U struct {
	OpToken
	Align  uint
	Offset uint
}
type Ins_I64Load8S struct {
	OpToken
	Align  uint
	Offset uint
}
type Ins_I64Load8U struct {
	OpToken
	Align  uint
	Offset uint
}
type Ins_I64Load16S struct {
	OpToken
	Align  uint
	Offset uint
}
type Ins_I64Load16U struct {
	OpToken
	Align  uint
	Offset uint
}
type Ins_I64Load32S struct {
	OpToken
	Align  uint
	Offset uint
}
type Ins_I64Load32U struct {
	OpToken
	Align  uint
	Offset uint
}

type Ins_I32Store struct {
	OpToken
	Align  uint
	Offset uint
}

type Ins_I64Store struct {
	OpToken
	Align  uint
	Offset uint
}
type Ins_F32Store struct {
	OpToken
	Align  uint
	Offset uint
}
type Ins_F64Store struct {
	OpToken
	Align  uint
	Offset uint
}

type Ins_I32Store8 struct {
	OpToken
	Align  uint
	Offset uint
}
type Ins_I32Store16 struct {
	OpToken
	Align  uint
	Offset uint
}
type Ins_I64Store8 struct {
	OpToken
	Align  uint
	Offset uint
}
type Ins_I64Store16 struct {
	OpToken
	Align  uint
	Offset uint
}
type Ins_I64Store32 struct {
	OpToken
	Align  uint
	Offset uint
}

type Ins_MemorySize struct{ OpToken }
type Ins_MemoryGrow struct{ OpToken }

type Ins_MemoryInit struct {
	OpToken
	DataIdx int32
}
type Ins_MemoryCopy struct{ OpToken }
type Ins_MemoryFill struct{ OpToken }

type Ins_I32Const struct {
	OpToken
	X int32
}
type Ins_I64Const struct {
	OpToken
	X int64
}
type Ins_F32Const struct {
	OpToken
	X float32
}
type Ins_F64Const struct {
	OpToken
	X float64
}

type Ins_I32Eqz struct{ OpToken }
type Ins_I32Eq struct{ OpToken }
type Ins_I32Ne struct{ OpToken }
type Ins_I32LtS struct{ OpToken }
type Ins_I32LtU struct{ OpToken }
type Ins_I32GtS struct{ OpToken }
type Ins_I32GtU struct{ OpToken }
type Ins_I32LeS struct{ OpToken }
type Ins_I32LeU struct{ OpToken }
type Ins_I32GeS struct{ OpToken }
type Ins_I32GeU struct{ OpToken }

type Ins_I64Eqz struct{ OpToken }
type Ins_I64Eq struct{ OpToken }
type Ins_I64Ne struct{ OpToken }
type Ins_I64LtS struct{ OpToken }
type Ins_I64LtU struct{ OpToken }
type Ins_I64GtS struct{ OpToken }
type Ins_I64GtU struct{ OpToken }
type Ins_I64LeS struct{ OpToken }
type Ins_I64LeU struct{ OpToken }
type Ins_I64GeS struct{ OpToken }
type Ins_I64GeU struct{ OpToken }

type Ins_F32Eq struct{ OpToken }
type Ins_F32Ne struct{ OpToken }
type Ins_F32Lt struct{ OpToken }
type Ins_F32Gt struct{ OpToken }
type Ins_F32Le struct{ OpToken }
type Ins_F32Ge struct{ OpToken }

type Ins_F64Eq struct{ OpToken }
type Ins_F64Ne struct{ OpToken }
type Ins_F64Lt struct{ OpToken }
type Ins_F64Gt struct{ OpToken }
type Ins_F64Le struct{ OpToken }
type Ins_F64Ge struct{ OpToken }

type Ins_I32Clz struct{ OpToken }
type Ins_I32Ctz struct{ OpToken }
type Ins_I32Popcnt struct{ OpToken }
type Ins_I32Add struct{ OpToken }
type Ins_I32Sub struct{ OpToken }
type Ins_I32Mul struct{ OpToken }
type Ins_I32DivS struct{ OpToken }
type Ins_I32DivU struct{ OpToken }
type Ins_I32RemS struct{ OpToken }
type Ins_I32RemU struct{ OpToken }
type Ins_I32And struct{ OpToken }
type Ins_I32Or struct{ OpToken }
type Ins_I32Xor struct{ OpToken }
type Ins_I32Shl struct{ OpToken }
type Ins_I32ShrS struct{ OpToken }
type Ins_I32ShrU struct{ OpToken }
type Ins_I32Rotl struct{ OpToken }
type Ins_I32Rotr struct{ OpToken }

type Ins_I64Clz struct{ OpToken }
type Ins_I64Ctz struct{ OpToken }
type Ins_I64Popcnt struct{ OpToken }
type Ins_I64Add struct{ OpToken }
type Ins_I64Sub struct{ OpToken }
type Ins_I64Mul struct{ OpToken }
type Ins_I64DivS struct{ OpToken }
type Ins_I64DivU struct{ OpToken }
type Ins_I64RemS struct{ OpToken }
type Ins_I64RemU struct{ OpToken }
type Ins_I64And struct{ OpToken }
type Ins_I64Or struct{ OpToken }
type Ins_I64Xor struct{ OpToken }
type Ins_I64Shl struct{ OpToken }
type Ins_I64ShrS struct{ OpToken }
type Ins_I64ShrU struct{ OpToken }
type Ins_I64Rotl struct{ OpToken }
type Ins_I64Rotr struct{ OpToken }

type Ins_F32Abs struct{ OpToken }
type Ins_F32Neg struct{ OpToken }
type Ins_F32Ceil struct{ OpToken }
type Ins_F32Floor struct{ OpToken }
type Ins_F32Trunc struct{ OpToken }
type Ins_F32Nearest struct{ OpToken }
type Ins_F32Sqrt struct{ OpToken }
type Ins_F32Add struct{ OpToken }
type Ins_F32Sub struct{ OpToken }
type Ins_F32Mul struct{ OpToken }
type Ins_F32Div struct{ OpToken }
type Ins_F32Min struct{ OpToken }
type Ins_F32Max struct{ OpToken }
type Ins_F32Copysign struct{ OpToken }

type Ins_F64Abs struct{ OpToken }
type Ins_F64Neg struct{ OpToken }
type Ins_F64Ceil struct{ OpToken }
type Ins_F64Floor struct{ OpToken }
type Ins_F64Trunc struct{ OpToken }
type Ins_F64Nearest struct{ OpToken }
type Ins_F64Sqrt struct{ OpToken }
type Ins_F64Add struct{ OpToken }
type Ins_F64Sub struct{ OpToken }
type Ins_F64Mul struct{ OpToken }
type Ins_F64Div struct{ OpToken }
type Ins_F64Min struct{ OpToken }
type Ins_F64Max struct{ OpToken }
type Ins_F64Copysign struct{ OpToken }

type Ins_I32WrapI64 struct{ OpToken }
type Ins_I32TruncF32S struct{ OpToken }
type Ins_I32TruncF32U struct{ OpToken }
type Ins_I32TruncF64S struct{ OpToken }
type Ins_I32TruncF64U struct{ OpToken }
type Ins_I64ExtendI32S struct{ OpToken }
type Ins_I64ExtendI32U struct{ OpToken }
type Ins_I64TruncF32S struct{ OpToken }
type Ins_I64TruncF32U struct{ OpToken }
type Ins_I64TruncF64S struct{ OpToken }
type Ins_I64TruncF64U struct{ OpToken }
type Ins_F32ConvertI32S struct{ OpToken }
type Ins_F32ConvertI32U struct{ OpToken }
type Ins_F32ConvertI64S struct{ OpToken }
type Ins_F32ConvertI64U struct{ OpToken }
type Ins_F32DemoteF64 struct{ OpToken }
type Ins_F64ConvertI32S struct{ OpToken }
type Ins_F64ConvertI32U struct{ OpToken }
type Ins_F64ConvertI64S struct{ OpToken }
type Ins_F64ConvertI64U struct{ OpToken }
type Ins_F64PromoteF32 struct{ OpToken }
type Ins_I32ReintepretF32 struct{ OpToken }
type Ins_I64ReintepretF64 struct{ OpToken }
type Ins_F32ReintepretI32 struct{ OpToken }
type Ins_F64ReintepretI64 struct{ OpToken }
