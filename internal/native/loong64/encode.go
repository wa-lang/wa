package loong64

import (
	"fmt"

	"wa-lang.org/wa/internal/native/abi"
)

// 返回寄存器机器码编号
func RegI(r abi.RegType) uint32 {
	return (*_OpContextType)(nil).regI(r)
}

// 返回浮点数寄存器机器码编号
func RegF(r abi.RegType) uint32 {
	return (*_OpContextType)(nil).regF(r)
}

// 编码龙芯指令
func Encode(cpu abi.CPUType, as abi.As, arg *abi.AsArgument) (uint32, error) {
	switch cpu {
	case abi.RISCV64:
		return EncodeLA64(as, arg)
	default:
		return 0, fmt.Errorf("unknonw cpu: %v", cpu)
	}
}

// 编码龙芯64指令
func EncodeLA64(as abi.As, arg *abi.AsArgument) (uint32, error) {
	if as <= 0 || as >= ALAST {
		return 0, fmt.Errorf("loong64.EncodeLA64: bad as(%v), arg(%v)", as, arg)
	}

	ctx := _AOpContextTable[as]
	assert(ctx.mask != 0)

	return ctx.encodeRaw(as, arg)
}

func (ctx *_OpContextType) encodeRaw(as abi.As, arg *abi.AsArgument) (x uint32, err error) {
	assert(ctx != nil)
	assert(ctx.op == as)
	assert(arg != nil)

	x = ctx.mask & ctx.value

	switch ctx.fmt {
	case OpFormatType_NULL:
		return
	case OpFormatType_2R:
		panic("TODO")
	case OpFormatType_2F:
		panic("TODO")
	case OpFormatType_1F_1R:
		panic("TODO")
	case OpFormatType_1R_1F:
		panic("TODO")
	case OpFormatType_3R:
		panic("TODO")
	case OpFormatType_3F:
		panic("TODO")
	case OpFormatType_1F_2R:
		panic("TODO")
	case OpFormatType_4F:
		panic("TODO")
	case OpFormatType_2R_ui5:
		panic("TODO")
	case OpFormatType_2R_ui6:
		panic("TODO")
	case OpFormatType_2R_si12:
		panic("TODO")
	case OpFormatType_2R_ui12:
		panic("TODO")
	case OpFormatType_2R_si14:
		panic("TODO")
	case OpFormatType_2R_si16:
		panic("TODO")
	case OpFormatType_1R_si20:
		panic("TODO")
	case OpFormatType_0_2R:
		panic("TODO")
	case OpFormatType_3R_sa2:
		panic("TODO")
	case OpFormatType_3R_sa3:
		panic("TODO")
	case OpFormatType_code:
		panic("TODO")
	case OpFormatType_code_1R_si12:
		panic("TODO")
	case OpFormatType_2R_msbw_lsbw:
		panic("TODO")
	case OpFormatType_2R_msbd_lsbd:
		panic("TODO")
	case OpFormatType_fcsr_1R:
		panic("TODO")
	case OpFormatType_1R_fcsr:
		panic("TODO")
	case OpFormatType_cd_1R:
		panic("TODO")
	case OpFormatType_cd_1F:
		panic("TODO")
	case OpFormatType_cd_2R:
		panic("TODO")
	case OpFormatType_cd_2F:
		panic("TODO")
	case OpFormatType_1R_cj:
		panic("TODO")
	case OpFormatType_1F_cj:
		panic("TODO")
	case OpFormatType_1R_csr:
		panic("TODO")
	case OpFormatType_2R_csr:
		panic("TODO")
	case OpFormatType_2R_level:
		panic("TODO")
	case OpFormatType_level:
		panic("TODO")
	case OpFormatType_0_1R_seq:
		panic("TODO")
	case OpFormatType_op_2R:
		panic("TODO")
	case OpFormatType_3R_ca:
		panic("TODO")
	case OpFormatType_hint_1R_si12:
		panic("TODO")
	case OpFormatType_hint_2R:
		panic("TODO")
	case OpFormatType_hint:
		panic("TODO")
	case OpFormatType_cj_offset:
		panic("TODO")
	case OpFormatType_rj_offset:
		panic("TODO")
	case OpFormatType_rj_rd_offset:
		panic("TODO")
	case OpFormatType_rd_rj_offset:
		panic("TODO")
	case OpFormatType_offset:
		panic("TODO")
	default:
		panic("unreachable")
	}

	return x, nil
}
