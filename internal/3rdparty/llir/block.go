// The basic block API of the ir package was heavily inspired by
// https://github.com/mrbenshef/goory.

package llir

import (
	"fmt"
	"strings"

	"github.com/wa-lang/wa/internal/3rdparty/llir/internal/enc"
	types "github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
)

// === [ Basic blocks ] ========================================================

// Block is an LLVM IR basic block; a sequence of non-branching instructions
// terminated by a control flow instruction (e.g. br or ret).
type Block struct {
	// Name of local variable associated with the basic block.
	LocalIdent
	// Instructions of the basic block.
	Insts []Instruction
	// Terminator of the basic block.
	Term Terminator

	// extra.

	// Parent function; field set by ir.Func.NewBlock.
	Parent *Func `json:"-"`
}

// NewBlock returns a new basic block based on the given label name. An empty
// label name indicates an unnamed basic block.
func NewBlock(name string) *Block {
	block := &Block{}
	block.SetName(name)
	return block
}

// String returns the LLVM syntax representation of the basic block as a
// type-value pair.
func (block *Block) String() string {
	return fmt.Sprintf("%s %s", block.Type(), block.Ident())
}

// Type returns the type of the basic block.
func (block *Block) Type() types.Type {
	return types.Label
}

// LLString returns the LLVM syntax representation of the basic block
// definition.
//
// Name=LabelIdentopt Insts=Instruction* Term=Terminator
func (block *Block) LLString() string {
	buf := &strings.Builder{}
	if block.IsUnnamed() {
		//fmt.Fprintf(buf, "; <label>:%d\n", block.LocalID)
		// Explicitly print basic block label to conform with Clang 9.0, and
		// because it's the sane thing to do.
		fmt.Fprintf(buf, "%s\n", enc.LabelID(block.LocalID))
	} else {
		fmt.Fprintf(buf, "%s\n", enc.LabelName(block.LocalName))
	}
	for _, inst := range block.Insts {
		fmt.Fprintf(buf, "\t%s\n", inst.LLString())
	}
	if block.Term == nil {
		panic(fmt.Sprintf("missing terminator in basic block %q.\ncurrent instructions:\n%s", block.Name(), buf.String()))
	}
	fmt.Fprintf(buf, "\t%s", block.Term.LLString())
	return buf.String()
}
