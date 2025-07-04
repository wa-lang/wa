package compiler

import (
	"wa-lang.org/wa/internal/3rdparty/wazero/internal/asm/amd64"
	"wa-lang.org/wa/internal/3rdparty/wazero/internal/wazeroir"
)

// init initializes variables for the amd64 architecture
func init() {
	newArchContext = newArchContextImpl
	registerNameFn = amd64.RegisterName
	unreservedGeneralPurposeRegisters = amd64UnreservedGeneralPurposeRegisters
	unreservedVectorRegisters = amd64UnreservedVectorRegisters
}

// archContext is embedded in callEngine in order to store architecture-specific data.
// For amd64, this is empty.
type archContext struct{}

// newArchContextImpl implements newArchContext for amd64 architecture.
func newArchContextImpl() (ret archContext) { return }

// newCompiler returns a new compiler interface which can be used to compile the given function instance.
// Note: ir param can be nil for host functions.
func newCompiler(ir *wazeroir.CompilationResult, withListener bool) (compiler, error) {
	return newAmd64Compiler(ir, withListener)
}
