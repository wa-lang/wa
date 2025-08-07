package arm

import (
	"wa-lang.org/wa/internal/p9asm/obj"
	"wa-lang.org/wa/internal/p9asm/objabi"
)

var unaryDst = map[objabi.As]bool{
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
