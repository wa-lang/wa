package llir

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/wa-lang/wa/internal/3rdparty/llir/internal/enc"
	enum "github.com/wa-lang/wa/internal/3rdparty/llir/llenum"
	metadata "github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata"
	types "github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
	value "github.com/wa-lang/wa/internal/3rdparty/llir/llvalue"
)

// Align is a memory alignment attribute.
type Align uint64

// String returns the string representation of the alignment attribute.
func (align Align) String() string {
	// Note, alignment is printed as `align = 8` in attribute groups.
	return fmt.Sprintf("align %d", uint64(align))
}

// AlignStack is a stack alignment attribute.
type AlignStack uint64

// String returns the string representation of the stack alignment attribute.
func (align AlignStack) String() string {
	// Note, stack alignment is printed as `alignstack = 8` in attribute groups.
	return fmt.Sprintf("alignstack(%d)", uint64(align))
}

// AllocSize is an attribute for functions like malloc. If the second parameter
// is omitted, NElemsIndex will be -1.
type AllocSize struct {
	// Element size parameter index.
	ElemSizeIndex int
	// Number of elements parameter index; -1 if not present.
	NElemsIndex int
}

// String returns the string representation of the allocsize attribute.
func (a AllocSize) String() string {
	if a.NElemsIndex == -1 {
		return fmt.Sprintf("allocsize(%d)", a.ElemSizeIndex)
	}
	return fmt.Sprintf("allocsize(%d, %d)", a.ElemSizeIndex, a.NElemsIndex)
}

// Arg is a function argument with optional parameter attributes.
type Arg struct {
	// Argument value.
	value.Value
	// (optional) Parameter attributes.
	Attrs []ParamAttribute
}

// NewArg returns a new function argument based on the given value and parameter
// attributes.
func NewArg(x value.Value, attrs ...ParamAttribute) *Arg {
	return &Arg{Value: x, Attrs: attrs}
}

// String returns a string representation of the function argument.
func (arg *Arg) String() string {
	// Typ=ConcreteType Attrs=ParamAttribute* Val=Value
	buf := &strings.Builder{}
	buf.WriteString(arg.Type().String())
	for _, attr := range arg.Attrs {
		fmt.Fprintf(buf, " %s", attr)
	}
	fmt.Fprintf(buf, " %s", arg.Ident())
	return buf.String()
}

// AttrPair is an attribute key-value pair (used in function, parameter and
// return attributes).
type AttrPair struct {
	Key   string
	Value string
}

// String returns the string representation of the attribute key-value pair.
func (a AttrPair) String() string {
	// Key=StringLit '=' Val=StringLit
	return fmt.Sprintf("%s=%s", quote(a.Key), quote(a.Value))
}

// AttrString is an attribute string (used in function, parameter and return
// attributes).
type AttrString string

// String returns the string representation of the attribute string.
func (a AttrString) String() string {
	return quote(string(a))
}

// Byval is a byval parameter attribute.
type Byval struct {
	// (optional) Parameter type.
	Typ types.Type
}

// String returns the string representation of the byval parameter attribute.
func (b Byval) String() string {
	// 'byval'
	//
	// 'byval' '(' Typ=Type ')'
	if b.Typ != nil {
		return fmt.Sprintf("byval(%s)", b.Typ)
	}
	return "byval"
}

// Dereferenceable is a dereferenceable memory attribute.
type Dereferenceable struct {
	// Number of bytes known to be dereferenceable.
	N uint64
	// (optional) Either dereferenceable or null if set.
	DerefOrNull bool
}

// String returns the string representation of the dereferenceable memory
// attribute.
func (d Dereferenceable) String() string {
	// 'dereferenceable' '(' N=UintLit ')'
	//
	// 'dereferenceable_or_null' '(' N=UintLit ')'
	if d.DerefOrNull {
		return fmt.Sprintf("dereferenceable_or_null(%d)", d.N)
	}
	return fmt.Sprintf("dereferenceable(%d)", d.N)
}

// Preallocated is a func/param attribute.
type Preallocated struct {
	Typ types.Type
}

// String returns a string representation of the Preallocated attribute.
func (p Preallocated) String() string {
	return fmt.Sprintf("preallocated (%v)", p.Typ)
}

// TODO: check if *ir.InstLandingPad is a valid ExceptionPad.

// ExceptionPad is an exception pad or the none token.
//
// An ExceptionPad has one of the following underlying types.
//
//    *ir.InstCatchPad
//    *ir.InstCleanupPad
//    *constant.NoneToken
type ExceptionPad interface {
	value.Value
}

// FuncAttribute is a function attribute.
//
// A FuncAttribute has one of the following underlying types.
//
//    ir.AttrString
//    ir.AttrPair
//    *ir.AttrGroupDef
//    ir.Align
//    ir.AlignStack
//    ir.AllocSize
//    enum.FuncAttr
type FuncAttribute interface {
	fmt.Stringer
	// IsFuncAttribute ensures that only function attributes can be assigned to
	// the ir.FuncAttribute interface.
	IsFuncAttribute()
}

// GlobalIdent is a global identifier.
type GlobalIdent struct {
	GlobalName string
	GlobalID   int64
}

// Ident returns the identifier associated with the global identifier.
func (i GlobalIdent) Ident() string {
	if i.IsUnnamed() {
		return enc.GlobalID(i.GlobalID)
	}
	return enc.GlobalName(i.GlobalName)
}

// Name returns the name of the global identifier.
//
// If unnamed, the global ID is returned. To distinguish numeric names from
// unnamed IDs, numeric names are quoted.
func (i GlobalIdent) Name() string {
	if i.IsUnnamed() {
		return strconv.FormatInt(i.GlobalID, 10)
	}
	if x, err := strconv.ParseInt(i.GlobalName, 10, 64); err == nil {
		// Print GlobalName with quotes if it is a number; e.g. "42".
		return fmt.Sprintf(`"%d"`, x)
	}
	return i.GlobalName
}

// SetName sets the name of the global identifier.
func (i *GlobalIdent) SetName(name string) {
	i.GlobalName = name
	i.GlobalID = 0
}

// ID returns the ID of the global identifier.
func (i GlobalIdent) ID() int64 {
	return i.GlobalID
}

// SetID sets the ID of the global identifier.
func (i *GlobalIdent) SetID(id int64) {
	i.GlobalID = id
}

// IsUnnamed reports whether the global identifier is unnamed.
func (i GlobalIdent) IsUnnamed() bool {
	return len(i.GlobalName) == 0
}

// LocalIdent is a local identifier.
type LocalIdent struct {
	LocalName string
	LocalID   int64
}

// NewLocalIdent returns a new local identifier based on the given string. An
// unnamed local ID is used if ident is an integer, and a local name otherwise.
func NewLocalIdent(ident string) LocalIdent {
	if id, err := strconv.ParseInt(ident, 10, 64); err == nil {
		return LocalIdent{LocalID: id}
	}
	return LocalIdent{LocalName: ident}
}

// Ident returns the identifier associated with the local identifier.
func (i LocalIdent) Ident() string {
	if i.IsUnnamed() {
		return enc.LocalID(i.LocalID)
	}
	return enc.LocalName(i.LocalName)
}

// Name returns the name of the local identifier.
//
// If unnamed, the local ID is returned. To distinguish numeric names from
// unnamed IDs, numeric names are quoted.
func (i LocalIdent) Name() string {
	if i.IsUnnamed() {
		return strconv.FormatInt(i.LocalID, 10)
	}
	if x, err := strconv.ParseInt(i.LocalName, 10, 64); err == nil {
		// Print LocalName with quotes if it is a number; e.g. "42".
		return fmt.Sprintf(`"%d"`, x)
	}
	return i.LocalName
}

// SetName sets the name of the local identifier.
func (i *LocalIdent) SetName(name string) {
	i.LocalName = name
	i.LocalID = 0
}

// ID returns the ID of the local identifier.
func (i LocalIdent) ID() int64 {
	return i.LocalID
}

// SetID sets the ID of the local identifier.
func (i *LocalIdent) SetID(id int64) {
	i.LocalID = id
}

// IsUnnamed reports whether the local identifier is unnamed.
func (i LocalIdent) IsUnnamed() bool {
	return len(i.LocalName) == 0
}

// Metadata is a list of metadata attachments.
type Metadata []*metadata.Attachment

// MDAttachments returns the metadata attachments of the value.
func (mds Metadata) MDAttachments() []*metadata.Attachment {
	return mds
}

// OperandBundle is a tagged set of SSA values associated with a call-site.
type OperandBundle struct {
	Tag    string
	Inputs []value.Value
}

// NewOperandBundle returns a new operand bundle based on the given tag and
// input values.
func NewOperandBundle(tag string, inputs ...value.Value) *OperandBundle {
	return &OperandBundle{Tag: tag, Inputs: inputs}
}

// String returns a string representation of the operand bundle.
func (o *OperandBundle) String() string {
	// Tag=StringLit '(' Inputs=(TypeValue separator ',')* ')'
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "%s(", quote(o.Tag))
	for i, input := range o.Inputs {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(input.String())
	}
	buf.WriteString(")")
	return buf.String()
}

// ParamAttribute is a parameter attribute.
//
// A ParamAttribute has one of the following underlying types.
//
//    ir.AttrString
//    ir.AttrPair
//    ir.Align
//    ir.Dereferenceable
//    enum.ParamAttr
type ParamAttribute interface {
	fmt.Stringer
	// IsParamAttribute ensures that only parameter attributes can be assigned to
	// the ir.ParamAttribute interface.
	IsParamAttribute()
}

// ReturnAttribute is a return attribute.
//
// A ReturnAttribute has one of the following underlying types.
//
//    ir.AttrString
//    ir.AttrPair
//    ir.Align
//    ir.Dereferenceable
//    enum.ReturnAttr
type ReturnAttribute interface {
	fmt.Stringer
	// IsReturnAttribute ensures that only return attributes can be assigned to
	// the ir.ReturnAttribute interface.
	IsReturnAttribute()
}

// ___ [ Function parameter ] __________________________________________________

// Param is an LLVM IR function parameter.
type Param struct {
	// (optional) Parameter name (without '%' prefix).
	LocalIdent
	// Parameter type.
	Typ types.Type

	// extra.

	// (optional) Parameter attributes.
	Attrs []ParamAttribute
}

// NewParam returns a new function parameter based on the given name and type.
func NewParam(name string, typ types.Type) *Param {
	return &Param{
		LocalIdent: LocalIdent{LocalName: name},
		Typ:        typ,
	}
}

// String returns the LLVM syntax representation of the function parameter as a
// type-value pair.
func (p *Param) String() string {
	return fmt.Sprintf("%s %s", p.Type(), p.Ident())
}

// Type returns the type of the function parameter.
func (p *Param) Type() types.Type {
	return p.Typ
}

// LLString returns the LLVM syntax representation of the function parameter.
//
// Typ=Type Attrs=ParamAttribute* Name=LocalIdent?
func (p *Param) LLString() string {
	buf := &strings.Builder{}
	buf.WriteString(p.Typ.String())
	for _, attr := range p.Attrs {
		fmt.Fprintf(buf, " %s", attr)
	}
	fmt.Fprintf(buf, " %s", p.Ident())
	return buf.String()
}

// SRet is an sret parameter attribute.
type SRet struct {
	Typ types.Type
}

// String returns the string representation of the sret parameter attribute.
func (s SRet) String() string {
	// 'sret' '(' Typ=Type ')'
	return fmt.Sprintf("sret(%s)", s.Typ)
}

// ### [ Helper functions ] ####################################################

// quote returns s as a double-quoted string literal.
func quote(s string) string {
	return enc.Quote([]byte(s))
}

// callingConvString returns the string representation of the given calling
// convention.
func callingConvString(callingConv enum.CallingConv) string {
	s := callingConv.String()
	cc := uint(callingConv)
	if unknown := fmt.Sprintf("CallingConv(%d)", cc); s == unknown {
		return fmt.Sprintf("cc %d", cc)
	}
	return s
}

// tlsModelString returns the string representation of the given thread local
// storage model.
func tlsModelString(model enum.TLSModel) string {
	if model == enum.TLSModelGeneric {
		return "thread_local"
	}
	return fmt.Sprintf("thread_local(%s)", model)
}

// --- [ Formatted I/O writer ] ------------------------------------------------

// fmtWriter is a formatted I/O writer.
//
// A formatted I/O writer keeps track of the total number of bytes written to w
// and the first non-nil error encountered.
type fmtWriter struct {
	// underlying io.Writer.
	w io.Writer
	// Number of bytes written to w.
	size int64
	// First non-nil error encountered.
	err error
}

// Fprint formats using the default formats for its operands and writes to w.
// Spaces are added between operands when neither is a string. It returns the
// number of bytes written and any write error encountered.
func (fw *fmtWriter) Fprint(a ...interface{}) (n int, err error) {
	if fw.err != nil {
		// early return if a previous error has been encountered.
		return 0, nil
	}
	n, err = fmt.Fprint(fw.w, a...)
	fw.size += int64(n)
	fw.err = err
	return n, err
}

// Fprintf formats according to a format specifier and writes to w. It returns
// the number of bytes written and any write error encountered.
func (fw *fmtWriter) Fprintf(format string, a ...interface{}) (n int, err error) {
	if fw.err != nil {
		// early return if a previous error has been encountered.
		return 0, nil
	}
	n, err = fmt.Fprintf(fw.w, format, a...)
	fw.size += int64(n)
	fw.err = err
	return n, err
}

// Fprintln formats using the default formats for its operands and writes to w.
// Spaces are always added between operands and a newline is appended. It
// returns the number of bytes written and any write error encountered.
func (fw *fmtWriter) Fprintln(a ...interface{}) (n int, err error) {
	if fw.err != nil {
		// early return if a previous error has been encountered.
		return 0, nil
	}
	n, err = fmt.Fprintln(fw.w, a...)
	fw.size += int64(n)
	fw.err = err
	return n, err
}

// namedVar is a named variable.
type namedVar interface {
	value.Named
	// ID returns the ID of the local identifier.
	ID() int64
	// SetID sets the ID of the local identifier.
	SetID(id int64)
	// IsUnnamed reports whether the local identifier is unnamed.
	IsUnnamed() bool
}
