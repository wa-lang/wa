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

func (p *wat2laWorker) comment(w io.Writer, format string, a ...interface{}) {
	fmt.Fprintln(w, "   ", fmt.Sprintf(format, a...))
}

func (p *wat2laWorker) st_w(w io.Writer, rd, rj string, offset int32, tmpReg string) {
	assert(rd != "")
	assert(rj != "")
	assert(tmpReg != rd && tmpReg != rj)

	if !p.isS12Overflow(offset) {
		fmt.Fprintf(w, "    st.w %s, %s, %d\n", rd, rj, offset)
	} else {
		hi20 := uint32(offset) >> 12
		lo12 := uint32(offset) & 0xFFF
		fmt.Fprintf(w, "    # st.w %s, %s, %d\n", rd, rj, offset)
		fmt.Fprintf(w, "    lu12i.w %s, %d\n", tmpReg, i32SignExtend(hi20, 20))
		fmt.Fprintf(w, "    ori     %s, %s, 0x%X\n", tmpReg, tmpReg, lo12)
		fmt.Fprintf(w, "    add.d   %s, %s, %s\n", tmpReg, tmpReg, rj)
		fmt.Fprintf(w, "    st.w    %s, %s, 0\n", rd, tmpReg)
	}
}

func (p *wat2laWorker) st_d(w io.Writer, rd, rj string, offset int32, tmpReg string) {
	assert(rd != "")
	assert(rj != "")
	assert(tmpReg != rd && tmpReg != rj)

	if !p.isS12Overflow(offset) {
		fmt.Fprintf(w, "    st.d %s, %s, %d\n", rd, rj, offset)
	} else {
		hi20 := uint32(offset) >> 12
		lo12 := uint32(offset) & 0xFFF
		fmt.Fprintf(w, "    # st.d %s, %s, %d\n", rd, rj, offset)
		fmt.Fprintf(w, "    lu12i.w %s, %d\n", tmpReg, i32SignExtend(hi20, 20))
		fmt.Fprintf(w, "    ori     %s, %s, 0x%X\n", tmpReg, tmpReg, lo12)
		fmt.Fprintf(w, "    add.d   %s, %s, %s\n", tmpReg, tmpReg, rj)
		fmt.Fprintf(w, "    st.d    %s, %s, 0\n", rd, tmpReg)
	}
}

func (p *wat2laWorker) fst_s(w io.Writer, rd, rj string, offset int32, tmpReg string) {
	assert(rd != "")
	assert(rj != "")
	assert(tmpReg != rd && tmpReg != rj)

	if !p.isS12Overflow(offset) {
		fmt.Fprintf(w, "    fst.s %s, %s, %d\n", rd, rj, offset)
	} else {
		hi20 := uint32(offset) >> 12
		lo12 := uint32(offset) & 0xFFF
		fmt.Fprintf(w, "    # fst.s %s, %s, %d\n", rd, rj, offset)
		fmt.Fprintf(w, "    lu12i.w %s, %d\n", tmpReg, i32SignExtend(hi20, 20))
		fmt.Fprintf(w, "    ori     %s, %s, 0x%X\n", tmpReg, tmpReg, lo12)
		fmt.Fprintf(w, "    add.d   %s, %s, %s\n", tmpReg, tmpReg, rj)
		fmt.Fprintf(w, "    fst.s   %s, %s, 0\n", rd, tmpReg)
	}
}

func (p *wat2laWorker) fst_d(w io.Writer, rd, rj string, offset int32, tmpReg string) {
	assert(rd != "")
	assert(rj != "")
	assert(tmpReg != rd && tmpReg != rj)

	if !p.isS12Overflow(offset) {
		fmt.Fprintf(w, "    fst.d %s, %s, %d\n", rd, rj, offset)
	} else {
		hi20 := uint32(offset) >> 12
		lo12 := uint32(offset) & 0xFFF
		fmt.Fprintf(w, "    # fst.d %s, %s, %d\n", rd, rj, offset)
		fmt.Fprintf(w, "    lu12i.w %s, %d\n", tmpReg, i32SignExtend(hi20, 20))
		fmt.Fprintf(w, "    ori     %s, %s, 0x%X\n", tmpReg, tmpReg, lo12)
		fmt.Fprintf(w, "    add.d   %s, %s, %s\n", tmpReg, tmpReg, rj)
		fmt.Fprintf(w, "    fst.d   %s, %s, 0\n", rd, tmpReg)
	}
}

func (p *wat2laWorker) ld_b(w io.Writer, rd, rj string, offset int32, tmpReg string) {
	assert(rd != "")
	assert(rj != "")
	assert(tmpReg != rd && tmpReg != rj)

	if !p.isS12Overflow(offset) {
		fmt.Fprintf(w, "    ld.b %s, %s, %d\n", rd, rj, offset)
	} else {
		hi20 := uint32(offset) >> 12
		lo12 := uint32(offset) & 0xFFF
		fmt.Fprintf(w, "    # ld.b %s, %s, %d\n", rd, rj, offset)
		fmt.Fprintf(w, "    lu12i.w %s, %d\n", tmpReg, i32SignExtend(hi20, 20))
		fmt.Fprintf(w, "    ori     %s, %s, 0x%X\n", tmpReg, tmpReg, lo12)
		fmt.Fprintf(w, "    add.d   %s, %s, $fp\n", tmpReg, tmpReg)
		fmt.Fprintf(w, "    ld.b    %s, %s, 0\n", rd, tmpReg)
	}
}

func (p *wat2laWorker) ld_bu(w io.Writer, rd, rj string, offset int32, tmpReg string) {
	assert(rd != "")
	assert(rj != "")
	assert(tmpReg != rd && tmpReg != rj)

	if !p.isS12Overflow(offset) {
		fmt.Fprintf(w, "    ld.bu %s, %s, %d\n", rd, rj, offset)
	} else {
		hi20 := uint32(offset) >> 12
		lo12 := uint32(offset) & 0xFFF
		fmt.Fprintf(w, "    # ld.bu %s, %s, %d\n", rd, rj, offset)
		fmt.Fprintf(w, "    lu12i.w %s, %d\n", tmpReg, i32SignExtend(hi20, 20))
		fmt.Fprintf(w, "    ori     %s, %s, 0x%X\n", tmpReg, tmpReg, lo12)
		fmt.Fprintf(w, "    add.d   %s, %s, $fp\n", tmpReg, tmpReg)
		fmt.Fprintf(w, "    ld.bu   %s, %s, 0\n", rd, tmpReg)
	}
}

func (p *wat2laWorker) ld_h(w io.Writer, rd, rj string, offset int32, tmpReg string) {
	assert(rd != "")
	assert(rj != "")
	assert(tmpReg != rd && tmpReg != rj)

	if !p.isS12Overflow(offset) {
		fmt.Fprintf(w, "    ld.h %s, %s, %d\n", rd, rj, offset)
	} else {
		hi20 := uint32(offset) >> 12
		lo12 := uint32(offset) & 0xFFF
		fmt.Fprintf(w, "    # ld.h %s, %s, %d\n", rd, rj, offset)
		fmt.Fprintf(w, "    lu12i.w %s, %d\n", tmpReg, i32SignExtend(hi20, 20))
		fmt.Fprintf(w, "    ori     %s, %s, 0x%X\n", tmpReg, tmpReg, lo12)
		fmt.Fprintf(w, "    add.d   %s, %s, $fp\n", tmpReg, tmpReg)
		fmt.Fprintf(w, "    ld.h    %s, %s, 0\n", rd, tmpReg)
	}
}

func (p *wat2laWorker) ld_hu(w io.Writer, rd, rj string, offset int32, tmpReg string) {
	assert(rd != "")
	assert(rj != "")
	assert(tmpReg != rd && tmpReg != rj)

	if !p.isS12Overflow(offset) {
		fmt.Fprintf(w, "    ld.hu %s, %s, %d\n", rd, rj, offset)
	} else {
		hi20 := uint32(offset) >> 12
		lo12 := uint32(offset) & 0xFFF
		fmt.Fprintf(w, "    # ld.hu %s, %s, %d\n", rd, rj, offset)
		fmt.Fprintf(w, "    lu12i.w %s, %d\n", tmpReg, i32SignExtend(hi20, 20))
		fmt.Fprintf(w, "    ori     %s, %s, 0x%X\n", tmpReg, tmpReg, lo12)
		fmt.Fprintf(w, "    add.d   %s, %s, $fp\n", tmpReg, tmpReg)
		fmt.Fprintf(w, "    ld.hu   %s, %s, 0\n", rd, tmpReg)
	}
}

func (p *wat2laWorker) ld_w(w io.Writer, rd, rj string, offset int32, tmpReg string) {
	assert(rd != "")
	assert(rj != "")
	assert(tmpReg != rd && tmpReg != rj)

	if !p.isS12Overflow(offset) {
		fmt.Fprintf(w, "    ld.w %s, %s, %d\n", rd, rj, offset)
	} else {
		hi20 := uint32(offset) >> 12
		lo12 := uint32(offset) & 0xFFF
		fmt.Fprintf(w, "    # ld.w %s, %s, %d\n", rd, rj, offset)
		fmt.Fprintf(w, "    lu12i.w %s, %d\n", tmpReg, i32SignExtend(hi20, 20))
		fmt.Fprintf(w, "    ori     %s, %s, 0x%X\n", tmpReg, tmpReg, lo12)
		fmt.Fprintf(w, "    add.d   %s, %s, $fp\n", tmpReg, tmpReg)
		fmt.Fprintf(w, "    ld.w    %s, %s, 0\n", rd, tmpReg)
	}
}

func (p *wat2laWorker) ld_wu(w io.Writer, rd, rj string, offset int32, tmpReg string) {
	assert(rd != "")
	assert(rj != "")
	assert(tmpReg != rd && tmpReg != rj)

	if !p.isS12Overflow(offset) {
		fmt.Fprintf(w, "    ld.wu %s, %s, %d\n", rd, rj, offset)
	} else {
		hi20 := uint32(offset) >> 12
		lo12 := uint32(offset) & 0xFFF
		fmt.Fprintf(w, "    # ld.wu %s, %s, %d\n", rd, rj, offset)
		fmt.Fprintf(w, "    lu12i.w %s, %d\n", tmpReg, i32SignExtend(hi20, 20))
		fmt.Fprintf(w, "    ori     %s, %s, 0x%X\n", tmpReg, tmpReg, lo12)
		fmt.Fprintf(w, "    add.d   %s, %s, $fp\n", tmpReg, tmpReg)
		fmt.Fprintf(w, "    ld.wu   %s, %s, 0\n", rd, tmpReg)
	}
}

func (p *wat2laWorker) ld_d(w io.Writer, rd, rj string, offset int32, tmpReg string) {
	assert(rd != "")
	assert(rj != "")
	assert(tmpReg != rd && tmpReg != rj)

	if !p.isS12Overflow(offset) {
		fmt.Fprintf(w, "    ld.d %s, %s, %d\n", rd, rj, offset)
	} else {
		hi20 := uint32(offset) >> 12
		lo12 := uint32(offset) & 0xFFF
		fmt.Fprintf(w, "    # ld.d %s, %s, %d\n", rd, rj, offset)
		fmt.Fprintf(w, "    lu12i.w %s, %d\n", tmpReg, i32SignExtend(hi20, 20))
		fmt.Fprintf(w, "    ori     %s, %s, 0x%X\n", tmpReg, tmpReg, lo12)
		fmt.Fprintf(w, "    add.d   %s, %s, $fp\n", tmpReg, tmpReg)
		fmt.Fprintf(w, "    ld.d    %s, %s, 0\n", rd, tmpReg)
	}
}

func (p *wat2laWorker) fld_s(w io.Writer, rd, rj string, offset int32, tmpReg string) {
	assert(rd != "")
	assert(rj != "")
	assert(tmpReg != rd && tmpReg != rj)

	if !p.isS12Overflow(offset) {
		fmt.Fprintf(w, "    fld.s %s, %s, %d\n", rd, rj, offset)
	} else {
		hi20 := uint32(offset) >> 12
		lo12 := uint32(offset) & 0xFFF
		fmt.Fprintf(w, "    # fld.s %s, %s, %d\n", rd, rj, offset)
		fmt.Fprintf(w, "    lu12i.w %s, %d\n", tmpReg, i32SignExtend(hi20, 20))
		fmt.Fprintf(w, "    ori     %s, %s, 0x%X\n", tmpReg, tmpReg, lo12)
		fmt.Fprintf(w, "    add.d   %s, %s, $fp\n", tmpReg, tmpReg)
		fmt.Fprintf(w, "    fld.s   %s, %s, 0\n", rd, tmpReg)
	}
}

func (p *wat2laWorker) fld_d(w io.Writer, rd, rj string, offset int32, tmpReg string) {
	assert(rd != "")
	assert(rj != "")
	assert(tmpReg != rd && tmpReg != rj)

	if !p.isS12Overflow(offset) {
		fmt.Fprintf(w, "    fld.d %s, %s, %d\n", rd, rj, offset)
	} else {
		hi20 := uint32(offset) >> 12
		lo12 := uint32(offset) & 0xFFF
		fmt.Fprintf(w, "    # fld.d %s, %s, %d\n", rd, rj, offset)
		fmt.Fprintf(w, "    lu12i.w %s, %d\n", tmpReg, i32SignExtend(hi20, 20))
		fmt.Fprintf(w, "    ori     %s, %s, 0x%X\n", tmpReg, tmpReg, lo12)
		fmt.Fprintf(w, "    add.d   %s, %s, $fp\n", tmpReg, tmpReg)
		fmt.Fprintf(w, "    fld.d   %s, %s, 0\n", rd, tmpReg)
	}
}
