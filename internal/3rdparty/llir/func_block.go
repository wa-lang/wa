package llir

// NewBlock appends a new basic block to the function based on the given label
// name. An empty label name indicates an unnamed basic block.
//
// The Parent field of the block is set to f.
func (f *Func) NewBlock(name string) *Block {
	block := NewBlock(name)
	block.Parent = f
	f.Blocks = append(f.Blocks, block)
	return block
}
