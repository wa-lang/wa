package llutil

import (
	"github.com/wa-lang/wa/internal/3rdparty/llir"
	value "github.com/wa-lang/wa/internal/3rdparty/llir/llvalue"
)

// Ident is the global or local identifier of a named value.
type Ident interface {
	value.Named
	// ID returns the ID of the identifier.
	ID() int64
	// SetID sets the ID of the identifier.
	SetID(id int64)
	// IsUnnamed reports whether the identifier is unnamed.
	IsUnnamed() bool
}

// ResetNames resets the IDs of unnamed local variables in the given function.
func ResetNames(f *llir.Func) {
	for _, param := range f.Params {
		// clear ID of unnamed function parameter.
		if param.IsUnnamed() {
			param.SetName("")
		}
	}
	for _, block := range f.Blocks {
		// clear ID of unnamed basic block.
		if block.IsUnnamed() {
			block.SetName("")
		}
		for _, inst := range block.Insts {
			if inst, ok := inst.(Ident); ok {
				// clear ID of unnamed variable.
				if inst.IsUnnamed() {
					inst.SetName("")
				}
			}
		}
	}
}
