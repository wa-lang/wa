package x64

// X64 汇编指令参数类型
type OpFormatType int

const (
	OpFormatType_NULL    OpFormatType = iota
	OpFormatType_NoArgs               // 无参数
	OpFormatType_Imm                  // 一元立即数
	OpFormatType_Reg                  // 一元寄存器
	OpFormatType_Mem                  // 一元内  存
	OpFormatType_Imm2Reg              // 立即数 => 寄存器
	OpFormatType_Imm2Mem              // 立即数 => 内  存
	OpFormatType_Reg2Reg              // 寄存器 => 寄存器
	OpFormatType_Mem2Reg              // 内  存 => 寄存器
	OpFormatType_Reg2Mem              // 寄存器 => 内  存
	OpFormatType_Any                  // 一元未知
	OpFormatType_Any2Any              // 二元未知, 对应 mov 指令
)

// X64 指令类型
var x64ModeTable = [...]OpFormatType{
	AADD:       OpFormatType_Any2Any,
	AADDSD:     OpFormatType_Any2Any,
	AADDSS:     OpFormatType_Any2Any,
	AAND:       OpFormatType_Any2Any,
	ACALL:      OpFormatType_Any,
	ACDQ:       OpFormatType_NoArgs,
	ACMP:       OpFormatType_Any2Any,
	ACVTSI2SD:  OpFormatType_Reg2Reg,
	ACVTSS2SD:  OpFormatType_Reg2Reg,
	ACVTTSD2SI: OpFormatType_Reg2Reg,
	ADEC:       OpFormatType_Any,
	ADIV:       OpFormatType_Any,
	ADIVSD:     OpFormatType_Any2Any,
	ADIVSS:     OpFormatType_Any2Any,
	AIDIV:      OpFormatType_Any,
	AIMUL:      OpFormatType_Any2Any,
	AINC:       OpFormatType_Any,
	AJA:        OpFormatType_Imm,
	AJB:        OpFormatType_Any,
	AJE:        OpFormatType_Imm,
	AJGE:       OpFormatType_Any,
	AJMP:       OpFormatType_Any,
	AJNS:       OpFormatType_Any,
	AJNZ:       OpFormatType_Any,
	AJZ:        OpFormatType_Any,
	ALEA:       OpFormatType_Mem2Reg,
	AMOV:       OpFormatType_Any2Any,
	AMOVABS:    OpFormatType_Imm2Reg,
	AMOVQ:      OpFormatType_Any2Any,
	AMOVSD:     OpFormatType_Any2Any,
	AMOVSS:     OpFormatType_Any2Any,
	AMOVSXD:    OpFormatType_Any2Any,
	AMOVZX:     OpFormatType_Any2Any,
	AMULSD:     OpFormatType_Any2Any,
	AMULSS:     OpFormatType_Any2Any,
	ANEG:       OpFormatType_Any,
	ANOP:       OpFormatType_NoArgs,
	AOR:        OpFormatType_Any2Any,
	APOP:       OpFormatType_Any,
	APUSH:      OpFormatType_Any,
	ARET:       OpFormatType_NoArgs,
	ASAR:       OpFormatType_Reg2Reg,
	ASETA:      OpFormatType_Any,
	ASETAE:     OpFormatType_Any,
	ASETB:      OpFormatType_Any,
	ASETBE:     OpFormatType_Any,
	ASETE:      OpFormatType_Any,
	ASETG:      OpFormatType_Any,
	ASETGE:     OpFormatType_Any,
	ASETL:      OpFormatType_Any,
	ASETLE:     OpFormatType_Any,
	ASETNE:     OpFormatType_Any,
	ASETNP:     OpFormatType_Any,
	ASHL:       OpFormatType_Any2Any,
	ASTD:       OpFormatType_NoArgs,
	ASUB:       OpFormatType_Any2Any,
	ASUBSD:     OpFormatType_Any2Any,
	ASUBSS:     OpFormatType_Any2Any,
	ASYSCALL:   OpFormatType_NoArgs,
	ATEST:      OpFormatType_Any2Any,
	AUCOMISD:   OpFormatType_Reg2Reg,
	AXOR:       OpFormatType_Reg2Reg,

	ALAST: OpFormatType_NULL,
}
