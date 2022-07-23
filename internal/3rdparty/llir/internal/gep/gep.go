// Package gep computes the result type of getelementptr instructions and
// constant expressions.
//
// ref: https://llvm.org/docs/GetElementPtr.html
// ref: https://llvm.org/docs/LangRef.html#getelementptr-instruction
package gep

import (
	"fmt"
	"log"

	types "github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
)

// Index is a gep index.
type Index struct {
	// HasVal specifies whether Val has a valid value. If index is a constant
	// integer or a constant integer vector of which all elements have the same
	// value, then HasVal is set. Note, this is a requirement to index into
	// structure types.
	HasVal bool
	// Index integer value. Val is only valid if HasVal is set.
	Val int64
	// Length of index vector; or 0 if index is scalar. VectorLen may be non-zero
	// even if HasVal is false.
	VectorLen uint64
}

// NewIndex returns a new constant index with the given value.
func NewIndex(val int64) Index {
	return Index{
		HasVal: true,
		Val:    val,
	}
}

// ResultType computes the result type of a getelementptr instruction or
// constant expression.
//
//    getelementptr (ElemType, Src, Indices)
func ResultType(elemType, src types.Type, indices []Index) types.Type {
	// ref: http://llvm.org/docs/GetElementPtr.html#what-effect-do-address-spaces-have-on-geps
	//
	// > the address space qualifier on the second operand pointer type always
	// > matches the address space qualifier on the result type.
	var (
		// Address space of src pointer type or src vector element pointer type.
		addrSpace types.AddrSpace
		// Length of vector of pointers result type; or 0 if pointer result type.
		resultVectorLength uint64
	)
	// ref: https://llvm.org/docs/LangRef.html#getelementptr-instruction
	//
	// > The second argument is always a pointer or a vector of pointers.
	switch src := src.(type) {
	case *types.PointerType:
		addrSpace = src.AddrSpace
	case *types.VectorType:
		vectorElemType, ok := src.ElemType.(*types.PointerType)
		if !ok {
			panic(fmt.Errorf("invalid gep source vector element type; expected *types.PointerType, got %T", src.ElemType))
		}
		addrSpace = vectorElemType.AddrSpace
		resultVectorLength = src.Len
	default:
		panic(fmt.Errorf("invalid gep source type; expected pointer or vector of pointers type, got %T", src))
	}
	// ref: https://llvm.org/docs/LangRef.html#getelementptr-instruction
	//
	// > The first argument is always a type used as the basis for the
	// > calculations.
	e := elemType
	for i, index := range indices {
		// ref: https://llvm.org/docs/LangRef.html#getelementptr-instruction
		//
		// > The getelementptr returns a vector of pointers, instead of a single
		// > address, when one or more of its arguments is a vector. In such
		// > cases, all vector arguments should have the same number of elements,
		// > and every scalar argument will be effectively broadcast into a vector
		// > during address calculation.
		if index.VectorLen != 0 && resultVectorLength != 0 && index.VectorLen != resultVectorLength {
			panic(fmt.Errorf("vector length mismatch of index vector (%d) and result type vector (%d)", index.VectorLen, resultVectorLength))
		}
		if resultVectorLength == 0 && index.VectorLen != 0 {
			resultVectorLength = index.VectorLen
		}
		// ref: https://llvm.org/docs/GetElementPtr.html#why-is-the-extra-0-index-required
		//
		// > Since the second argument to the GEP instruction must always be a
		// > value of pointer type, the first index steps through that pointer.
		if i == 0 {
			continue
		}
		switch elm := e.(type) {
		case *types.PointerType:
			panic(fmt.Errorf("cannot index into pointer type at %d:th gep index, only valid at 0:th gep index; see https://llvm.org/docs/GetElementPtr.html#what-is-dereferenced-by-gep", i))
		case *types.VectorType:
			// ref: https://llvm.org/docs/GetElementPtr.html#can-gep-index-into-vector-elements
			//
			// > This hasn’t always been forcefully disallowed, though it’s not
			// > recommended. It leads to awkward special cases in the optimizers,
			// > and fundamental inconsistency in the IR. In the future, it will
			// > probably be outright disallowed.
			log.Printf("using gep to index into vector types will be disallowed in a future release or llir/llvm; see https://llvm.org/docs/GetElementPtr.html#can-gep-index-into-vector-elements")
			e = elm.ElemType
		case *types.ArrayType:
			e = elm.ElemType
		case *types.StructType:
			// ref: https://llvm.org/docs/LangRef.html#getelementptr-instruction
			//
			// > When indexing into a (optionally packed) structure, only i32
			// > integer constants are allowed (when using a vector of indices they
			// > must all be the same i32 integer constant).
			if !index.HasVal {
				panic(fmt.Errorf("unable to index into struct type `%v` using gep with non-constant index", e))
			}
			e = elm.Fields[index.Val]
		default:
			panic(fmt.Errorf("cannot index into type %T using gep", e))
		}
	}
	ptr := types.NewPointer(e)
	ptr.AddrSpace = addrSpace
	if resultVectorLength != 0 {
		vec := types.NewVector(resultVectorLength, ptr)
		return vec
	}
	return ptr
}
