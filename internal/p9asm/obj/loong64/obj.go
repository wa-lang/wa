// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package loong64

import "wa-lang.org/wa/internal/p9asm/obj"

var Linkloong64 = obj.LinkArch{
	//Arch:           sys.ArchLoong64,
	//Init:           buildop,
	//Preprocess:     preprocess,
	//Assemble:       span0,
	//Progedit:       progedit,
	UnaryDst: map[obj.As]bool{},
	//DWARFRegisters: LOONG64DWARFRegisters,
}
