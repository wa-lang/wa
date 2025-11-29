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
