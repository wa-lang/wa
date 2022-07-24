package llir

import (
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/wa-lang/wa/internal/3rdparty/errors"
	"github.com/wa-lang/wa/internal/3rdparty/llir/internal/enc"
	"github.com/wa-lang/wa/internal/3rdparty/llir/internal/natsort"
	enum "github.com/wa-lang/wa/internal/3rdparty/llir/llenum"
	metadata "github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata"
	types "github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
	value "github.com/wa-lang/wa/internal/3rdparty/llir/llvalue"
)

// === [ Modules ] =============================================================

// Module is an LLVM IR module, which consists of top-level declarations and
// definitions.
type Module struct {
	// Type definitions.
	TypeDefs []types.Type
	// Global variable declarations and definitions.
	Globals []*Global
	// Function declarations and definitions.
	Funcs []*Func

	// extra.

	// (optional) Source filename; or empty if not present.
	SourceFilename string
	// (optional) Data layout; or empty if not present.
	DataLayout string
	// (optional) Target triple; or empty if not present.
	TargetTriple string
	// (optional) Module-level inline assembly.
	ModuleAsms []string
	// (optional) Comdat definitions.
	ComdatDefs []*ComdatDef
	// (optional) Aliases.
	Aliases []*Alias
	// (optional) IFuncs.
	IFuncs []*IFunc
	// (optional) Attribute group definitions.
	AttrGroupDefs []*AttrGroupDef
	// (optional) Named metadata definitions.
	NamedMetadataDefs map[string]*metadata.NamedDef
	// (optional) Metadata definitions.
	MetadataDefs []metadata.Definition
	// (optional) Use-list order directives.
	UseListOrders []*UseListOrder
	// (optional) Basic block specific use-list order directives.
	UseListOrderBBs []*UseListOrderBB

	// mu prevents races on AssignGlobalIDs and AssignMetadataIDs.
	mu sync.Mutex
}

// NewModule returns a new LLVM IR module.
func NewModule() *Module {
	return &Module{
		NamedMetadataDefs: make(map[string]*metadata.NamedDef),
	}
}

// String returns the string representation of the module in LLVM IR assembly
// syntax.
func (m *Module) String() string {
	buf := &strings.Builder{}
	if _, err := m.WriteTo(buf); err != nil {
		panic(fmt.Errorf("unable to write to string buffer; %v", err))
	}
	return buf.String()
}

// WriteTo write the string representation of the module in LLVM IR assembly
// syntax to w.
func (m *Module) WriteTo(w io.Writer) (n int64, err error) {
	fw := &fmtWriter{w: w}
	// Assign global IDs.
	if err := m.AssignGlobalIDs(); err != nil {
		panic(fmt.Errorf("unable to assign globals IDs of module; %v", err))
	}
	// Assign metadata IDs.
	if err := m.AssignMetadataIDs(); err != nil {
		panic(fmt.Errorf("unable to assign metadata IDs of module; %v", err))
	}
	// Source filename.
	if len(m.SourceFilename) > 0 {
		// 'source_filename' '=' Name=StringLit
		fw.Fprintf("source_filename = %s\n", quote(m.SourceFilename))
	}
	// Data layout.
	if len(m.DataLayout) > 0 {
		// 'target' 'datalayout' '=' DataLayout=StringLit
		fw.Fprintf("target datalayout = %s\n", quote(m.DataLayout))
	}
	// Target triple.
	if len(m.TargetTriple) > 0 {
		// 'target' 'triple' '=' TargetTriple=StringLit
		fw.Fprintf("target triple = %s\n", quote(m.TargetTriple))
	}
	// Module-level inline assembly.
	if len(m.ModuleAsms) > 0 && fw.size > 0 {
		fw.Fprint("\n")
	}
	for _, asm := range m.ModuleAsms {
		// 'module' 'asm' Asm=StringLit
		fw.Fprintf("module asm %s\n", quote(asm))
	}
	// Type definitions.
	if len(m.TypeDefs) > 0 && fw.size > 0 {
		fw.Fprint("\n")
	}
	for _, t := range m.TypeDefs {
		// Name=LocalIdent '=' 'type' Typ=OpaqueType
		//
		// Name=LocalIdent '=' 'type' Typ=Type
		fw.Fprintf("%s = type %s\n", t, t.LLString())
	}
	// Comdat definitions.
	if len(m.ComdatDefs) > 0 && fw.size > 0 {
		fw.Fprint("\n")
	}
	for _, def := range m.ComdatDefs {
		fw.Fprintln(def.LLString())
	}
	// Global declarations and definitions.
	if len(m.Globals) > 0 && fw.size > 0 {
		fw.Fprint("\n")
	}
	for _, g := range m.Globals {
		fw.Fprintln(g.LLString())
	}
	// Aliases.
	if len(m.Aliases) > 0 && fw.size > 0 {
		fw.Fprint("\n")
	}
	for _, alias := range m.Aliases {
		fw.Fprintln(alias.LLString())
	}
	// IFuncs.
	if len(m.IFuncs) > 0 && fw.size > 0 {
		fw.Fprint("\n")
	}
	for _, ifunc := range m.IFuncs {
		fw.Fprintln(ifunc.LLString())
	}
	// Function declarations and definitions.
	if len(m.Funcs) > 0 && fw.size > 0 {
		fw.Fprint("\n")
	}
	for i, f := range m.Funcs {
		if i != 0 {
			fw.Fprint("\n")
		}
		fw.Fprintln(f.LLString())
	}
	// Attribute group definitions.
	if len(m.AttrGroupDefs) > 0 && fw.size > 0 {
		fw.Fprint("\n")
	}
	for _, a := range m.AttrGroupDefs {
		fw.Fprintln(a.LLString())
	}
	// Named metadata definitions; output in natural sorting order.
	var mdNames []string
	for mdName := range m.NamedMetadataDefs {
		mdNames = append(mdNames, mdName)
	}
	natsort.Strings(mdNames)
	if len(m.NamedMetadataDefs) > 0 && fw.size > 0 {
		fw.Fprint("\n")
	}
	for _, mdName := range mdNames {
		// Name=MetadataName '=' '!' '{' MDNodes=(MetadataNode separator ',')* '}'
		md := m.NamedMetadataDefs[mdName]
		fw.Fprintf("%s = %s\n", md.Ident(), md.LLString())
	}
	// Metadata definitions.
	if len(m.MetadataDefs) > 0 && fw.size > 0 {
		fw.Fprint("\n")
	}
	for _, md := range m.MetadataDefs {
		// ID=MetadataID '=' Distinctopt MDNode=MDTuple
		//
		// ID=MetadataID '=' Distinctopt MDNode=SpecializedMDNode
		fw.Fprintf("%s = %s\n", md.Ident(), md.LLString())
	}
	// Use-list orders.
	if len(m.UseListOrders) > 0 && fw.size > 0 {
		fw.Fprint("\n")
	}
	for _, u := range m.UseListOrders {
		fw.Fprintln(u)
	}
	// Basic block specific use-list orders.
	if len(m.UseListOrderBBs) > 0 && fw.size > 0 {
		fw.Fprint("\n")
	}
	for _, u := range m.UseListOrderBBs {
		fw.Fprintln(u)
	}
	return fw.size, fw.err
}

// ~~~ [ Comdat Definition ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ComdatDef is a comdat definition top-level entity.
type ComdatDef struct {
	// Comdat name (without '$' prefix).
	Name string
	// Comdat kind.
	Kind enum.SelectionKind
}

// String returns the string representation of the Comdat definition.
func (c *ComdatDef) String() string {
	return fmt.Sprintf("comdat(%s)", enc.ComdatName(c.Name))
}

// LLString returns the LLVM syntax representation of the Comdat definition.
//
// Name=ComdatName '=' 'comdat' Kind=SelectionKind
func (c *ComdatDef) LLString() string {
	return fmt.Sprintf("%s = comdat %s", enc.ComdatName(c.Name), c.Kind)
}

// ~~~ [ Attribute Group Definition ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// AttrGroupDef is an attribute group definition.
type AttrGroupDef struct {
	// Attribute group ID (without '#' prefix).
	ID int64
	// Function attributes.
	FuncAttrs []FuncAttribute
}

// String returns the string representation of the attribute group definition.
func (a *AttrGroupDef) String() string {
	return enc.AttrGroupID(a.ID)
}

// LLString returns the LLVM syntax representation of the attribute group
// definition.
//
// 'attributes' ID=AttrGroupID '=' '{' FuncAttrs=FuncAttribute* '}'
func (a *AttrGroupDef) LLString() string {
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "attributes %s = { ", enc.AttrGroupID(a.ID))
	for i, attr := range a.FuncAttrs {
		if i != 0 {
			buf.WriteString(" ")
		}
		switch attr := attr.(type) {
		case Align:
			// Note, alignment is printed as `align = 8` in attribute groups.
			fmt.Fprintf(buf, "align = %d", uint64(attr))
		case AlignStack:
			// Note, stack alignment is printed as `alignstack = 8` in attribute
			// groups.
			fmt.Fprintf(buf, "alignstack = %d", uint64(attr))
		default:
			buf.WriteString(attr.String())
		}
	}
	buf.WriteString(" }")
	return buf.String()
}

// ~~~ [ Use-list Order Directives ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// UseListOrder is a use-list order directive.
type UseListOrder struct {
	// Value.
	Value value.Value
	// Use-list order.
	Indices []uint64
}

// String returns the string representation of the use-list order directive
// definition.
func (u *UseListOrder) String() string {
	//  'uselistorder' TypeValue ',' '{' Indices=(UintLit separator ',')+ '}'
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "uselistorder %s, { ", u.Value)
	for i, index := range u.Indices {
		if i != 0 {
			buf.WriteString(", ")
		}
		fmt.Fprintf(buf, "%d", index)
	}
	buf.WriteString(" }")
	return buf.String()
}

// UseListOrderBB is a basic block specific use-list order directive.
type UseListOrderBB struct {
	// Function.
	Func *Func
	// Basic block.
	Block *Block
	// Use-list order.
	Indices []uint64
}

// String returns the string representation of the basic block specific use-
// list order directive definition.
func (u *UseListOrderBB) String() string {
	//  'uselistorder_bb' Func=GlobalIdent ',' Block=LocalIdent ',' '{'
	//  Indices=(UintLit separator ',')+ '}'
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "uselistorder_bb %s, %s, { ", u.Func.Ident(), u.Block.Ident())
	for i, index := range u.Indices {
		if i != 0 {
			buf.WriteString(", ")
		}
		fmt.Fprintf(buf, "%d", index)
	}
	buf.WriteString(" }")
	return buf.String()
}

// ### [ Helper functions ] ####################################################

// AssignGlobalIDs assigns IDs to unnamed global variables.
func (m *Module) AssignGlobalIDs() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	id := int64(0)
	setName := func(n namedVar) error {
		if n.IsUnnamed() {
			if n.ID() != 0 && id != n.ID() {
				want := id
				got := n.ID()
				return errors.Errorf("invalid global ID, expected %s, got %s", enc.GlobalID(want), enc.GlobalID(got))
			}
			n.SetID(id)
			id++
		}
		return nil
	}
	// Assign global IDs to unnamed global variables.
	for _, n := range m.Globals {
		if err := setName(n); err != nil {
			return errors.WithStack(err)
		}
	}
	// Assign global IDs to unnamed aliases.
	for _, n := range m.Aliases {
		if err := setName(n); err != nil {
			return errors.WithStack(err)
		}
	}
	// Assign global IDs to unnamed IFuncs.
	for _, n := range m.IFuncs {
		if err := setName(n); err != nil {
			return errors.WithStack(err)
		}
	}
	// Assign global IDs to unnamed functions.
	for _, n := range m.Funcs {
		if err := setName(n); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// AssignMetadataIDs assigns metadata IDs to the unnamed metadata definitions of
// the module.
func (m *Module) AssignMetadataIDs() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	// Index used IDs.
	used := make(map[int64]bool)
	for _, md := range m.MetadataDefs {
		id := md.ID()
		if id != -1 {
			if _, ok := used[id]; ok {
				return errors.Errorf("metadata ID %s already in use", enc.MetadataID(id))
			}
			used[id] = true
		}
	}
	// nextID returns the next unused metdata ID.
	curID := int64(-1)
	nextID := func() int64 {
		for {
			curID++
			if !used[curID] {
				return curID
			}
		}
	}
	// Assign IDs to unnamed metadata definitions.
	for _, md := range m.MetadataDefs {
		id := md.ID()
		if id != -1 {
			// Metadata definition already has ID.
			continue
		}
		newID := nextID()
		md.SetID(newID)
	}
	return nil
}
