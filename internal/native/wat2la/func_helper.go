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

func (p *wat2laWorker) emitFP_st_w(w io.Writer, srcReg, tmpReg string, offset int32) {
	assert(srcReg != "")
	assert(tmpReg != "")
	assert(srcReg != tmpReg)

	if !p.isS12Overflow(offset) {
		fmt.Fprintf(w, "    st.w %s, $fp, %d\n", srcReg, offset)
	} else {
		hi20 := uint32(offset) >> 12
		lo12 := uint32(offset) & 0xFFF
		fmt.Fprintf(w, "    lu12i.w %s, %d\n", tmpReg, i32SignExtend(hi20, 20))
		fmt.Fprintf(w, "    ori     %s, %s, 0x%X\n", tmpReg, tmpReg, lo12)
		fmt.Fprintf(w, "    add.d   %s, %s, $fp\n", tmpReg, tmpReg)
		fmt.Fprintf(w, "    st.w    %s, %s, 0\n", srcReg, tmpReg)
	}
}

func (p *wat2laWorker) emitFP_st_d(w io.Writer, srcReg, tmpReg string, offset int32) {
	assert(srcReg != "")
	assert(tmpReg != "")
	assert(srcReg != tmpReg)

	if !p.isS12Overflow(offset) {
		fmt.Fprintf(w, "    st.d %s, $fp, %d\n", srcReg, offset)
	} else {
		hi20 := uint32(offset) >> 12
		lo12 := uint32(offset) & 0xFFF
		fmt.Fprintf(w, "    lu12i.w %s, %d\n", tmpReg, i32SignExtend(hi20, 20))
		fmt.Fprintf(w, "    ori     %s, %s, 0x%X\n", tmpReg, tmpReg, lo12)
		fmt.Fprintf(w, "    add.d   %s, %s, $fp\n", tmpReg, tmpReg)
		fmt.Fprintf(w, "    st.d    %s, %s, 0\n", srcReg, tmpReg)
	}
}
func (p *wat2laWorker) emitFP_fst_s(w io.Writer, srcReg, tmpReg string, offset int32) {
	assert(srcReg != "")
	assert(tmpReg != "")
	assert(srcReg != tmpReg)

	if !p.isS12Overflow(offset) {
		fmt.Fprintf(w, "    fst.s %s, $fp, %d\n", srcReg, offset)
	} else {
		hi20 := uint32(offset) >> 12
		lo12 := uint32(offset) & 0xFFF
		fmt.Fprintf(w, "    lu12i.w %s, %d\n", tmpReg, i32SignExtend(hi20, 20))
		fmt.Fprintf(w, "    ori     %s, %s, 0x%X\n", tmpReg, tmpReg, lo12)
		fmt.Fprintf(w, "    add.d   %s, %s, $fp\n", tmpReg, tmpReg)
		fmt.Fprintf(w, "    fst.s   %s, %s, 0\n", srcReg, tmpReg)
	}
}
func (p *wat2laWorker) emitFP_fst_d(w io.Writer, srcReg, tmpReg string, offset int32) {
	assert(srcReg != "")
	assert(tmpReg != "")
	assert(srcReg != tmpReg)

	if !p.isS12Overflow(offset) {
		fmt.Fprintf(w, "    fst.d %s, $fp, %d\n", srcReg, offset)
	} else {
		hi20 := uint32(offset) >> 12
		lo12 := uint32(offset) & 0xFFF
		fmt.Fprintf(w, "    lu12i.w %s, %d\n", tmpReg, i32SignExtend(hi20, 20))
		fmt.Fprintf(w, "    ori     %s, %s, 0x%X\n", tmpReg, tmpReg, lo12)
		fmt.Fprintf(w, "    add.d   %s, %s, $fp\n", tmpReg, tmpReg)
		fmt.Fprintf(w, "    fst.d   %s, %s, 0\n", srcReg, tmpReg)
	}
}

func (p *wat2laWorker) emitFP_ld_b(w io.Writer, srcReg, tmpReg string, offset int32) {
	assert(srcReg != "")
	assert(tmpReg != "")
	assert(srcReg != tmpReg)

	if !p.isS12Overflow(offset) {
		fmt.Fprintf(w, "    ld.b %s, $fp, %d\n", srcReg, offset)
	} else {
		hi20 := uint32(offset) >> 12
		lo12 := uint32(offset) & 0xFFF
		fmt.Fprintf(w, "    lu12i.w %s, %d\n", tmpReg, i32SignExtend(hi20, 20))
		fmt.Fprintf(w, "    ori     %s, %s, 0x%X\n", tmpReg, tmpReg, lo12)
		fmt.Fprintf(w, "    add.d   %s, %s, $fp\n", tmpReg, tmpReg)
		fmt.Fprintf(w, "    ld.b    %s, %s, 0\n", srcReg, tmpReg)
	}
}
func (p *wat2laWorker) emitFP_ld_h(w io.Writer, srcReg, tmpReg string, offset int32) {
	assert(srcReg != "")
	assert(tmpReg != "")
	assert(srcReg != tmpReg)

	if !p.isS12Overflow(offset) {
		fmt.Fprintf(w, "    ld.h %s, $fp, %d\n", srcReg, offset)
	} else {
		hi20 := uint32(offset) >> 12
		lo12 := uint32(offset) & 0xFFF
		fmt.Fprintf(w, "    lu12i.w %s, %d\n", tmpReg, i32SignExtend(hi20, 20))
		fmt.Fprintf(w, "    ori     %s, %s, 0x%X\n", tmpReg, tmpReg, lo12)
		fmt.Fprintf(w, "    add.d   %s, %s, $fp\n", tmpReg, tmpReg)
		fmt.Fprintf(w, "    ld.h    %s, %s, 0\n", srcReg, tmpReg)
	}
}
func (p *wat2laWorker) emitFP_ld_w(w io.Writer, srcReg, tmpReg string, offset int32) {
	assert(srcReg != "")
	assert(tmpReg != "")
	assert(srcReg != tmpReg)

	if !p.isS12Overflow(offset) {
		fmt.Fprintf(w, "    ld.w %s, $fp, %d\n", srcReg, offset)
	} else {
		hi20 := uint32(offset) >> 12
		lo12 := uint32(offset) & 0xFFF
		fmt.Fprintf(w, "    lu12i.w %s, %d\n", tmpReg, i32SignExtend(hi20, 20))
		fmt.Fprintf(w, "    ori     %s, %s, 0x%X\n", tmpReg, tmpReg, lo12)
		fmt.Fprintf(w, "    add.d   %s, %s, $fp\n", tmpReg, tmpReg)
		fmt.Fprintf(w, "    ld.w    %s, %s, 0\n", srcReg, tmpReg)
	}
}
func (p *wat2laWorker) emitFP_ld_d(w io.Writer, srcReg, tmpReg string, offset int32) {
	assert(srcReg != "")
	assert(tmpReg != "")
	assert(srcReg != tmpReg)

	if !p.isS12Overflow(offset) {
		fmt.Fprintf(w, "    ld.d %s, $fp, %d\n", srcReg, offset)
	} else {
		hi20 := uint32(offset) >> 12
		lo12 := uint32(offset) & 0xFFF
		fmt.Fprintf(w, "    lu12i.w %s, %d\n", tmpReg, i32SignExtend(hi20, 20))
		fmt.Fprintf(w, "    ori     %s, %s, 0x%X\n", tmpReg, tmpReg, lo12)
		fmt.Fprintf(w, "    add.d   %s, %s, $fp\n", tmpReg, tmpReg)
		fmt.Fprintf(w, "    ld.d    %s, %s, 0\n", srcReg, tmpReg)
	}
}

func (p *wat2laWorker) emitFP_ld_bu(w io.Writer, srcReg, tmpReg string, offset int32) {
	assert(srcReg != "")
	assert(tmpReg != "")
	assert(srcReg != tmpReg)

	if !p.isS12Overflow(offset) {
		fmt.Fprintf(w, "    ld.bu %s, $fp, %d\n", srcReg, offset)
	} else {
		hi20 := uint32(offset) >> 12
		lo12 := uint32(offset) & 0xFFF
		fmt.Fprintf(w, "    lu12i.w %s, %d\n", tmpReg, i32SignExtend(hi20, 20))
		fmt.Fprintf(w, "    ori     %s, %s, 0x%X\n", tmpReg, tmpReg, lo12)
		fmt.Fprintf(w, "    add.d   %s, %s, $fp\n", tmpReg, tmpReg)
		fmt.Fprintf(w, "    ld.bu   %s, %s, 0\n", srcReg, tmpReg)
	}
}
func (p *wat2laWorker) emitFP_ld_hu(w io.Writer, srcReg, tmpReg string, offset int32) {
	assert(srcReg != "")
	assert(tmpReg != "")
	assert(srcReg != tmpReg)

	if !p.isS12Overflow(offset) {
		fmt.Fprintf(w, "    ld.hu %s, $fp, %d\n", srcReg, offset)
	} else {
		hi20 := uint32(offset) >> 12
		lo12 := uint32(offset) & 0xFFF
		fmt.Fprintf(w, "    lu12i.w %s, %d\n", tmpReg, i32SignExtend(hi20, 20))
		fmt.Fprintf(w, "    ori     %s, %s, 0x%X\n", tmpReg, tmpReg, lo12)
		fmt.Fprintf(w, "    add.d   %s, %s, $fp\n", tmpReg, tmpReg)
		fmt.Fprintf(w, "    ld.hu   %s, %s, 0\n", srcReg, tmpReg)
	}
}
func (p *wat2laWorker) emitFP_ld_wu(w io.Writer, srcReg, tmpReg string, offset int32) {
	assert(srcReg != "")
	assert(tmpReg != "")
	assert(srcReg != tmpReg)

	if !p.isS12Overflow(offset) {
		fmt.Fprintf(w, "    ld.wu %s, $fp, %d\n", srcReg, offset)
	} else {
		hi20 := uint32(offset) >> 12
		lo12 := uint32(offset) & 0xFFF
		fmt.Fprintf(w, "    lu12i.w %s, %d\n", tmpReg, i32SignExtend(hi20, 20))
		fmt.Fprintf(w, "    ori     %s, %s, 0x%X\n", tmpReg, tmpReg, lo12)
		fmt.Fprintf(w, "    add.d   %s, %s, $fp\n", tmpReg, tmpReg)
		fmt.Fprintf(w, "    ld.wu   %s, %s, 0\n", srcReg, tmpReg)
	}
}

func (p *wat2laWorker) emitFP_fld_s(w io.Writer, srcReg, tmpReg string, offset int32) {
	assert(srcReg != "")
	assert(tmpReg != "")
	assert(srcReg != tmpReg)

	if !p.isS12Overflow(offset) {
		fmt.Fprintf(w, "    fld.s %s, $fp, %d\n", srcReg, offset)
	} else {
		hi20 := uint32(offset) >> 12
		lo12 := uint32(offset) & 0xFFF
		fmt.Fprintf(w, "    lu12i.w %s, %d\n", tmpReg, i32SignExtend(hi20, 20))
		fmt.Fprintf(w, "    ori     %s, %s, 0x%X\n", tmpReg, tmpReg, lo12)
		fmt.Fprintf(w, "    add.d   %s, %s, $fp\n", tmpReg, tmpReg)
		fmt.Fprintf(w, "    fld.s   %s, %s, 0\n", srcReg, tmpReg)
	}
}
func (p *wat2laWorker) emitFP_fld_d(w io.Writer, srcReg, tmpReg string, offset int32) {
	assert(srcReg != "")
	assert(tmpReg != "")
	assert(srcReg != tmpReg)

	if !p.isS12Overflow(offset) {
		fmt.Fprintf(w, "    fld.d %s, $fp, %d\n", srcReg, offset)
	} else {
		hi20 := uint32(offset) >> 12
		lo12 := uint32(offset) & 0xFFF
		fmt.Fprintf(w, "    lu12i.w %s, %d\n", tmpReg, i32SignExtend(hi20, 20))
		fmt.Fprintf(w, "    ori     %s, %s, 0x%X\n", tmpReg, tmpReg, lo12)
		fmt.Fprintf(w, "    add.d   %s, %s, $fp\n", tmpReg, tmpReg)
		fmt.Fprintf(w, "    fld.d   %s, %s, 0\n", srcReg, tmpReg)
	}
}
