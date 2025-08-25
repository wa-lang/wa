// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package riscv

import "testing"

func TestAOpContextTable_opcode(t *testing.T) {
	for i := AADDI; i < ALAST; i++ {
		ctx := AOpContextTable[i]
		if ctx.PseudoAs != 0 {
			continue
		}
		if ctx.Opcode == 0 {
			t.Fatalf("invalid AOpContextTable[%s]", AsString(i))
		}
	}
}

func TestAOpContextTable(t *testing.T) {
	for i := AADDI; i < ALAST; i++ {
		ctx := AOpContextTable[i]
		if ctx.PseudoAs == 0 {
			if ctx.Opcode == 0 {
				t.Fatalf("%s: invalid opcode", AsString(i))
			}
		}
		if ctx.ArgMarks&ARG_FUNCT3 != 0 {
			if ctx.Funct3 > 0b_111 {
				t.Fatalf("%s: invalid funct3", AsString(i))
			}
		}
		if ctx.ArgMarks&ARG_FUNCT7 != 0 {
			if ctx.Funct7 > 0b_111_1111 {
				t.Fatalf("%s: invalid funct7", AsString(i))
			}
		}
		if ctx.ArgMarks&ARG_FUNCT2 != 0 {
			if ctx.Funct7 > 0b_11 {
				t.Fatalf("%s: invalid funct2", AsString(i))
			}
		}
		if ctx.ArgMarks&ARG_IMM != 0 {
			if ctx.ImmMin >= ctx.ImmMax {
				// TODO: 填充 imm 范围信息
				// t.Fatalf("%s: invalid imm range: [%d, %d]", AsString(i), ctx.ImmMin, ctx.ImmMax)
			}
		}
	}
}
