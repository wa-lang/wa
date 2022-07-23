package llir

import (
	"fmt"
	"strings"

	constant "github.com/wa-lang/wa/internal/3rdparty/llir/llconstant"
	enum "github.com/wa-lang/wa/internal/3rdparty/llir/llenum"
	types "github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
)

// === [ Global variables ] ====================================================

// Global is a global variable declaration or definition.
type Global struct {
	// Global variable name (without '@' prefix).
	GlobalIdent
	// Immutability of global variable (constant or global).
	Immutable bool
	// Content type.
	ContentType types.Type
	// Initial value; or nil if declaration.
	Init constant.Constant

	// extra.

	// Pointer type to global variable, including an optional address space. If
	// Typ is nil, the first invocation of Type stores a pointer type with
	// ContentType as element.
	Typ *types.PointerType
	// (optional) Linkage; zero value if not present.
	Linkage enum.Linkage
	// (optional) Preemption; zero value if not present.
	Preemption enum.Preemption
	// (optional) Visibility; zero value if not present.
	Visibility enum.Visibility
	// (optional) DLL storage class; zero value if not present.
	DLLStorageClass enum.DLLStorageClass
	// (optional) Thread local storage model; zero value if not present.
	TLSModel enum.TLSModel
	// (optional) Unnamed address; zero value if not present.
	UnnamedAddr enum.UnnamedAddr
	// (optional) Address space; zero if not present.
	AddrSpace types.AddrSpace
	// (optional) Externally initialized; false if not present.
	ExternallyInitialized bool
	// (optional) Section name; empty if not present.
	Section string
	// (optional) Partition name; empty if not present.
	Partition string
	// (optional) Comdat; nil if not present.
	Comdat *ComdatDef
	// (optional) Alignment; zero if not present.
	Align Align
	// (optional) Function attributes.
	FuncAttrs []FuncAttribute
	// (optional) Metadata.
	Metadata
}

// NewGlobal returns a new global variable declaration based on the given global
// variable name and content type.
func NewGlobal(name string, contentType types.Type) *Global {
	global := &Global{ContentType: contentType}
	global.SetName(name)
	// Compute type.
	global.Type()
	return global
}

// NewGlobalDef returns a new global variable definition based on the given
// global variable name and initial value.
func NewGlobalDef(name string, init constant.Constant) *Global {
	global := &Global{ContentType: init.Type(), Init: init}
	global.SetName(name)
	// Compute type.
	global.Type()
	return global
}

// String returns the LLVM syntax representation of the global variable as a
// type-value pair.
func (g *Global) String() string {
	return fmt.Sprintf("%s %s", g.Type(), g.Ident())
}

// Type returns the type of the global variable.
func (g *Global) Type() types.Type {
	// Cache type if not present.
	if g.Typ == nil {
		g.Typ = types.NewPointer(g.ContentType)
		g.Typ.AddrSpace = g.AddrSpace
	}
	return g.Typ
}

// LLString returns the LLVM syntax representation of the global variable
// definition or declaration.
//
// Global declaration.
//
//    Name=GlobalIdent '=' Linkage=ExternLinkage Preemptionopt Visibilityopt DLLStorageClassopt ThreadLocalopt UnnamedAddropt AddrSpaceopt ExternallyInitializedopt Immutable ContentType=Type (',' Section)? (',' Partition)? (',' Comdat)? (',' Align)? Metadata=(',' MetadataAttachment)+? FuncAttrs=FuncAttribute+?
//
// Global definition.
//
//    Name=GlobalIdent '=' Linkage=Linkageopt Preemptionopt Visibilityopt DLLStorageClassopt ThreadLocalopt UnnamedAddropt AddrSpaceopt ExternallyInitializedopt Immutable ContentType=Type Init=Constant (',' Section)? (',' Partition)? (',' Comdat)? (',' Align)? Metadata=(',' MetadataAttachment)+? FuncAttrs=FuncAttribute+?
func (g *Global) LLString() string {
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "%s =", g.Ident())
	if g.Linkage != enum.LinkageNone {
		fmt.Fprintf(buf, " %s", g.Linkage)
	}
	if g.Preemption != enum.PreemptionNone {
		fmt.Fprintf(buf, " %s", g.Preemption)
	}
	if g.Visibility != enum.VisibilityNone {
		fmt.Fprintf(buf, " %s", g.Visibility)
	}
	if g.DLLStorageClass != enum.DLLStorageClassNone {
		fmt.Fprintf(buf, " %s", g.DLLStorageClass)
	}
	if g.TLSModel != enum.TLSModelNone {
		fmt.Fprintf(buf, " %s", tlsModelString(g.TLSModel))
	}
	if g.UnnamedAddr != enum.UnnamedAddrNone {
		fmt.Fprintf(buf, " %s", g.UnnamedAddr)
	}
	if g.AddrSpace != 0 {
		fmt.Fprintf(buf, " %s", g.AddrSpace)
	}
	if g.ExternallyInitialized {
		buf.WriteString(" externally_initialized")
	}
	if g.Immutable {
		buf.WriteString(" constant")
	} else {
		buf.WriteString(" global")
	}
	fmt.Fprintf(buf, " %s", g.ContentType)
	if g.Init != nil {
		// Global definition.
		fmt.Fprintf(buf, " %s", g.Init.Ident())
	}
	if g.Section != "" {
		fmt.Fprintf(buf, ", section %s", quote(g.Section))
	}
	if g.Partition != "" {
		fmt.Fprintf(buf, ", partition %s", quote(g.Partition))
	}
	if g.Comdat != nil {
		if g.Comdat.Name == g.Name() {
			buf.WriteString(", comdat")
		} else {
			fmt.Fprintf(buf, ", %s", g.Comdat)
		}
	}
	if g.Align != 0 {
		fmt.Fprintf(buf, ", %s", g.Align)
	}
	for _, md := range g.Metadata {
		fmt.Fprintf(buf, ", %s", md)
	}
	for _, attr := range g.FuncAttrs {
		fmt.Fprintf(buf, " %s", attr)
	}
	return buf.String()
}
