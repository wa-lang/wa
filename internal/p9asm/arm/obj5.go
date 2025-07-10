package arm

import "wa-lang.org/wa/internal/p9asm/obj"

var unaryDst = map[int]bool{
	ASWI:  true,
	AWORD: true,
}

var Linkarm = obj.LinkArch{
	//Arch:           sys.ArchARM,
	//Init:           buildop,
	//Preprocess:     preprocess,
	//Assemble:       span5,
	//Progedit:       progedit,
	UnaryDst: unaryDst,
	//DWARFRegisters: ARMDWARFRegisters,
}
