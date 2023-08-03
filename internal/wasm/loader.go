package wasm

// Sym encapsulates a global symbol index, used to identify a specific
// Go symbol. The 0-valued Sym is corresponds to an invalid symbol.
type Sym int

// objSym represents a symbol in an object file. It is a tuple of
// the object and the symbol's local index.
// For external symbols, objidx is the index of l.extReader (extObj),
// s is its index into the payload array.
// {0, 0} represents the nil symbol.
type objSym struct {
	objidx uint32 // index of the object (in l.objs array)
	s      uint32 // local index
}

type nameVer struct {
	name string
	v    int
}

type Loader struct {
	sects    []*Section // sections
	symSects []uint16   // symbol's section, index to sects array

	symsByName    [2]map[string]Sym // map symbol name to index, two maps are for ABI0 and ABIInternal
	extStaticSyms map[nameVer]Sym   // externally defined static symbols, keyed by name

	values []int64 // symbol values, indexed by global sym index
}

// SymValue returns the section of the i-th symbol. i is global index.
func (l *Loader) SymSect(i Sym) *Section {
	if int(i) >= len(l.symSects) {
		// symSects is extended lazily -- it the sym in question is
		// outside the range of the existing slice, then we assume its
		// section has not yet been set.
		return nil
	}
	return l.sects[l.symSects[i]]
}

// Look up a symbol by name, return global index, or 0 if not found.
// This is more like Syms.ROLookup than Lookup -- it doesn't create
// new symbol.
func (l *Loader) Lookup(name string, ver int) Sym {
	return l.symsByName[ver][name]
}

// SymValue returns the value of the i-th symbol. i is global index.
func (l *Loader) SymValue(i Sym) int64 {
	return l.values[i]
}
