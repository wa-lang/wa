package syscall

import (
	"github.com/wa-lang/wa/internal/3rdparty/llir"
	"github.com/wa-lang/wa/internal/3rdparty/llir/llconstant"
	"github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
	"github.com/wa-lang/wa/internal/3rdparty/llir/llvalue"
	"github.com/wa-lang/wa/internal/backends/compiler_ll/llutil"
)

func Print(block *llir.Block, value llvalue.Value, goos string) {
	asmFunc := llir.NewInlineAsm(
		lltypes.NewPointer(lltypes.NewFunc(lltypes.I64)),
		"syscall", "=r,{rax},{rdi},{rsi},{rdx}",
	)
	asmFunc.SideEffect = true

	strPtr := llutil.StrToI8Ptr(block, value)
	strLen := llutil.StrLen(block, value)

	block.NewCall(asmFunc,
		llconstant.NewInt(lltypes.I64, Convert(WRITE, goos)), // rax
		llconstant.NewInt(lltypes.I64, 1),                    // rdi, stdout
		strPtr,                                               // rsi
		strLen,                                               // rdx
	)
}
