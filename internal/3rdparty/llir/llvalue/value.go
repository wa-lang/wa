// Package value provides a definition of LLVM IR values.
package llvalue

import (
	"fmt"

	types "github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
)

// Value is an LLVM IR value, which may be used as an operand of instructions
// and terminators.
//
// A Value has one of the following underlying types.
//
//    constant.Constant   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llconstant#Constant
//    value.Named         // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/value#Named
//    TODO: add literal metadata value?
type Value interface {
	// String returns the LLVM syntax representation of the value as a type-value
	// pair.
	fmt.Stringer
	// Type returns the type of the value.
	Type() types.Type
	// Ident returns the identifier associated with the value.
	Ident() string
}

// Named is a named LLVM IR value.
//
// A Named value has one of the following underlying types.
//
//    *ir.Global            // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#Global
//    *ir.Func              // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#Func
//    *ir.Param             // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#Param
//    *ir.Block             // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#Block
//    TODO: add named metadata value?
//    ir.Instruction        // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#Instruction (except store and fence)
//    *ir.TermInvoke        // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#TermInvoke
//    *ir.TermCatchSwitch   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#TermCatchSwitch (token result used by catchpad)
type Named interface {
	Value
	// Name returns the name of the value.
	Name() string
	// SetName sets the name of the value.
	SetName(name string)
}
