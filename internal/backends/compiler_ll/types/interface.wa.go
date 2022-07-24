package types

import (
	"fmt"
	"sort"

	"github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
	"github.com/wa-lang/wa/internal/3rdparty/llir/llvalue"
)

type Interface struct {
	backingType

	SourceName      string
	RequiredMethods map[string]InterfaceMethod
}

func (i Interface) Name() string {
	return fmt.Sprintf("interface(%s)", i.SourceName)
}

// SortedRequiredMethods returns a sorted slice of all method names
// The returned order is the order the methods will be layed out in the JumpTable
func (i Interface) SortedRequiredMethods() []string {
	var orderedMethods []string
	for methodName := range i.RequiredMethods {
		orderedMethods = append(orderedMethods, methodName)
	}
	sort.Strings(orderedMethods)
	return orderedMethods
}

func (i Interface) JumpTable() *lltypes.StructType {
	orderedMethods := i.SortedRequiredMethods()

	var ifaceTableMethods []lltypes.Type

	for _, methodName := range orderedMethods {
		methodSignature := i.RequiredMethods[methodName]

		var retType lltypes.Type = lltypes.Void
		if len(methodSignature.ReturnTypes) > 0 {
			retType = methodSignature.ReturnTypes[0].LLVM()
		}

		paramTypes := []lltypes.Type{lltypes.NewPointer(lltypes.I8)}
		for _, argType := range methodSignature.ArgumentTypes {
			paramTypes = append(paramTypes, argType.LLVM())
		}

		ifaceTableMethods = append(ifaceTableMethods, lltypes.NewPointer(lltypes.NewFunc(retType, paramTypes...)))
	}

	return lltypes.NewStruct(ifaceTableMethods...)
}

func (i Interface) LLVM() lltypes.Type {
	return lltypes.NewStruct(
		// Pointer to the backing data
		lltypes.NewPointer(lltypes.I8),

		// Backing data type
		lltypes.I32,

		// Interface table
		// Used for method resolving
		lltypes.NewPointer(i.JumpTable()),
	)
}

func (Interface) Size() int64 {
	return 64 / 8 * 3
}

type InterfaceMethod struct {
	backingType

	LlvmJumpFunction llvalue.Named

	ArgumentTypes []Type
	ReturnTypes   []Type
}

func (InterfaceMethod) LLVM() lltypes.Type {
	panic("InterfaceMethod has no LLVM value")
}

func (InterfaceMethod) Name() string {
	return "InterfaceMethod"
}
