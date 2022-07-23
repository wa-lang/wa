package llir

import (
	"fmt"
	"strings"

	constant "github.com/wa-lang/wa/internal/3rdparty/llir/llconstant"
	enum "github.com/wa-lang/wa/internal/3rdparty/llir/llenum"
	types "github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
)

// === [ IFuncs ] ==============================================================

// IFunc is an indirect function (a special kind of function alias).
type IFunc struct {
	// IFunc name (without '@' prefix).
	GlobalIdent
	// Resolver.
	Resolver constant.Constant

	// Pointer type of resolver.
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

// NewIFunc returns a new indirect function based on the given IFunc name and
// resolver.
func NewIFunc(name string, resolver constant.Constant) *IFunc {
	ifunc := &IFunc{Resolver: resolver}
	ifunc.SetName(name)
	// Compute type.
	ifunc.Type()
	return ifunc
}

// String returns the LLVM syntax representation of the IFunc as a type-value
// pair.
func (i *IFunc) String() string {
	return fmt.Sprintf("%s %s", i.Type(), i.Ident())
}

// Type returns the type of the IFunc.
func (i *IFunc) Type() types.Type {
	// Cache type if not present.
	if i.Typ == nil {
		typ, ok := i.Resolver.Type().(*types.PointerType)
		if !ok {
			panic(fmt.Errorf("invalid resolver type of %q; expected *types.PointerType, got %T", i.Ident(), i.Resolver.Type()))
		}
		i.Typ = typ
	}
	return i.Typ
}

// LLString returns the LLVM syntax representation of the IFunc definition.
//
// Name=GlobalIdent '=' (ExternLinkage | Linkageopt) Preemptionopt Visibilityopt DLLStorageClassopt ThreadLocalopt UnnamedAddropt IndirectSymbolKind ContentType=Type ',' IndirectSymbol Partitions=(',' Partition)*
func (i *IFunc) LLString() string {
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "%s =", i.Ident())
	if i.Linkage != enum.LinkageNone {
		fmt.Fprintf(buf, " %s", i.Linkage)
	}
	if i.Preemption != enum.PreemptionNone {
		fmt.Fprintf(buf, " %s", i.Preemption)
	}
	if i.Visibility != enum.VisibilityNone {
		fmt.Fprintf(buf, " %s", i.Visibility)
	}
	if i.DLLStorageClass != enum.DLLStorageClassNone {
		fmt.Fprintf(buf, " %s", i.DLLStorageClass)
	}
	if i.TLSModel != enum.TLSModelNone {
		fmt.Fprintf(buf, " %s", tlsModelString(i.TLSModel))
	}
	if i.UnnamedAddr != enum.UnnamedAddrNone {
		fmt.Fprintf(buf, " %s", i.UnnamedAddr)
	}
	buf.WriteString(" ifunc")
	fmt.Fprintf(buf, " %s, %s", i.Typ.ElemType, i.Resolver)
	if len(i.Partition) > 0 {
		fmt.Fprintf(buf, ", partition %s", quote(i.Partition))
	}
	return buf.String()
}
