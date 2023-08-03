package wasm

// Link holds the context for writing object code from a compiler
// or for reading that input into the linker.
type ld_Link struct {
	Out *OutBuf
}

// Used only on Wasm for now.
func ld_DatblkBytes(ctxt *ld_Link, addr int64, size int64) []byte {
	buf := make([]byte, size)
	out := &OutBuf{heap: buf}
	writeDatblkToOutBuf(ctxt, out, addr, size)
	return buf
}

func writeDatblkToOutBuf(ctxt *ld_Link, out *OutBuf, addr int64, size int64) {
	//writeBlocks(ctxt, out, ctxt.outSem, ctxt.loader, ctxt.datap, addr, size, zeros[:])
}
