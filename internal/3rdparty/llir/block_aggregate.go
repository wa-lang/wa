package llir

import (
	value "github.com/wa-lang/wa/internal/3rdparty/llir/llvalue"
)

// --- [ Aggregate instructions ] ----------------------------------------------

// ~~~ [ extractvalue ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// NewExtractValue appends a new extractvalue instruction to the basic block
// based on the given aggregate value and indicies.
func (block *Block) NewExtractValue(x value.Value, indices ...uint64) *InstExtractValue {
	inst := NewExtractValue(x, indices...)
	block.Insts = append(block.Insts, inst)
	return inst
}

// ~~~ [ insertvalue ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// NewInsertValue appends a new insertvalue instruction to the basic block based
// on the given aggregate value, element and indicies.
func (block *Block) NewInsertValue(x, elem value.Value, indices ...uint64) *InstInsertValue {
	inst := NewInsertValue(x, elem, indices...)
	block.Insts = append(block.Insts, inst)
	return inst
}
