package llir

import (
	"fmt"
	"strings"

	"github.com/wa-lang/wa/internal/3rdparty/llir/internal/gep"
	constant "github.com/wa-lang/wa/internal/3rdparty/llir/llconstant"
	enum "github.com/wa-lang/wa/internal/3rdparty/llir/llenum"
	types "github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
	value "github.com/wa-lang/wa/internal/3rdparty/llir/llvalue"
)

// --- [ Memory instructions ] -------------------------------------------------

// ~~~ [ alloca ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// InstAlloca is an LLVM IR alloca instruction.
type InstAlloca struct {
	// Name of local variable associated with the result.
	LocalIdent
	// Element type.
	ElemType types.Type
	// (optional) Number of elements; nil if not present.
	NElems value.Value

	// extra.

	// Type of result produced by the instruction, including an optional address
	// space.
	Typ *types.PointerType
	// (optional) In-alloca.
	InAlloca bool
	// (optional) Swift error.
	SwiftError bool
	// (optional) Alignment; zero if not present.
	Align Align
	// (optional) Address space; zero if not present.
	AddrSpace types.AddrSpace
	// (optional) Metadata.
	Metadata
}

// NewAlloca returns a new alloca instruction based on the given element type.
func NewAlloca(elemType types.Type) *InstAlloca {
	inst := &InstAlloca{ElemType: elemType}
	// Compute type.
	inst.Type()
	return inst
}

// String returns the LLVM syntax representation of the instruction as a
// type-value pair.
func (inst *InstAlloca) String() string {
	return fmt.Sprintf("%s %s", inst.Type(), inst.Ident())
}

// Type returns the type of the instruction.
func (inst *InstAlloca) Type() types.Type {
	// Cache type if not present.
	if inst.Typ == nil {
		inst.Typ = types.NewPointer(inst.ElemType)
		inst.Typ.AddrSpace = inst.AddrSpace
	}
	return inst.Typ
}

// LLString returns the LLVM syntax representation of the instruction.
//
// 'alloca' InAllocaopt SwiftErroropt ElemType=Type NElems=(',' TypeValue)? (',' Align)? (',' AddrSpace)? Metadata=(',' MetadataAttachment)+?
func (inst *InstAlloca) LLString() string {
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "%s = ", inst.Ident())
	buf.WriteString("alloca")
	if inst.InAlloca {
		buf.WriteString(" inalloca")
	}
	if inst.SwiftError {
		buf.WriteString(" swifterror")
	}
	fmt.Fprintf(buf, " %s", inst.ElemType)
	if inst.NElems != nil {
		fmt.Fprintf(buf, ", %s", inst.NElems)
	}
	if inst.Align != 0 {
		fmt.Fprintf(buf, ", %s", inst.Align)
	}
	if inst.AddrSpace != 0 {
		fmt.Fprintf(buf, ", %s", inst.AddrSpace)
	}
	for _, md := range inst.Metadata {
		fmt.Fprintf(buf, ", %s", md)
	}
	return buf.String()
}

// ~~~ [ load ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// InstLoad is an LLVM IR load instruction.
type InstLoad struct {
	// Name of local variable associated with the result.
	LocalIdent
	// Element type of src.
	ElemType types.Type
	// Source address.
	Src value.Value

	// extra.

	// (optional) Atomic.
	Atomic bool
	// (optional) Volatile.
	Volatile bool
	// (optional) Sync scope; empty if not present.
	SyncScope string
	// (optional) Atomic memory ordering constraints; zero if not present.
	Ordering enum.AtomicOrdering
	// (optional) Alignment; zero if not present.
	Align Align
	// (optional) Metadata.
	Metadata
}

// NewLoad returns a new load instruction based on the given element type and
// source address.
func NewLoad(elemType types.Type, src value.Value) *InstLoad {
	inst := &InstLoad{ElemType: elemType, Src: src}
	return inst
}

// String returns the LLVM syntax representation of the instruction as a
// type-value pair.
func (inst *InstLoad) String() string {
	return fmt.Sprintf("%s %s", inst.Type(), inst.Ident())
}

// Type returns the type of the instruction.
func (inst *InstLoad) Type() types.Type {
	return inst.ElemType
}

// LLString returns the LLVM syntax representation of the instruction.
//
// Load instruction.
//
//    'load' Volatileopt ElemType=Type ',' Src=TypeValue (',' Align)? Metadata=(',' MetadataAttachment)+?
//
// Atomic load instruction.
//
//    'load' Atomic Volatileopt ElemType=Type ',' Src=TypeValue SyncScopeopt Ordering=AtomicOrdering (',' Align)? Metadata=(',' MetadataAttachment)+?
func (inst *InstLoad) LLString() string {
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "%s = ", inst.Ident())
	buf.WriteString("load")
	if inst.Atomic {
		buf.WriteString(" atomic")
	}
	if inst.Volatile {
		buf.WriteString(" volatile")
	}
	fmt.Fprintf(buf, " %s, %s", inst.ElemType, inst.Src)
	if len(inst.SyncScope) > 0 {
		fmt.Fprintf(buf, " syncscope(%s)", quote(inst.SyncScope))
	}
	if inst.Ordering != enum.AtomicOrderingNone {
		fmt.Fprintf(buf, " %s", inst.Ordering)
	}
	if inst.Align != 0 {
		fmt.Fprintf(buf, ", %s", inst.Align)
	}
	for _, md := range inst.Metadata {
		fmt.Fprintf(buf, ", %s", md)
	}
	return buf.String()
}

// ~~~ [ store ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// InstStore is an LLVM IR store instruction.
type InstStore struct {
	// Source value.
	Src value.Value
	// Destination address.
	Dst value.Value

	// extra.

	// (optional) Atomic.
	Atomic bool
	// (optional) Volatile.
	Volatile bool
	// (optional) Sync scope; empty if not present.
	SyncScope string
	// (optional) Atomic memory ordering constraints; zero if not present.
	Ordering enum.AtomicOrdering
	// (optional) Alignment; zero if not present.
	Align Align
	// (optional) Metadata.
	Metadata
}

// NewStore returns a new store instruction based on the given source value and
// destination address.
func NewStore(src, dst value.Value) *InstStore {
	// Type-check operands.
	dstPtrType, ok := dst.Type().(*types.PointerType)
	if !ok {
		panic(fmt.Errorf("invalid store dst operand type; expected *types.Pointer, got %T", dst.Type()))
	}
	if !src.Type().Equal(dstPtrType.ElemType) {
		panic(fmt.Errorf("store operands are not compatible: src=%v; dst=%v", src.Type(), dst.Type()))
	}
	return &InstStore{Src: src, Dst: dst}
}

// LLString returns the LLVM syntax representation of the instruction.
//
// Store instruction.
//
//    'store' Volatileopt Src=TypeValue ',' Dst=TypeValue (',' Align)? Metadata=(',' MetadataAttachment)+?
//
// Atomic store instruction.
//
//    'store' Atomic Volatileopt Src=TypeValue ',' Dst=TypeValue SyncScopeopt Ordering=AtomicOrdering (',' Align)? Metadata=(',' MetadataAttachment)+?
func (inst *InstStore) LLString() string {
	buf := &strings.Builder{}
	buf.WriteString("store")
	if inst.Atomic {
		buf.WriteString(" atomic")
	}
	if inst.Volatile {
		buf.WriteString(" volatile")
	}
	fmt.Fprintf(buf, " %s, %s", inst.Src, inst.Dst)
	if len(inst.SyncScope) > 0 {
		fmt.Fprintf(buf, " syncscope(%s)", quote(inst.SyncScope))
	}
	if inst.Ordering != enum.AtomicOrderingNone {
		fmt.Fprintf(buf, " %s", inst.Ordering)
	}
	if inst.Align != 0 {
		fmt.Fprintf(buf, ", %s", inst.Align)
	}
	for _, md := range inst.Metadata {
		fmt.Fprintf(buf, ", %s", md)
	}
	return buf.String()
}

// ~~~ [ fence ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// InstFence is an LLVM IR fence instruction.
type InstFence struct {
	// Atomic memory ordering constraints.
	Ordering enum.AtomicOrdering

	// extra.

	// (optional) Sync scope; empty if not present.
	SyncScope string
	// (optional) Metadata.
	Metadata
}

// NewFence returns a new fence instruction based on the given atomic ordering.
func NewFence(ordering enum.AtomicOrdering) *InstFence {
	return &InstFence{Ordering: ordering}
}

// LLString returns the LLVM syntax representation of the instruction.
//
// 'fence' SyncScopeopt Ordering=AtomicOrdering Metadata=(',' MetadataAttachment)+?
func (inst *InstFence) LLString() string {
	buf := &strings.Builder{}
	buf.WriteString("fence")
	if len(inst.SyncScope) > 0 {
		fmt.Fprintf(buf, " syncscope(%s)", quote(inst.SyncScope))
	}
	fmt.Fprintf(buf, " %s", inst.Ordering)
	for _, md := range inst.Metadata {
		fmt.Fprintf(buf, ", %s", md)
	}
	return buf.String()
}

// ~~~ [ cmpxchg ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// InstCmpXchg is an LLVM IR cmpxchg instruction.
type InstCmpXchg struct {
	// Name of local variable associated with the result.
	LocalIdent
	// Address to read from, compare against and store to.
	Ptr value.Value
	// Value to compare against.
	Cmp value.Value
	// New value to store.
	New value.Value
	// Atomic memory ordering constraints on success.
	SuccessOrdering enum.AtomicOrdering
	// Atomic memory ordering constraints on failure.
	FailureOrdering enum.AtomicOrdering

	// extra.

	// Type of result produced by the instruction; the first field of the struct
	// holds the old value, and the second field indicates success.
	Typ *types.StructType
	// (optional) Weak.
	Weak bool
	// (optional) Volatile.
	Volatile bool
	// (optional) Sync scope; empty if not present.
	SyncScope string
	// (optional) Metadata.
	Metadata
}

// NewCmpXchg returns a new cmpxchg instruction based on the given address,
// value to compare against, new value to store, and atomic orderings for
// success and failure.
func NewCmpXchg(ptr, cmp, new value.Value, successOrdering, failureOrdering enum.AtomicOrdering) *InstCmpXchg {
	inst := &InstCmpXchg{Ptr: ptr, Cmp: cmp, New: new, SuccessOrdering: successOrdering, FailureOrdering: failureOrdering}
	// Compute type.
	inst.Type()
	return inst
}

// String returns the LLVM syntax representation of the instruction as a
// type-value pair.
func (inst *InstCmpXchg) String() string {
	return fmt.Sprintf("%s %s", inst.Type(), inst.Ident())
}

// Type returns the type of the instruction. The result type is a struct type
// with two fields, the first field has the type of the old value and the second
// field has boolean type.
func (inst *InstCmpXchg) Type() types.Type {
	// Cache type if not present.
	if inst.Typ == nil {
		oldType := inst.New.Type()
		inst.Typ = types.NewStruct(oldType, types.I1)
	}
	return inst.Typ
}

// LLString returns the LLVM syntax representation of the instruction.
//
// 'cmpxchg' Weakopt Volatileopt Ptr=TypeValue ',' Cmp=TypeValue ',' New=TypeValue SyncScopeopt SuccessOrdering=AtomicOrdering FailureOrdering=AtomicOrdering Metadata=(',' MetadataAttachment)+?
func (inst *InstCmpXchg) LLString() string {
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "%s = ", inst.Ident())
	buf.WriteString("cmpxchg")
	if inst.Weak {
		buf.WriteString(" weak")
	}
	if inst.Volatile {
		buf.WriteString(" volatile")
	}
	fmt.Fprintf(buf, " %s, %s, %s", inst.Ptr, inst.Cmp, inst.New)
	if len(inst.SyncScope) > 0 {
		fmt.Fprintf(buf, " syncscope(%s)", quote(inst.SyncScope))
	}
	fmt.Fprintf(buf, " %s", inst.SuccessOrdering)
	fmt.Fprintf(buf, " %s", inst.FailureOrdering)
	for _, md := range inst.Metadata {
		fmt.Fprintf(buf, ", %s", md)
	}
	return buf.String()
}

// ~~~ [ atomicrmw ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// InstAtomicRMW is an LLVM IR atomicrmw instruction.
type InstAtomicRMW struct {
	// Name of local variable associated with the result.
	LocalIdent
	// Atomic operation.
	Op enum.AtomicOp
	// Destination address.
	Dst value.Value
	// Operand.
	X value.Value
	// Atomic memory ordering constraints.
	Ordering enum.AtomicOrdering

	// extra.

	// Type of result produced by the instruction.
	Typ types.Type
	// (optional) Volatile.
	Volatile bool
	// (optional) Sync scope; empty if not present.
	SyncScope string
	// (optional) Metadata.
	Metadata
}

// NewAtomicRMW returns a new atomicrmw instruction based on the given atomic
// operation, destination address, operand and atomic ordering.
func NewAtomicRMW(op enum.AtomicOp, dst, x value.Value, ordering enum.AtomicOrdering) *InstAtomicRMW {
	inst := &InstAtomicRMW{Op: op, Dst: dst, X: x, Ordering: ordering}
	// Compute type.
	inst.Type()
	return inst
}

// String returns the LLVM syntax representation of the instruction as a
// type-value pair.
func (inst *InstAtomicRMW) String() string {
	return fmt.Sprintf("%s %s", inst.Type(), inst.Ident())
}

// Type returns the type of the instruction.
func (inst *InstAtomicRMW) Type() types.Type {
	// Cache type if not present.
	if inst.Typ == nil {
		t, ok := inst.Dst.Type().(*types.PointerType)
		if !ok {
			panic(fmt.Errorf("invalid destination type; expected *types.PointerType, got %T", inst.Dst.Type()))
		}
		inst.Typ = t.ElemType
	}
	return inst.Typ
}

// LLString returns the LLVM syntax representation of the instruction.
//
// 'atomicrmw' Volatileopt Op=AtomicOp Dst=TypeValue ',' X=TypeValue SyncScopeopt Ordering=AtomicOrdering Metadata=(',' MetadataAttachment)+?
func (inst *InstAtomicRMW) LLString() string {
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "%s = ", inst.Ident())
	buf.WriteString("atomicrmw")
	if inst.Volatile {
		buf.WriteString(" volatile")
	}
	fmt.Fprintf(buf, " %s %s, %s", inst.Op, inst.Dst, inst.X)
	if len(inst.SyncScope) > 0 {
		fmt.Fprintf(buf, " syncscope(%s)", quote(inst.SyncScope))
	}
	fmt.Fprintf(buf, " %s", inst.Ordering)
	for _, md := range inst.Metadata {
		fmt.Fprintf(buf, ", %s", md)
	}
	return buf.String()
}

// ~~~ [ getelementptr ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// InstGetElementPtr is an LLVM IR getelementptr instruction.
type InstGetElementPtr struct {
	// Name of local variable associated with the result.
	LocalIdent
	// Element type.
	ElemType types.Type
	// Source address.
	Src value.Value
	// Element indicies.
	Indices []value.Value

	// extra.

	// Type of result produced by the instruction.
	Typ types.Type // *types.PointerType or *types.VectorType (with elements of pointer type)
	// (optional) In-bounds.
	InBounds bool
	// (optional) Metadata.
	Metadata
}

// NewGetElementPtr returns a new getelementptr instruction based on the given
// element type, source address and element indices.
func NewGetElementPtr(elemType types.Type, src value.Value, indices ...value.Value) *InstGetElementPtr {
	inst := &InstGetElementPtr{ElemType: elemType, Src: src, Indices: indices}
	// Compute type.
	inst.Type()
	return inst
}

// String returns the LLVM syntax representation of the instruction as a
// type-value pair.
func (inst *InstGetElementPtr) String() string {
	return fmt.Sprintf("%s %s", inst.Type(), inst.Ident())
}

// Type returns the type of the instruction.
func (inst *InstGetElementPtr) Type() types.Type {
	// Cache type if not present.
	if inst.Typ == nil {
		inst.Typ = gepInstType(inst.ElemType, inst.Src.Type(), inst.Indices)
	}
	return inst.Typ
}

// LLString returns the LLVM syntax representation of the instruction.
//
// 'getelementptr' InBoundsopt ElemType=Type ',' Src=TypeValue Indices=(',' TypeValue)* Metadata=(',' MetadataAttachment)+?
func (inst *InstGetElementPtr) LLString() string {
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "%s = ", inst.Ident())
	buf.WriteString("getelementptr")
	if inst.InBounds {
		buf.WriteString(" inbounds")
	}
	fmt.Fprintf(buf, " %s, %s", inst.ElemType, inst.Src)
	for _, index := range inst.Indices {
		fmt.Fprintf(buf, ", %s", index)
	}
	for _, md := range inst.Metadata {
		fmt.Fprintf(buf, ", %s", md)
	}
	return buf.String()
}

// ### [ Helper functions ] ####################################################

// gepInstType computes the result type of a getelementptr instruction.
//
//    getelementptr ElemType, Src, Indices
func gepInstType(elemType, src types.Type, indices []value.Value) types.Type {
	var idxs []gep.Index
	for _, index := range indices {
		var idx gep.Index
		switch index := index.(type) {
		case constant.Constant:
			idx = getIndex(index)
		default:
			idx = gep.Index{HasVal: false}
			// Check if index is of vector type.
			if indexType, ok := index.Type().(*types.VectorType); ok {
				idx.VectorLen = indexType.Len
			}
		}
		idxs = append(idxs, idx)
	}
	return gep.ResultType(elemType, src, idxs)
}

// NOTE: keep getIndex in sync with getIndex in:
//
//    * ast/inst_memory.go
//    * ir/inst_memory.go
//    * ir/constant/expr_memory.go
//
// The reference point and source of truth is in ir/constant/expr_memory.go.

// getIndex returns the gep index corresponding to the given constant index.
func getIndex(index constant.Constant) gep.Index {
	// unpack inrange indices.
	if idx, ok := index.(*constant.Index); ok {
		index = idx.Constant
	}
	// TODO: figure out how to simplify expressions for GEP instructions without
	// creating import cycle on irutil.

	// Use index.Simplify() to simplify the constant expression to a concrete
	// integer constant or vector of integers constant.
	//if idx, ok := index.(constant.Expression); ok {
	//	index = idx.Simplify()
	//}
	switch index := index.(type) {
	case *constant.Int:
		val := index.X.Int64()
		return gep.NewIndex(val)
	case *constant.ZeroInitializer:
		return gep.NewIndex(0)
	case *constant.Vector:
		// ref: https://llvm.org/docs/LangRef.html#getelementptr-instruction
		//
		// > The getelementptr returns a vector of pointers, instead of a single
		// > address, when one or more of its arguments is a vector. In such
		// > cases, all vector arguments should have the same number of elements,
		// > and every scalar argument will be effectively broadcast into a vector
		// > during address calculation.
		if len(index.Elems) == 0 {
			return gep.Index{HasVal: false}
		}
		// Sanity check. All vector elements must be integers, and must have the
		// same value.
		var val int64
		for i, elem := range index.Elems {
			switch elem := elem.(type) {
			case *constant.Int:
				x := elem.X.Int64()
				if i == 0 {
					val = x
				} else if x != val {
					// since all elements were not identical, we must conclude that
					// the index vector does not have a concrete value.
					return gep.Index{
						HasVal:    false,
						VectorLen: uint64(len(index.Elems)),
					}
				}
			default:
				// TODO: remove debug output.
				panic(fmt.Errorf("support for gep index vector element type %T not yet implemented", elem))
				//return gep.Index{HasVal: false}
			}
		}
		return gep.Index{
			HasVal:    true,
			Val:       val,
			VectorLen: uint64(len(index.Elems)),
		}
	case *constant.Undef:
		return gep.Index{HasVal: false}
	case *constant.Poison:
		return gep.Index{HasVal: false}
	case constant.Expression:
		// should already have been simplified to a form we can handle.
		return gep.Index{HasVal: false}
	default:
		// TODO: add support for more constant expressions.
		// TODO: remove debug output.
		panic(fmt.Errorf("support for gep index type %T not yet implemented", index))
		//return gep.Index{HasVal: false}
	}
}
