// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

// 重要: 请保持 RV32 和 RV64 版本完全一致, LAUInt 在另一个文件定义.

package loong64

import (
	"fmt"

	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/loong64"
	"wa-lang.org/wa/internal/native/wemu/device"
)

// 处理器
type CPU struct {
	RegX [32]LAUInt  // 整数寄存器
	RegF [32]float64 // 浮点数寄存器
	PC   LAUInt      // PC指针
}

// 构建处理器
func NewCPU() *CPU {
	return new(CPU)
}

func (p *CPU) GetPC() uint64  { return uint64(p.PC) }
func (p *CPU) SetPC(v uint64) { p.PC = LAUInt(v) }

func (p *CPU) XRegNum() int            { return len(p.RegX) }
func (p *CPU) GetXReg(i int) uint64    { return uint64(p.RegX[i]) }
func (p *CPU) SetXReg(i int, v uint64) { p.RegX[i] = LAUInt(v) }

func (p *CPU) FRegNum() int             { return len(p.RegF) }
func (p *CPU) GetFReg(i int) float64    { return p.RegF[i] }
func (p *CPU) SetFReg(i int, v float64) { p.RegF[i] = v }

// 重置CPU
func (p *CPU) Reset(pc, sp uint64) {
	*p = CPU{}
	p.RegX[loong64.RegI(loong64.REG_SP)] = LAUInt(sp)
	p.PC = LAUInt(pc)
}

// 单步执行
func (p *CPU) StepRun(bus *device.Bus) error {
	inst, err := bus.Read(uint64(p.PC), 4)
	if err != nil {
		return err
	}
	as, arg, argRaw, err := loong64.DecodeEx(uint32(inst))
	if err != nil {
		return fmt.Errorf("fetch instruntion failed at 0x%08X: %v", p.PC, err)
	}

	if err := p.execInst(bus, as, argRaw); err != nil {
		return fmt.Errorf("exec instruntion(%08d:%s) failed at 0x%08X: %v", inst,
			loong64.AsmSyntax(as, "", arg),
			p.PC,
			err,
		)
	}

	return nil
}

func (p *CPU) execInst(bus *device.Bus, as abi.As, arg *abi.AsRawArgument) error {
	// 重置0寄存器
	p.RegX[0] = 0
	p.RegF[0] = 0

	// 当前的PC
	curPC := p.PC
	_ = curPC

	// 调整PC(少数跳转指令会覆盖)
	p.PC += 4

	// 执行指令
	// 基于指令编码分类
	switch loong64.AsFormatType(as) {
	case loong64.OpFormatType_NULL:
		switch as {
		case loong64.AERTN:
			return fmt.Errorf("privileged instruction encountered: %s", loong64.AsString(as, ""))
		case loong64.ATLBCLR:
			return fmt.Errorf("privileged instruction encountered: %s", loong64.AsString(as, ""))
		case loong64.ATLBFILL:
			return fmt.Errorf("privileged instruction encountered: %s", loong64.AsString(as, ""))
		case loong64.ATLBFLUSH:
			return fmt.Errorf("privileged instruction encountered: %s", loong64.AsString(as, ""))
		case loong64.ATLBRD:
			return fmt.Errorf("privileged instruction encountered: %s", loong64.AsString(as, ""))
		case loong64.ATLBSRCH:
			return fmt.Errorf("privileged instruction encountered: %s", loong64.AsString(as, ""))
		case loong64.ATLBWR:
			return fmt.Errorf("privileged instruction encountered: %s", loong64.AsString(as, ""))
		default:
			panic("unreachable")
		}
	case loong64.OpFormatType_2R:
		switch as {
		case loong64.ABITREV_4B:
			panic("TODO")
		case loong64.ABITREV_8B:
			panic("TODO")
		case loong64.ABITREV_D:
			panic("TODO")
		case loong64.ABITREV_W:
			panic("TODO")
		case loong64.ACLO_D:
			panic("TODO")
		case loong64.ACLO_W:
			panic("TODO")
		case loong64.ACLZ_D:
			panic("TODO")
		case loong64.ACLZ_W:
			panic("TODO")
		case loong64.ACPUCFG:
			panic("TODO")
		case loong64.ACTO_D:
			panic("TODO")
		case loong64.ACTO_W:
			panic("TODO")
		case loong64.ACTZ_D:
			panic("TODO")
		case loong64.ACTZ_W:
			panic("TODO")
		case loong64.AEXT_W_B:
			panic("TODO")
		case loong64.AEXT_W_H:
			panic("TODO")
		case loong64.AIOCSRRD_B:
			panic("TODO")
		case loong64.AIOCSRRD_D:
			panic("TODO")
		case loong64.AIOCSRRD_H:
			panic("TODO")
		case loong64.AIOCSRRD_W:
			panic("TODO")
		case loong64.AIOCSRWR_B:
			panic("TODO")
		case loong64.AIOCSRWR_D:
			panic("TODO")
		case loong64.AIOCSRWR_H:
			panic("TODO")
		case loong64.AIOCSRWR_W:
			panic("TODO")
		case loong64.ALLACQ_D:
			panic("TODO")
		case loong64.ALLACQ_W:
			panic("TODO")
		case loong64.ARDTIMEH_W:
			panic("TODO")
		case loong64.ARDTIMEL_W:
			panic("TODO")
		case loong64.ARDTIME_D:
			panic("TODO")
		case loong64.AREVB_2H:
			panic("TODO")
		case loong64.AREVB_2W:
			panic("TODO")
		case loong64.AREVB_4H:
			panic("TODO")
		case loong64.AREVB_D:
			panic("TODO")
		case loong64.AREVH_2W:
			panic("TODO")
		case loong64.AREVH_D:
			panic("TODO")
		case loong64.ASCREL_D:
			panic("TODO")
		case loong64.ASCREL_W:
			panic("TODO")
		default:
			panic("unreachable")
		}
	case loong64.OpFormatType_2F:
		switch as {
		default:
			return fmt.Errorf("unsupport: %s", loong64.AsString(as, ""))
		case loong64.AFABS_D:
			panic("TODO")
		case loong64.AFABS_S:
			panic("TODO")
		case loong64.AFCLASS_D:
			panic("TODO")
		case loong64.AFCLASS_S:
			panic("TODO")
		case loong64.AFCVT_D_S:
			panic("TODO")
		case loong64.AFCVT_S_D:
			panic("TODO")
		case loong64.AFFINT_D_L:
			panic("TODO")
		case loong64.AFFINT_D_W:
			panic("TODO")
		case loong64.AFFINT_S_L:
			panic("TODO")
		case loong64.AFFINT_S_W:
			panic("TODO")
		case loong64.AFLOGB_D:
			panic("TODO")
		case loong64.AFLOGB_S:
			panic("TODO")
		case loong64.AFMOV_D:
			panic("TODO")
		case loong64.AFMOV_S:
			panic("TODO")
		case loong64.AFNEG_D:
			panic("TODO")
		case loong64.AFNEG_S:
			panic("TODO")
		case loong64.AFRECIPE_D:
			panic("TODO")
		case loong64.AFRECIPE_S:
			panic("TODO")
		case loong64.AFRECIP_D:
			panic("TODO")
		case loong64.AFRECIP_S:
			panic("TODO")
		case loong64.AFRINT_D:
			panic("TODO")
		case loong64.AFRINT_S:
			panic("TODO")
		case loong64.AFRSQRTE_D:
			panic("TODO")
		case loong64.AFRSQRTE_S:
			panic("TODO")
		case loong64.AFRSQRT_D:
			panic("TODO")
		case loong64.AFRSQRT_S:
			panic("TODO")
		case loong64.AFSQRT_D:
			panic("TODO")
		case loong64.AFSQRT_S:
			panic("TODO")
		case loong64.AFTINTRM_L_D:
			panic("TODO")
		case loong64.AFTINTRM_L_S:
			panic("TODO")
		case loong64.AFTINTRM_W_D:
			panic("TODO")
		case loong64.AFTINTRM_W_S:
			panic("TODO")
		case loong64.AFTINTRNE_L_D:
			panic("TODO")
		case loong64.AFTINTRNE_L_S:
			panic("TODO")
		case loong64.AFTINTRNE_W_D:
			panic("TODO")
		case loong64.AFTINTRNE_W_S:
			panic("TODO")
		case loong64.AFTINTRP_L_D:
			panic("TODO")
		case loong64.AFTINTRP_L_S:
			panic("TODO")
		case loong64.AFTINTRP_W_D:
			panic("TODO")
		case loong64.AFTINTRP_W_S:
			panic("TODO")
		case loong64.AFTINTRZ_L_D:
			panic("TODO")
		case loong64.AFTINTRZ_L_S:
			panic("TODO")
		case loong64.AFTINTRZ_W_D:
			panic("TODO")
		case loong64.AFTINTRZ_W_S:
			panic("TODO")
		case loong64.AFTINT_L_D:
			panic("TODO")
		case loong64.AFTINT_L_S:
			panic("TODO")
		case loong64.AFTINT_W_D:
			panic("TODO")
		case loong64.AFTINT_W_S:
			panic("TODO")
		}
	case loong64.OpFormatType_1F_1R:
		switch as {
		case loong64.AMOVGR2FRH_W:
			panic("TODO")
		case loong64.AMOVGR2FR_D:
			panic("TODO")
		case loong64.AMOVGR2FR_W:
			panic("TODO")
		default:
			return fmt.Errorf("unsupport: %s", loong64.AsString(as, ""))
		}
	case loong64.OpFormatType_1R_1F:
		switch as {
		case loong64.AMOVFR2GR_D:
			panic("TODO")
		case loong64.AMOVFR2GR_S:
			panic("TODO")
		case loong64.AMOVFRH2GR_S:
			panic("TODO")
		default:
			return fmt.Errorf("unsupport: %s", loong64.AsString(as, ""))
		}
	case loong64.OpFormatType_3R:
		switch as {
		case loong64.AADD_D:
			p.RegX[arg.Rd] = p.RegX[arg.Rs1] + p.RegX[arg.Rs2]
			return nil
		case loong64.AADD_W:
			result := int32(p.RegX[arg.Rs1]) + int32(p.RegX[arg.Rs2])
			p.RegX[arg.Rd] = LAUInt(int64(result))
			return nil
		case loong64.AAMADD_B:
			panic("TODO")
		case loong64.AAMADD_D:
			panic("TODO")
		case loong64.AAMADD_DB_B:
			panic("TODO")
		case loong64.AAMADD_DB_D:
			panic("TODO")
		case loong64.AAMADD_DB_H:
			panic("TODO")
		case loong64.AAMADD_DB_W:
			panic("TODO")
		case loong64.AAMADD_H:
			panic("TODO")
		case loong64.AAMADD_W:
			panic("TODO")
		case loong64.AAMAND_D:
			panic("TODO")
		case loong64.AAMAND_DB_D:
			panic("TODO")
		case loong64.AAMAND_DB_W:
			panic("TODO")
		case loong64.AAMAND_W:
			panic("TODO")
		case loong64.AAMCAS_B:
			panic("TODO")
		case loong64.AAMCAS_D:
			panic("TODO")
		case loong64.AAMCAS_DB_B:
			panic("TODO")
		case loong64.AAMCAS_DB_D:
			panic("TODO")
		case loong64.AAMCAS_DB_H:
			panic("TODO")
		case loong64.AAMCAS_DB_W:
			panic("TODO")
		case loong64.AAMCAS_H:
			panic("TODO")
		case loong64.AAMCAS_W:
			panic("TODO")
		case loong64.AAMMAX_D:
			panic("TODO")
		case loong64.AAMMAX_DB_D:
			panic("TODO")
		case loong64.AAMMAX_DB_DU:
			panic("TODO")
		case loong64.AAMMAX_DB_W:
			panic("TODO")
		case loong64.AAMMAX_DB_WU:
			panic("TODO")
		case loong64.AAMMAX_DU:
			panic("TODO")
		case loong64.AAMMAX_W:
			panic("TODO")
		case loong64.AAMMAX_WU:
			panic("TODO")
		case loong64.AAMMIN_D:
			panic("TODO")
		case loong64.AAMMIN_DB_D:
			panic("TODO")
		case loong64.AAMMIN_DB_DU:
			panic("TODO")
		case loong64.AAMMIN_DB_W:
			panic("TODO")
		case loong64.AAMMIN_DB_WU:
			panic("TODO")
		case loong64.AAMMIN_DU:
			panic("TODO")
		case loong64.AAMMIN_W:
			panic("TODO")
		case loong64.AAMMIN_WU:
			panic("TODO")
		case loong64.AAMOR_D:
			panic("TODO")
		case loong64.AAMOR_DB_D:
			panic("TODO")
		case loong64.AAMOR_DB_W:
			panic("TODO")
		case loong64.AAMOR_W:
			panic("TODO")
		case loong64.AAMSWAP_B:
			panic("TODO")
		case loong64.AAMSWAP_D:
			panic("TODO")
		case loong64.AAMSWAP_DB_B:
			panic("TODO")
		case loong64.AAMSWAP_DB_D:
			panic("TODO")
		case loong64.AAMSWAP_DB_H:
			panic("TODO")
		case loong64.AAMSWAP_DB_W:
			panic("TODO")
		case loong64.AAMSWAP_H:
			panic("TODO")
		case loong64.AAMSWAP_W:
			panic("TODO")
		case loong64.AAMXOR_D:
			panic("TODO")
		case loong64.AAMXOR_DB_D:
			panic("TODO")
		case loong64.AAMXOR_DB_W:
			panic("TODO")
		case loong64.AAMXOR_W:
			panic("TODO")
		case loong64.AAND:
			p.RegX[arg.Rd] = p.RegX[arg.Rs1] & p.RegX[arg.Rs2]
			return nil
		case loong64.AANDN:
			panic("TODO")
		case loong64.ACRCC_W_B_W:
			panic("TODO")
		case loong64.ACRCC_W_D_W:
			panic("TODO")
		case loong64.ACRCC_W_H_W:
			panic("TODO")
		case loong64.ACRCC_W_W_W:
			panic("TODO")
		case loong64.ACRC_W_B_W:
			panic("TODO")
		case loong64.ACRC_W_D_W:
			panic("TODO")
		case loong64.ACRC_W_H_W:
			panic("TODO")
		case loong64.ACRC_W_W_W:
			panic("TODO")
		case loong64.ADIV_D:
			panic("TODO")
		case loong64.ADIV_DU:
			panic("TODO")
		case loong64.ADIV_W:
			panic("TODO")
		case loong64.ADIV_WU:
			panic("TODO")
		case loong64.ALDGT_B:
			panic("TODO")
		case loong64.ALDGT_D:
			panic("TODO")
		case loong64.ALDGT_H:
			panic("TODO")
		case loong64.ALDGT_W:
			panic("TODO")
		case loong64.ALDLE_B:
			panic("TODO")
		case loong64.ALDLE_D:
			panic("TODO")
		case loong64.ALDLE_H:
			panic("TODO")
		case loong64.ALDLE_W:
			panic("TODO")
		case loong64.ALDX_B:
			panic("TODO")
		case loong64.ALDX_BU:
			panic("TODO")
		case loong64.ALDX_D:
			panic("TODO")
		case loong64.ALDX_H:
			panic("TODO")
		case loong64.ALDX_HU:
			panic("TODO")
		case loong64.ALDX_W:
			panic("TODO")
		case loong64.ALDX_WU:
			panic("TODO")
		case loong64.AMASKEQZ:
			panic("TODO")
		case loong64.AMASKNEZ:
			panic("TODO")
		case loong64.AMOD_D:
			panic("TODO")
		case loong64.AMOD_DU:
			panic("TODO")
		case loong64.AMOD_W:
			panic("TODO")
		case loong64.AMOD_WU:
			panic("TODO")
		case loong64.AMULH_D:
			panic("TODO")
		case loong64.AMULH_DU:
			panic("TODO")
		case loong64.AMULH_W:
			panic("TODO")
		case loong64.AMULH_WU:
			panic("TODO")
		case loong64.AMULW_D_W:
			panic("TODO")
		case loong64.AMULW_D_WU:
			panic("TODO")
		case loong64.AMUL_D:
			panic("TODO")
		case loong64.AMUL_W:
			panic("TODO")
		case loong64.ANOR:
			panic("TODO")
		case loong64.AOR:
			p.RegX[arg.Rd] = p.RegX[arg.Rs1] | p.RegX[arg.Rs2]
			return nil
		case loong64.AORN:
			panic("TODO")
		case loong64.AROTR_D:
			panic("TODO")
		case loong64.AROTR_W:
			panic("TODO")
		case loong64.ASC_Q:
			panic("TODO")
		case loong64.ASLL_D:
			panic("TODO")
		case loong64.ASLL_W:
			panic("TODO")
		case loong64.ASLT:
			if int64(p.RegX[arg.Rs1]) < int64(p.RegX[arg.Rs2]) {
				p.RegX[arg.Rd] = 1
			} else {
				p.RegX[arg.Rd] = 0
			}
			return nil
		case loong64.ASLTU:
			panic("TODO")
		case loong64.ASRA_D:
			panic("TODO")
		case loong64.ASRA_W:
			panic("TODO")
		case loong64.ASRL_D:
			panic("TODO")
		case loong64.ASRL_W:
			panic("TODO")
		case loong64.ASTGT_B:
			panic("TODO")
		case loong64.ASTGT_D:
			panic("TODO")
		case loong64.ASTGT_H:
			panic("TODO")
		case loong64.ASTGT_W:
			panic("TODO")
		case loong64.ASTLE_B:
			panic("TODO")
		case loong64.ASTLE_D:
			panic("TODO")
		case loong64.ASTLE_H:
			panic("TODO")
		case loong64.ASTLE_W:
			panic("TODO")
		case loong64.ASTX_B:
			panic("TODO")
		case loong64.ASTX_D:
			panic("TODO")
		case loong64.ASTX_H:
			panic("TODO")
		case loong64.ASTX_W:
			panic("TODO")
		case loong64.ASUB_D:
			p.RegX[arg.Rd] = p.RegX[arg.Rs1] - p.RegX[arg.Rs2]
			return nil
		case loong64.ASUB_W:
			result := int32(p.RegX[arg.Rs1]) - int32(p.RegX[arg.Rs2])
			p.RegX[arg.Rd] = LAUInt(int64(result))
			return nil
		case loong64.AXOR:
			panic("TODO")

		default:
			return fmt.Errorf("unsupport: %s", loong64.AsString(as, ""))
		}
	case loong64.OpFormatType_3F:
		switch as {
		default:
			return fmt.Errorf("unsupport: %s", loong64.AsString(as, ""))
		case loong64.AFADD_D:
			panic("TODO")
		case loong64.AFADD_S:
			fs1 := float32(p.RegF[arg.Rs1])
			fs2 := float32(p.RegF[arg.Rs2])
			p.RegF[arg.Rd] = float64(fs1 + fs2)
			return nil
		case loong64.AFCOPYSIGN_D:
			panic("TODO")
		case loong64.AFCOPYSIGN_S:
			panic("TODO")
		case loong64.AFDIV_D:
			panic("TODO")
		case loong64.AFDIV_S:
			panic("TODO")
		case loong64.AFMAXA_D:
			panic("TODO")
		case loong64.AFMAXA_S:
			panic("TODO")
		case loong64.AFMAX_D:
			panic("TODO")
		case loong64.AFMAX_S:
			panic("TODO")
		case loong64.AFMINA_D:
			panic("TODO")
		case loong64.AFMINA_S:
			panic("TODO")
		case loong64.AFMIN_D:
			panic("TODO")
		case loong64.AFMIN_S:
			panic("TODO")
		case loong64.AFMUL_D:
			p.RegF[arg.Rd] = p.RegF[arg.Rs1] * p.RegF[arg.Rs2]
			return nil
		case loong64.AFMUL_S:
			panic("TODO")
		case loong64.AFSCALEB_D:
			panic("TODO")
		case loong64.AFSCALEB_S:
			panic("TODO")
		case loong64.AFSUB_D:
			p.RegF[arg.Rd] = p.RegF[arg.Rs1] - p.RegF[arg.Rs2]
			return nil
		case loong64.AFSUB_S:
			panic("TODO")
		}
	case loong64.OpFormatType_1F_2R:
		switch as {
		case loong64.AFLDGT_D:
			panic("TODO")
		case loong64.AFLDGT_S:
			panic("TODO")
		case loong64.AFLDLE_D:
			panic("TODO")
		case loong64.AFLDLE_S:
			panic("TODO")
		case loong64.AFLDX_D:
			panic("TODO")
		case loong64.AFLDX_S:
			panic("TODO")
		case loong64.AFSTGT_D:
			panic("TODO")
		case loong64.AFSTGT_S:
			panic("TODO")
		case loong64.AFSTLE_D:
			panic("TODO")
		case loong64.AFSTLE_S:
			panic("TODO")
		case loong64.AFSTX_D:
			panic("TODO")
		case loong64.AFSTX_S:
			panic("TODO")
		default:
			return fmt.Errorf("unsupport: %s", loong64.AsString(as, ""))
		}
	case loong64.OpFormatType_4F:
		switch as {
		case loong64.AFMADD_D:
			panic("TODO")
		case loong64.AFMADD_S:
			panic("TODO")
		case loong64.AFMSUB_D:
			panic("TODO")
		case loong64.AFMSUB_S:
			panic("TODO")
		case loong64.AFNMADD_D:
			panic("TODO")
		case loong64.AFNMADD_S:
			panic("TODO")
		case loong64.AFNMSUB_D:
			panic("TODO")
		case loong64.AFNMSUB_S:
			panic("TODO")
		default:
			return fmt.Errorf("unsupport: %s", loong64.AsString(as, ""))
		}
	case loong64.OpFormatType_2R_ui5:
		switch as {
		default:
			return fmt.Errorf("unsupport: %s", loong64.AsString(as, ""))

		case loong64.AROTRI_W:
			panic("TODO")
		case loong64.ASLLI_W:
			p.RegX[arg.Rd] = LAUInt(int64(int32(p.RegX[arg.Rs1]) << uint(arg.Imm)))
			return nil
		case loong64.ASRAI_W:
			p.RegX[arg.Rd] = LAUInt(int64(int32(p.RegX[arg.Rs1]) >> uint(arg.Imm)))
			return nil
		case loong64.ASRLI_W:
			p.RegX[arg.Rd] = LAUInt(int64(uint32(p.RegX[arg.Rs1]) >> uint(arg.Imm)))
			return nil
		}
	case loong64.OpFormatType_2R_ui6:
		switch as {
		case loong64.AROTRI_D:
			panic("TODO")
		case loong64.ASLLI_D:
			panic("TODO")
		case loong64.ASRAI_D:
			panic("TODO")
		case loong64.ASRLI_D:
			panic("TODO")
		default:
			return fmt.Errorf("unsupport: %s", loong64.AsString(as, ""))
		}
	case loong64.OpFormatType_2R_si12:
		switch as {
		default:
			return fmt.Errorf("unsupport: %s", loong64.AsString(as, ""))

		case loong64.AADDI_D:
			panic("TODO")
		case loong64.AADDI_W:
			p.RegX[arg.Rd] = p.RegX[arg.Rs1] + LAUInt(arg.Imm)
			return nil
		case loong64.AFLD_D:
			panic("TODO")
		case loong64.AFLD_S:
			panic("TODO")
		case loong64.AFST_D:
			panic("TODO")
		case loong64.AFST_S:
			panic("TODO")
		case loong64.ALD_B:
			panic("TODO")
		case loong64.ALD_BU:
			addr := p.RegX[arg.Rs1] + LAUInt(arg.Imm)
			value, err := bus.Read(uint64(addr), 1)
			if err != nil {
				return err
			}
			p.RegX[arg.Rd] = LAUInt(uint8(value))
			return nil
		case loong64.ALD_D:
			addr := p.RegX[arg.Rs1] + LAUInt(arg.Imm)
			value, err := bus.Read(uint64(addr), 8)
			if err != nil {
				return err
			}
			p.RegX[arg.Rd] = LAUInt(value)
			return nil
		case loong64.ALD_H:
			panic("TODO")
		case loong64.ALD_HU:
			panic("TODO")
		case loong64.ALD_W:
			panic("TODO")
		case loong64.ALD_WU:
			panic("TODO")
		case loong64.ALU52I_D:
			panic("TODO")
		case loong64.ASLTI:
			panic("TODO")
		case loong64.ASLTUI:
			panic("TODO")
		case loong64.AST_B:
			value := p.RegX[arg.Rd]
			addr := p.RegX[arg.Rs1] + LAUInt(arg.Imm)
			if err := bus.Write(uint64(addr), 1, uint64(value)); err != nil {
				return err
			}
			return nil
		case loong64.AST_D:
			value := p.RegX[arg.Rd]
			addr := p.RegX[arg.Rs1] + LAUInt(arg.Imm)
			if err := bus.Write(uint64(addr), 8, uint64(value)); err != nil {
				return err
			}
			return nil
		case loong64.AST_H:
			panic("TODO")
		case loong64.AST_W:
			value := p.RegX[arg.Rd]
			addr := p.RegX[arg.Rs1] + LAUInt(arg.Imm)
			if err := bus.Write(uint64(addr), 4, uint64(value)); err != nil {
				return err
			}
			return nil
		}
	case loong64.OpFormatType_2R_ui12:
		switch as {
		default:
			return fmt.Errorf("unsupport: %s", loong64.AsString(as, ""))
		case loong64.AANDI:
			panic("TODO")
		case loong64.AORI:
			p.RegX[arg.Rd] = p.RegX[arg.Rs1] | LAUInt(arg.Imm)
			return nil
		case loong64.AXORI:
			panic("TODO")
		}
	case loong64.OpFormatType_2R_si14:
		switch as {
		case loong64.AADDU16I_D:
			panic("TODO")
		case loong64.ALDPTR_D:
			panic("TODO")
		case loong64.ALDPTR_W:
			panic("TODO")
		case loong64.ALL_D:
			panic("TODO")
		case loong64.ALL_W:
			panic("TODO")
		case loong64.ASC_D:
			panic("TODO")
		case loong64.ASC_W:
			panic("TODO")
		case loong64.ASTPTR_D:
			panic("TODO")
		case loong64.ASTPTR_W:
			panic("TODO")
		default:
			return fmt.Errorf("unsupport: %s", loong64.AsString(as, ""))
		}
	case loong64.OpFormatType_1R_si20:
		switch as {
		default:
			return fmt.Errorf("unsupport: %s", loong64.AsString(as, ""))
		case loong64.APCADDU12I:
			p.RegX[arg.Rd] = curPC + LAUInt(arg.Imm)
			return nil
		case loong64.ALU12I_W:
			p.RegX[arg.Rd] = LAUInt(arg.Imm << 12)
			return nil
		}
	case loong64.OpFormatType_0_2R:
		switch as {
		case loong64.AASRTGT_D:
			panic("TODO")
		case loong64.AASRTLE_D:
			panic("TODO")
		default:
			return fmt.Errorf("unsupport: %s", loong64.AsString(as, ""))
		}
	case loong64.OpFormatType_3R_sa2:
		switch as {
		case loong64.AALSL_D:
			panic("TODO")
		case loong64.AALSL_W:
			panic("TODO")
		case loong64.AALSL_WU:
			panic("TODO")
		case loong64.ABYTEPICK_W:
			panic("TODO")
		default:
			return fmt.Errorf("unsupport: %s", loong64.AsString(as, ""))
		}
	case loong64.OpFormatType_3R_sa3:
		switch as {
		case loong64.ABYTEPICK_D:
			panic("TODO")
		default:
			return fmt.Errorf("unsupport: %s", loong64.AsString(as, ""))
		}
	case loong64.OpFormatType_code:
		switch as {
		case loong64.ABREAK:
			panic("TODO")
		case loong64.ADBCL:
			panic("TODO")
		case loong64.ASYSCALL:
			panic("TODO")
		default:
			return fmt.Errorf("unsupport: %s", loong64.AsString(as, ""))
		}
	case loong64.OpFormatType_code_1R_si12:
		switch as {
		case loong64.ACACOP:
			panic("TODO")
		default:
			return fmt.Errorf("unsupport: %s", loong64.AsString(as, ""))
		}
	case loong64.OpFormatType_2R_msbw_lsbw:
		switch as {
		case loong64.ABSTRINS_W:
			panic("TODO")
		case loong64.ABSTRPICK_W:
			panic("TODO")
		default:
			return fmt.Errorf("unsupport: %s", loong64.AsString(as, ""))
		}
	case loong64.OpFormatType_2R_msbd_lsbd:
		switch as {
		case loong64.ABSTRINS_D:
			panic("TODO")
		case loong64.ABSTRPICK_D:
			panic("TODO")
		default:
			return fmt.Errorf("unsupport: %s", loong64.AsString(as, ""))
		}
	case loong64.OpFormatType_fcsr_1R:
		switch as {
		case loong64.AMOVGR2FCSR:
			panic("TODO")
		default:
			return fmt.Errorf("unsupport: %s", loong64.AsString(as, ""))
		}
	case loong64.OpFormatType_1R_fcsr:
		switch as {
		case loong64.AMOVFCSR2GR:
			panic("TODO")
		default:
			return fmt.Errorf("unsupport: %s", loong64.AsString(as, ""))
		}
	case loong64.OpFormatType_cd_1R:
		switch as {
		case loong64.AMOVGR2CF:
			panic("TODO")
		default:
			return fmt.Errorf("unsupport: %s", loong64.AsString(as, ""))
		}
	case loong64.OpFormatType_cd_1F:
		switch as {
		case loong64.AMOVFR2CF:
			panic("TODO")
		default:
			return fmt.Errorf("unsupport: %s", loong64.AsString(as, ""))
		}
	case loong64.OpFormatType_cd_2F:
		switch as {
		case loong64.AFCMP_CAF_D:
			panic("TODO")
		case loong64.AFCMP_CAF_S:
			panic("TODO")
		case loong64.AFCMP_CEQ_D:
			panic("TODO")
		case loong64.AFCMP_CEQ_S:
			panic("TODO")
		case loong64.AFCMP_CLE_D:
			panic("TODO")
		case loong64.AFCMP_CLE_S:
			panic("TODO")
		case loong64.AFCMP_CLT_D:
			panic("TODO")
		case loong64.AFCMP_CLT_S:
			panic("TODO")
		case loong64.AFCMP_CNE_D:
			panic("TODO")
		case loong64.AFCMP_CNE_S:
			panic("TODO")
		case loong64.AFCMP_COR_D:
			panic("TODO")
		case loong64.AFCMP_COR_S:
			panic("TODO")
		case loong64.AFCMP_CUEQ_D:
			panic("TODO")
		case loong64.AFCMP_CUEQ_S:
			panic("TODO")
		case loong64.AFCMP_CULE_D:
			panic("TODO")
		case loong64.AFCMP_CULE_S:
			panic("TODO")
		case loong64.AFCMP_CULT_D:
			panic("TODO")
		case loong64.AFCMP_CULT_S:
			panic("TODO")
		case loong64.AFCMP_CUNE_D:
			panic("TODO")
		case loong64.AFCMP_CUNE_S:
			panic("TODO")
		case loong64.AFCMP_CUN_D:
			panic("TODO")
		case loong64.AFCMP_CUN_S:
			panic("TODO")
		case loong64.AFCMP_SAF_D:
			panic("TODO")
		case loong64.AFCMP_SAF_S:
			panic("TODO")
		case loong64.AFCMP_SEQ_D:
			panic("TODO")
		case loong64.AFCMP_SEQ_S:
			panic("TODO")
		case loong64.AFCMP_SLE_D:
			panic("TODO")
		case loong64.AFCMP_SLE_S:
			panic("TODO")
		case loong64.AFCMP_SLT_D:
			panic("TODO")
		case loong64.AFCMP_SLT_S:
			panic("TODO")
		case loong64.AFCMP_SNE_D:
			panic("TODO")
		case loong64.AFCMP_SNE_S:
			panic("TODO")
		case loong64.AFCMP_SOR_D:
			panic("TODO")
		case loong64.AFCMP_SOR_S:
			panic("TODO")
		case loong64.AFCMP_SUEQ_D:
			panic("TODO")
		case loong64.AFCMP_SUEQ_S:
			panic("TODO")
		case loong64.AFCMP_SULE_D:
			panic("TODO")
		case loong64.AFCMP_SULE_S:
			panic("TODO")
		case loong64.AFCMP_SULT_D:
			panic("TODO")
		case loong64.AFCMP_SULT_S:
			panic("TODO")
		case loong64.AFCMP_SUNE_D:
			panic("TODO")
		case loong64.AFCMP_SUNE_S:
			panic("TODO")
		case loong64.AFCMP_SUN_D:
			panic("TODO")
		case loong64.AFCMP_SUN_S:
			panic("TODO")
		default:
			return fmt.Errorf("unsupport: %s", loong64.AsString(as, ""))
		}
	case loong64.OpFormatType_1R_cj:
		switch as {
		case loong64.AMOVCF2GR:
			panic("TODO")
		default:
			return fmt.Errorf("unsupport: %s", loong64.AsString(as, ""))
		}
	case loong64.OpFormatType_1F_cj:
		switch as {
		case loong64.AMOVCF2FR:
			panic("TODO")
		default:
			return fmt.Errorf("unsupport: %s", loong64.AsString(as, ""))
		}
	case loong64.OpFormatType_1R_csr:
		switch as {
		case loong64.ACSRRD:
			panic("TODO")
		case loong64.ACSRWR:
			panic("TODO")
		default:
			return fmt.Errorf("unsupport: %s", loong64.AsString(as, ""))
		}
	case loong64.OpFormatType_2R_csr:
		switch as {
		case loong64.ACSRXCHG:
			panic("TODO")
		default:
			return fmt.Errorf("unsupport: %s", loong64.AsString(as, ""))
		}
	case loong64.OpFormatType_2R_level:
		switch as {
		case loong64.ALDDIR:
			panic("TODO")
		default:
			return fmt.Errorf("unsupport: %s", loong64.AsString(as, ""))
		}
	case loong64.OpFormatType_level:
		switch as {
		case loong64.AIDLE:
			panic("TODO")
		default:
			return fmt.Errorf("unsupport: %s", loong64.AsString(as, ""))
		}
	case loong64.OpFormatType_0_1R_seq:
		switch as {
		case loong64.ALDPTE:
			panic("TODO")
		default:
			return fmt.Errorf("unsupport: %s", loong64.AsString(as, ""))
		}
	case loong64.OpFormatType_op_2R:
		switch as {
		case loong64.AINVTLB:
			panic("TODO")
		default:
			return fmt.Errorf("unsupport: %s", loong64.AsString(as, ""))
		}
	case loong64.OpFormatType_3F_ca:
		switch as {
		case loong64.AFSEL:
			panic("TODO")
		default:
			return fmt.Errorf("unsupport: %s", loong64.AsString(as, ""))
		}
	case loong64.OpFormatType_hint_1R_si12:
		switch as {
		case loong64.APRELD:
			panic("TODO")
		default:
			return fmt.Errorf("unsupport: %s", loong64.AsString(as, ""))
		}
	case loong64.OpFormatType_hint_2R:
		switch as {
		case loong64.APRELDX:
			panic("TODO")
		default:
			return fmt.Errorf("unsupport: %s", loong64.AsString(as, ""))
		}
	case loong64.OpFormatType_hint:
		switch as {
		case loong64.ADBAR:
			panic("TODO")
		case loong64.AIBAR:
			panic("TODO")
		default:
			return fmt.Errorf("unsupport: %s", loong64.AsString(as, ""))
		}
	case loong64.OpFormatType_cj_offset:
		switch as {
		case loong64.ABCEQZ:
			panic("TODO")
		case loong64.ABCNEZ:
			panic("TODO")
		default:
			return fmt.Errorf("unsupport: %s", loong64.AsString(as, ""))
		}
	case loong64.OpFormatType_rj_offset:
		switch as {
		case loong64.ABEQZ:
			panic("TODO")
		case loong64.ABNEZ:
			panic("TODO")
		default:
			return fmt.Errorf("unsupport: %s", loong64.AsString(as, ""))
		}
	case loong64.OpFormatType_rj_rd_offset:
		switch as {
		default:
			return fmt.Errorf("unsupport: %s", loong64.AsString(as, ""))

		case loong64.ABEQ:
			if p.RegX[arg.Rs1] == p.RegX[arg.Rs2] {
				p.PC = curPC + LAUInt(arg.Imm)
			}
			return nil
		case loong64.ABGE:
			panic("TODO")
		case loong64.ABGEU:
			panic("TODO")
		case loong64.ABLT:
			if int64(p.RegX[arg.Rs1]) < int64(p.RegX[arg.Rs2]) {
				p.PC = curPC + LAUInt(arg.Imm)
			}
			return nil
		case loong64.ABLTU:
			panic("TODO")
		case loong64.ABNE:
			if p.RegX[arg.Rs1] != p.RegX[arg.Rs2] {
				p.PC = curPC + LAUInt(arg.Imm)
			}
			return nil
		}
	case loong64.OpFormatType_rd_rj_offset:
		switch as {
		case loong64.AJIRL:
			panic("TODO")
		default:
			return fmt.Errorf("unsupport: %s", loong64.AsString(as, ""))
		}
	case loong64.OpFormatType_offset:
		switch as {
		default:
			return fmt.Errorf("unsupport: %s", loong64.AsString(as, ""))
		case loong64.AB:
			p.PC = curPC + LAUInt(arg.Imm)
			return nil
		case loong64.ABL:
			p.RegX[1] = curPC + 4
			p.PC = curPC + LAUInt(arg.Imm)
			return nil
		}
	default:
		panic("unreachable")
	}
}

func (p *CPU) InstString(x uint32) (string, error) {
	as, arg, err := loong64.Decode(x)
	if err != nil {
		return "", err
	}
	return loong64.AsmSyntax(as, "", arg), nil
}

func (p *CPU) LookupRegister(regName string) (r abi.RegType, ok bool) {
	return loong64.LookupRegister(regName)
}

func (p *CPU) RegAliasString(r abi.RegType) string {
	return loong64.RegAliasString(r)
}
