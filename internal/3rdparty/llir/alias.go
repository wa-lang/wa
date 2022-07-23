package llir

import (
	"fmt"
	"strings"

	constant "github.com/wa-lang/wa/internal/3rdparty/llir/llconstant"
	enum "github.com/wa-lang/wa/internal/3rdparty/llir/llenum"
	types "github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
)

// === [ Aliases ] =============================================================

// Alias is an alias of a global identifier or constant expression.
type Alias struct {
	// Alias name (without '@' prefix).
	GlobalIdent
	// Aliasee.
	Aliasee constant.Constant

	// Pointer type of aliasee.
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
	// (optional) Partition name; empty if not present.
	Partition string
}

// NewAlias returns a new alias based on the given alias name and aliasee.
func NewAlias(name string, aliasee constant.Constant) *Alias {
	alias := &Alias{Aliasee: aliasee}
	alias.SetName(name)
	// Compute type.
	alias.Type()
	return alias
}

// String returns the LLVM syntax representation of the alias as a type-value
// pair.
func (a *Alias) String() string {
	return fmt.Sprintf("%s %s", a.Type(), a.Ident())
}

// Type returns the type of the alias.
func (a *Alias) Type() types.Type {
	// Cache type if not present.
	if a.Typ == nil {
		typ, ok := a.Aliasee.Type().(*types.PointerType)
		if !ok {
			panic(fmt.Errorf("invalid aliasee type of %q; expected *types.PointerType, got %T", a.Ident(), a.Aliasee.Type()))
		}
		a.Typ = typ
	}
	return a.Typ
}

// LLString returns the LLVM syntax representation of the alias definition.
//
// Name=GlobalIdent '=' (ExternLinkage | Linkageopt) Preemptionopt Visibilityopt DLLStorageClassopt ThreadLocalopt UnnamedAddropt IndirectSymbolKind ContentType=Type ',' IndirectSymbol Partitions=(',' Partition)*
func (a *Alias) LLString() string {
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "%s =", a.Ident())
	if a.Linkage != enum.LinkageNone {
		fmt.Fprintf(buf, " %s", a.Linkage)
	}
	if a.Preemption != enum.PreemptionNone {
		fmt.Fprintf(buf, " %s", a.Preemption)
	}
	if a.Visibility != enum.VisibilityNone {
		fmt.Fprintf(buf, " %s", a.Visibility)
	}
	if a.DLLStorageClass != enum.DLLStorageClassNone {
		fmt.Fprintf(buf, " %s", a.DLLStorageClass)
	}
	if a.TLSModel != enum.TLSModelNone {
		fmt.Fprintf(buf, " %s", tlsModelString(a.TLSModel))
	}
	if a.UnnamedAddr != enum.UnnamedAddrNone {
		fmt.Fprintf(buf, " %s", a.UnnamedAddr)
	}
	buf.WriteString(" alias")
	fmt.Fprintf(buf, " %s, ", a.Typ.ElemType)
	if expr, ok := a.Aliasee.(constant.Expression); ok {
		buf.WriteString(expr.Ident())
	} else {
		buf.WriteString(a.Aliasee.String())
	}
	if len(a.Partition) > 0 {
		fmt.Fprintf(buf, ", partition %s", quote(a.Partition))
	}
	return buf.String()
}
