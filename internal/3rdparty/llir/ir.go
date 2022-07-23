// Package ir declares the types used to represent LLVM IR modules.
package llir

// LLStringer is implemented by any value that has a LLString method, which
// defines the LLVM syntax for that value.
type LLStringer interface {
	// LLString returns the LLVM syntax representation of the value.
	LLString() string
}
