package llir

import (
	"fmt"
	"strings"
	"sync"

	"github.com/wa-lang/wa/internal/3rdparty/llir/internal/enc"
	constant "github.com/wa-lang/wa/internal/3rdparty/llir/llconstant"
	enum "github.com/wa-lang/wa/internal/3rdparty/llir/llenum"
	types "github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
)

// === [ Functions ] ===========================================================

// Func is an LLVM IR function. The body of a function definition consists of a
// set of basic blocks, interconnected by terminator control flow instructions.
type Func struct {
	// Function name (without '@' prefix).
	GlobalIdent
	// Function signature.
	Sig *types.FuncType
	// Function parameters.
	Params []*Param
	// Basic blocks.
	Blocks []*Block // nil if declaration.

	// extra.

	// Pointer type to function, including an optional address space. If Typ is
	// nil, the first invocation of Type stores a pointer type with Sig as
	// element.
	Typ *types.PointerType
	// (optional) Linkage.
	Linkage enum.Linkage
	// (optional) Preemption; zero value if not present.
	Preemption enum.Preemption
	// (optional) Visibility; zero value if not present.
	Visibility enum.Visibility
	// (optional) DLL storage class; zero value if not present.
	DLLStorageClass enum.DLLStorageClass
	// (optional) Calling convention; zero value if not present.
	CallingConv enum.CallingConv
	// (optional) Return attributes.
	ReturnAttrs []ReturnAttribute
	// (optional) Unnamed address.
	UnnamedAddr enum.UnnamedAddr
	// (optional) Address space; zero if not present.
	AddrSpace types.AddrSpace
	// (optional) Function attributes.
	FuncAttrs []FuncAttribute
	// (optional) Section name; empty if not present.
	Section string
	// (optional) Partition name; empty if not present.
	Partition string
	// (optional) Comdat; nil if not present.
	Comdat *ComdatDef
	// (optional) Alignment; zero if not present.
	Align Align
	// (optional) Garbage collection; empty if not present.
	GC string
	// (optional) Prefix; nil if not present.
	Prefix constant.Constant
	// (optional) Prologue; nil if not present.
	Prologue constant.Constant
	// (optional) Personality; nil if not present.
	Personality constant.Constant
	// (optional) Use list orders.
	UseListOrders []*UseListOrder
	// (optional) Metadata.
	Metadata

	// Parent module; field set by ir.Module.NewFunc.
	Parent *Module `json:"-"`

	// mu prevents races on AssignIDs.
	mu sync.Mutex
}

// NewFunc returns a new function based on the given function name, return type
// and function parameters.
func NewFunc(name string, retType types.Type, params ...*Param) *Func {
	paramTypes := make([]types.Type, len(params))
	for i, param := range params {
		paramTypes[i] = param.Type()
	}
	sig := types.NewFunc(retType, paramTypes...)
	f := &Func{Sig: sig, Params: params}
	f.SetName(name)
	// Compute type.
	f.Type()
	return f
}

// String returns the LLVM syntax representation of the function as a type-value
// pair.
func (f *Func) String() string {
	return fmt.Sprintf("%s %s", f.Type(), f.Ident())
}

// Type returns the type of the function.
func (f *Func) Type() types.Type {
	// Cache type if not present.
	if f.Typ == nil {
		f.Typ = types.NewPointer(f.Sig)
		f.Typ.AddrSpace = f.AddrSpace
	}
	return f.Typ
}

// LLString returns the LLVM syntax representation of the function definition or
// declaration.
//
// Function declaration.
//
//    'declare' Metadata=MetadataAttachment* Header=FuncHeader
//
// Function definition.
//
//    'define' Header=FuncHeader Metadata=MetadataAttachment* Body=FuncBody
func (f *Func) LLString() string {
	if err := f.AssignIDs(); err != nil {
		panic(fmt.Errorf("unable to assign IDs of function %q; %v", f.Ident(), err))
	}
	buf := &strings.Builder{}
	if len(f.Blocks) == 0 {
		// Function declaration.
		buf.WriteString("declare")
		for _, md := range f.Metadata {
			fmt.Fprintf(buf, " %s", md)
		}
		if f.Linkage != enum.LinkageNone {
			fmt.Fprintf(buf, " %s", f.Linkage)
		}
		buf.WriteString(headerString(f))
		return buf.String()
	} else {
		// Function definition.
		buf.WriteString("define")
		if f.Linkage != enum.LinkageNone {
			fmt.Fprintf(buf, " %s", f.Linkage)
		}
		buf.WriteString(headerString(f))
		for _, md := range f.Metadata {
			fmt.Fprintf(buf, " %s", md)
		}
		fmt.Fprintf(buf, " %s", bodyString(f))
		return buf.String()
	}
}

// AssignIDs assigns IDs to unnamed local variables.
func (f *Func) AssignIDs() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	id := int64(0)
	setName := func(n namedVar) error {
		if n.IsUnnamed() {
			if n.ID() != 0 && id != n.ID() {
				want := id
				got := n.ID()
				return fmt.Errorf("invalid local ID in function %q, expected %s, got %s", f.Ident(), enc.LocalID(want), enc.LocalID(got))
			}
			n.SetID(id)
			id++
		}
		return nil
	}
	for _, param := range f.Params {
		// Assign local IDs to unnamed parameters of function definitions.
		if err := setName(param); err != nil {
			return err
		}
	}
	for _, block := range f.Blocks {
		// Assign local IDs to unnamed basic blocks.
		if err := setName(block); err != nil {
			return err
		}
		for _, inst := range block.Insts {
			n, ok := inst.(namedVar)
			if !ok {
				continue
			}
			// Skip void instructions (call with void return).
			if types.Equal(n.Type(), types.Void) {
				continue
			}
			// Assign local IDs to unnamed local variables.
			if err := setName(n); err != nil {
				return err
			}
		}
		n, ok := block.Term.(namedVar)
		if !ok {
			continue
		}
		// Skip void terminators (invoke, callbr with void return).
		if types.Equal(n.Type(), types.Void) {
			continue
		}
		if err := setName(n); err != nil {
			return err
		}
	}
	return nil
}

// ### [ Helper functions ] ####################################################

// headerString returns the string representation of the function header.
func headerString(f *Func) string {
	// (Linkage | ExternLinkage)? Preemptionopt Visibilityopt DLLStorageClassopt CallingConvopt ReturnAttrs=ReturnAttribute* RetType=Type Name=GlobalIdent '(' Params ')' UnnamedAddropt AddrSpaceopt FuncAttrs=FuncAttribute* Sectionopt Partitionopt Comdatopt Alignopt GCopt Prefixopt Prologueopt Personalityopt
	buf := &strings.Builder{}
	if f.Preemption != enum.PreemptionNone {
		fmt.Fprintf(buf, " %s", f.Preemption)
	}
	if f.Visibility != enum.VisibilityNone {
		fmt.Fprintf(buf, " %s", f.Visibility)
	}
	if f.DLLStorageClass != enum.DLLStorageClassNone {
		fmt.Fprintf(buf, " %s", f.DLLStorageClass)
	}
	if f.CallingConv != enum.CallingConvNone {
		fmt.Fprintf(buf, " %s", callingConvString(f.CallingConv))
	}
	for _, attr := range f.ReturnAttrs {
		fmt.Fprintf(buf, " %s", attr)
	}
	fmt.Fprintf(buf, " %s", f.Sig.RetType)
	fmt.Fprintf(buf, " %s(", f.Ident())
	for i, param := range f.Params {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(param.LLString())
	}
	if f.Sig.Variadic {
		if len(f.Params) > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString("...")
	}
	buf.WriteString(")")
	if f.UnnamedAddr != enum.UnnamedAddrNone {
		fmt.Fprintf(buf, " %s", f.UnnamedAddr)
	}
	if f.AddrSpace != 0 {
		fmt.Fprintf(buf, " %s", f.AddrSpace)
	}
	for _, attr := range f.FuncAttrs {
		fmt.Fprintf(buf, " %s", attr)
	}
	if len(f.Section) > 0 {
		fmt.Fprintf(buf, " section %s", quote(f.Section))
	}
	if len(f.Partition) > 0 {
		fmt.Fprintf(buf, " partition %s", quote(f.Partition))
	}
	if f.Comdat != nil {
		if f.Comdat.Name == f.Name() {
			buf.WriteString(" comdat")
		} else {
			fmt.Fprintf(buf, " %s", f.Comdat)
		}
	}
	if f.Align != 0 {
		fmt.Fprintf(buf, " %s", f.Align)
	}
	if len(f.GC) > 0 {
		fmt.Fprintf(buf, " gc %s", quote(f.GC))
	}
	if f.Prefix != nil {
		fmt.Fprintf(buf, " prefix %s", f.Prefix)
	}
	if f.Prologue != nil {
		fmt.Fprintf(buf, " prologue %s", f.Prologue)
	}
	if f.Personality != nil {
		fmt.Fprintf(buf, " personality %s", f.Personality)
	}
	return buf.String()
}

// bodyString returns the string representation of the function body.
func bodyString(body *Func) string {
	// '{' Blocks=Block+ UseListOrders=UseListOrder* '}'
	buf := &strings.Builder{}
	buf.WriteString("{\n")
	for i, block := range body.Blocks {
		if i != 0 {
			buf.WriteString("\n")
		}
		fmt.Fprintf(buf, "%s\n", block.LLString())
	}
	if len(body.UseListOrders) > 0 {
		buf.WriteString("\n")
	}
	for _, u := range body.UseListOrders {
		fmt.Fprintf(buf, "\t%s\n", u)
	}
	buf.WriteString("}")
	return buf.String()
}
