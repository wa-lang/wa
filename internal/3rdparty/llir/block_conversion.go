package llir

import (
	types "github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
	value "github.com/wa-lang/wa/internal/3rdparty/llir/llvalue"
)

// --- [ Conversion instructions ] ---------------------------------------------

// ~~~ [ trunc ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// NewTrunc appends a new trunc instruction to the basic block based on the
// given source value and target type.
func (block *Block) NewTrunc(from value.Value, to types.Type) *InstTrunc {
	inst := NewTrunc(from, to)
	block.Insts = append(block.Insts, inst)
	return inst
}

// ~~~ [ zext ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// NewZExt appends a new zext instruction to the basic block based on the given
// source value and target type.
func (block *Block) NewZExt(from value.Value, to types.Type) *InstZExt {
	inst := NewZExt(from, to)
	block.Insts = append(block.Insts, inst)
	return inst
}

// ~~~ [ sext ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// NewSExt appends a new sext instruction to the basic block based on the given
// source value and target type.
func (block *Block) NewSExt(from value.Value, to types.Type) *InstSExt {
	inst := NewSExt(from, to)
	block.Insts = append(block.Insts, inst)
	return inst
}

// ~~~ [ fptrunc ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// NewFPTrunc appends a new fptrunc instruction to the basic block based on the
// given source value and target type.
func (block *Block) NewFPTrunc(from value.Value, to types.Type) *InstFPTrunc {
	inst := NewFPTrunc(from, to)
	block.Insts = append(block.Insts, inst)
	return inst
}

// ~~~ [ fpext ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// NewFPExt appends a new fpext instruction to the basic block based on the
// given source value and target type.
func (block *Block) NewFPExt(from value.Value, to types.Type) *InstFPExt {
	inst := NewFPExt(from, to)
	block.Insts = append(block.Insts, inst)
	return inst
}

// ~~~ [ fptoui ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// NewFPToUI appends a new fptoui instruction to the basic block based on the
// given source value and target type.
func (block *Block) NewFPToUI(from value.Value, to types.Type) *InstFPToUI {
	inst := NewFPToUI(from, to)
	block.Insts = append(block.Insts, inst)
	return inst
}

// ~~~ [ fptosi ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// NewFPToSI appends a new fptosi instruction to the basic block based on the
// given source value and target type.
func (block *Block) NewFPToSI(from value.Value, to types.Type) *InstFPToSI {
	inst := NewFPToSI(from, to)
	block.Insts = append(block.Insts, inst)
	return inst
}

// ~~~ [ uitofp ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// NewUIToFP appends a new uitofp instruction to the basic block based on the
// given source value and target type.
func (block *Block) NewUIToFP(from value.Value, to types.Type) *InstUIToFP {
	inst := NewUIToFP(from, to)
	block.Insts = append(block.Insts, inst)
	return inst
}

// ~~~ [ sitofp ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// NewSIToFP appends a new sitofp instruction to the basic block based on the
// given source value and target type.
func (block *Block) NewSIToFP(from value.Value, to types.Type) *InstSIToFP {
	inst := NewSIToFP(from, to)
	block.Insts = append(block.Insts, inst)
	return inst
}

// ~~~ [ ptrtoint ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// NewPtrToInt appends a new ptrtoint instruction to the basic block based on
// the given source value and target type.
func (block *Block) NewPtrToInt(from value.Value, to types.Type) *InstPtrToInt {
	inst := NewPtrToInt(from, to)
	block.Insts = append(block.Insts, inst)
	return inst
}

// ~~~ [ inttoptr ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// NewIntToPtr appends a new inttoptr instruction to the basic block based on
// the given source value and target type.
func (block *Block) NewIntToPtr(from value.Value, to types.Type) *InstIntToPtr {
	inst := NewIntToPtr(from, to)
	block.Insts = append(block.Insts, inst)
	return inst
}

// ~~~ [ bitcast ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// NewBitCast appends a new bitcast instruction to the basic block based on the
// given source value and target type.
func (block *Block) NewBitCast(from value.Value, to types.Type) *InstBitCast {
	inst := NewBitCast(from, to)
	block.Insts = append(block.Insts, inst)
	return inst
}

// ~~~ [ addrspacecast ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// NewAddrSpaceCast appends a new addrspacecast instruction to the basic block
// based on the given source value and target type.
func (block *Block) NewAddrSpaceCast(from value.Value, to types.Type) *InstAddrSpaceCast {
	inst := NewAddrSpaceCast(from, to)
	block.Insts = append(block.Insts, inst)
	return inst
}
