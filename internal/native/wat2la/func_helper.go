// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2la

import (
	"fmt"
	"io"
)

func (p *wat2laWorker) isS12Overflow(x int32) bool {
	return x < -2048 || x > 2047
}

func (p *wat2laWorker) emitFPStore(w io.Writer, stCmd string, srcReg, tmpReg string, offset int32) {
	assert(stCmd == "st.w" || stCmd == "st.d" || stCmd == "fst.s" || stCmd == "fst.d")
	assert(srcReg != "")
	assert(tmpReg != "")
	assert(srcReg != tmpReg)
	if !p.isS12Overflow(offset) {
		fmt.Fprintf(w, "    %s %s, $fp, %d\n", stCmd, srcReg, offset)
	} else {
		hi20 := uint32(offset) >> 12
		lo12 := uint32(offset) & 0xFFF
		fmt.Fprintf(w, "    lu12i.w %s, %d\n", tmpReg, i32SignExtend(hi20, 20))
		fmt.Fprintf(w, "    ori     %s, %s, 0x%X\n", tmpReg, tmpReg, lo12)
		fmt.Fprintf(w, "    add.d   %s, %s, $fp\n", tmpReg, tmpReg)
		if len(stCmd) == len("st.d") {
			fmt.Fprintf(w, "    %s    %s, %s, 0\n", stCmd, srcReg, tmpReg)
		} else {
			fmt.Fprintf(w, "    %s   %s, %s, 0\n", stCmd, srcReg, tmpReg)
		}
	}
}

func (p *wat2laWorker) emitFPLoad(w io.Writer, ldCmd string, srcReg, tmpReg string, offset int32) {
	assert(
		ldCmd == "ld.w" || ldCmd == "ld.d" ||
			ldCmd == "ld.bu" || ldCmd == "ld.hu" || ldCmd == "ld.wu" ||
			ldCmd == "fld.s" || ldCmd == "fld.d",
	)
	assert(srcReg != "")
	assert(tmpReg != "")
	assert(srcReg != tmpReg)
	if !p.isS12Overflow(offset) {
		fmt.Fprintf(w, "    %s %s, $fp, %d\n", ldCmd, srcReg, offset)
	} else {
		hi20 := uint32(offset) >> 12
		lo12 := uint32(offset) & 0xFFF
		fmt.Fprintf(w, "    lu12i.w %s, %d\n", tmpReg, i32SignExtend(hi20, 20))
		fmt.Fprintf(w, "    ori     %s, %s, 0x%X\n", tmpReg, tmpReg, lo12)
		fmt.Fprintf(w, "    add.d   %s, %s, $fp\n", tmpReg, tmpReg)
		if len(ldCmd) == len("ld.w") {
			fmt.Fprintf(w, "    %s    %s, %s, 0\n", ldCmd, srcReg, tmpReg)
		} else {
			fmt.Fprintf(w, "    %s   %s, %s, 0\n", ldCmd, srcReg, tmpReg)
		}
	}
}
