package riscv

import "wa-lang.org/wa/internal/p9asm/obj"

var LinkRISCV64 = obj.LinkArch{
	//Arch:           sys.ArchRISCV64,
	//Init:           buildop,
	//Preprocess:     preprocess,
	//Assemble:       assemble,
	//Progedit:       progedit,
	UnaryDst: unaryDst,
	//DWARFRegisters: RISCV64DWARFRegisters,
}
