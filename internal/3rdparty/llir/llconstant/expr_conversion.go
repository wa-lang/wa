package llconstant

import (
	"fmt"

	types "github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
)

// --- [ Conversion expressions ] ----------------------------------------------

// ~~~ [ trunc ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ExprTrunc is an LLVM IR trunc expression.
type ExprTrunc struct {
	// Value before conversion.
	From Constant
	// Type after conversion.
	To types.Type
}

// NewTrunc returns a new trunc expression based on the given source value and
// target type.
func NewTrunc(from Constant, to types.Type) *ExprTrunc {
	e := &ExprTrunc{From: from, To: to}
	// Compute type.
	e.Type()
	return e
}

// String returns the LLVM syntax representation of the constant expression as a
// type-value pair.
func (e *ExprTrunc) String() string {
	return fmt.Sprintf("%s %s", e.Type(), e.Ident())
}

// Type returns the type of the constant expression.
func (e *ExprTrunc) Type() types.Type {
	return e.To
}

// Ident returns the identifier associated with the constant expression.
func (e *ExprTrunc) Ident() string {
	// 'trunc' '(' From=TypeConst 'to' To=Type ')'
	return fmt.Sprintf("trunc (%s to %s)", e.From, e.To)
}

// ~~~ [ zext ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ExprZExt is an LLVM IR zext expression.
type ExprZExt struct {
	// Value before conversion.
	From Constant
	// Type after conversion.
	To types.Type
}

// NewZExt returns a new zext expression based on the given source value and
// target type.
func NewZExt(from Constant, to types.Type) *ExprZExt {
	e := &ExprZExt{From: from, To: to}
	// Compute type.
	e.Type()
	return e
}

// String returns the LLVM syntax representation of the constant expression as a
// type-value pair.
func (e *ExprZExt) String() string {
	return fmt.Sprintf("%s %s", e.Type(), e.Ident())
}

// Type returns the type of the constant expression.
func (e *ExprZExt) Type() types.Type {
	return e.To
}

// Ident returns the identifier associated with the constant expression.
func (e *ExprZExt) Ident() string {
	// 'zext' '(' From=TypeConst 'to' To=Type ')'
	return fmt.Sprintf("zext (%s to %s)", e.From, e.To)
}

// ~~~ [ sext ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ExprSExt is an LLVM IR sext expression.
type ExprSExt struct {
	// Value before conversion.
	From Constant
	// Type after conversion.
	To types.Type
}

// NewSExt returns a new sext expression based on the given source value and
// target type.
func NewSExt(from Constant, to types.Type) *ExprSExt {
	e := &ExprSExt{From: from, To: to}
	// Compute type.
	e.Type()
	return e
}

// String returns the LLVM syntax representation of the constant expression as a
// type-value pair.
func (e *ExprSExt) String() string {
	return fmt.Sprintf("%s %s", e.Type(), e.Ident())
}

// Type returns the type of the constant expression.
func (e *ExprSExt) Type() types.Type {
	return e.To
}

// Ident returns the identifier associated with the constant expression.
func (e *ExprSExt) Ident() string {
	// 'sext' '(' From=TypeConst 'to' To=Type ')'
	return fmt.Sprintf("sext (%s to %s)", e.From, e.To)
}

// ~~~ [ fptrunc ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ExprFPTrunc is an LLVM IR fptrunc expression.
type ExprFPTrunc struct {
	// Value before conversion.
	From Constant
	// Type after conversion.
	To types.Type
}

// NewFPTrunc returns a new fptrunc expression based on the given source value
// and target type.
func NewFPTrunc(from Constant, to types.Type) *ExprFPTrunc {
	e := &ExprFPTrunc{From: from, To: to}
	// Compute type.
	e.Type()
	return e
}

// String returns the LLVM syntax representation of the constant expression as a
// type-value pair.
func (e *ExprFPTrunc) String() string {
	return fmt.Sprintf("%s %s", e.Type(), e.Ident())
}

// Type returns the type of the constant expression.
func (e *ExprFPTrunc) Type() types.Type {
	return e.To
}

// Ident returns the identifier associated with the constant expression.
func (e *ExprFPTrunc) Ident() string {
	// 'fptrunc' '(' From=TypeConst 'to' To=Type ')'
	return fmt.Sprintf("fptrunc (%s to %s)", e.From, e.To)
}

// ~~~ [ fpext ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ExprFPExt is an LLVM IR fpext expression.
type ExprFPExt struct {
	// Value before conversion.
	From Constant
	// Type after conversion.
	To types.Type
}

// NewFPExt returns a new fpext expression based on the given source value and
// target type.
func NewFPExt(from Constant, to types.Type) *ExprFPExt {
	e := &ExprFPExt{From: from, To: to}
	// Compute type.
	e.Type()
	return e
}

// String returns the LLVM syntax representation of the constant expression as a
// type-value pair.
func (e *ExprFPExt) String() string {
	return fmt.Sprintf("%s %s", e.Type(), e.Ident())
}

// Type returns the type of the constant expression.
func (e *ExprFPExt) Type() types.Type {
	return e.To
}

// Ident returns the identifier associated with the constant expression.
func (e *ExprFPExt) Ident() string {
	// 'fpext' '(' From=TypeConst 'to' To=Type ')'
	return fmt.Sprintf("fpext (%s to %s)", e.From, e.To)
}

// ~~~ [ fptoui ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ExprFPToUI is an LLVM IR fptoui expression.
type ExprFPToUI struct {
	// Value before conversion.
	From Constant
	// Type after conversion.
	To types.Type
}

// NewFPToUI returns a new fptoui expression based on the given source value and
// target type.
func NewFPToUI(from Constant, to types.Type) *ExprFPToUI {
	e := &ExprFPToUI{From: from, To: to}
	// Compute type.
	e.Type()
	return e
}

// String returns the LLVM syntax representation of the constant expression as a
// type-value pair.
func (e *ExprFPToUI) String() string {
	return fmt.Sprintf("%s %s", e.Type(), e.Ident())
}

// Type returns the type of the constant expression.
func (e *ExprFPToUI) Type() types.Type {
	return e.To
}

// Ident returns the identifier associated with the constant expression.
func (e *ExprFPToUI) Ident() string {
	// 'fptoui' '(' From=TypeConst 'to' To=Type ')'
	return fmt.Sprintf("fptoui (%s to %s)", e.From, e.To)
}

// ~~~ [ fptosi ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ExprFPToSI is an LLVM IR fptosi expression.
type ExprFPToSI struct {
	// Value before conversion.
	From Constant
	// Type after conversion.
	To types.Type
}

// NewFPToSI returns a new fptosi expression based on the given source value and
// target type.
func NewFPToSI(from Constant, to types.Type) *ExprFPToSI {
	e := &ExprFPToSI{From: from, To: to}
	// Compute type.
	e.Type()
	return e
}

// String returns the LLVM syntax representation of the constant expression as a
// type-value pair.
func (e *ExprFPToSI) String() string {
	return fmt.Sprintf("%s %s", e.Type(), e.Ident())
}

// Type returns the type of the constant expression.
func (e *ExprFPToSI) Type() types.Type {
	return e.To
}

// Ident returns the identifier associated with the constant expression.
func (e *ExprFPToSI) Ident() string {
	// 'fptosi' '(' From=TypeConst 'to' To=Type ')'
	return fmt.Sprintf("fptosi (%s to %s)", e.From, e.To)
}

// ~~~ [ uitofp ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ExprUIToFP is an LLVM IR uitofp expression.
type ExprUIToFP struct {
	// Value before conversion.
	From Constant
	// Type after conversion.
	To types.Type
}

// NewUIToFP returns a new uitofp expression based on the given source value and
// target type.
func NewUIToFP(from Constant, to types.Type) *ExprUIToFP {
	e := &ExprUIToFP{From: from, To: to}
	// Compute type.
	e.Type()
	return e
}

// String returns the LLVM syntax representation of the constant expression as a
// type-value pair.
func (e *ExprUIToFP) String() string {
	return fmt.Sprintf("%s %s", e.Type(), e.Ident())
}

// Type returns the type of the constant expression.
func (e *ExprUIToFP) Type() types.Type {
	return e.To
}

// Ident returns the identifier associated with the constant expression.
func (e *ExprUIToFP) Ident() string {
	// 'uitofp' '(' From=TypeConst 'to' To=Type ')'
	return fmt.Sprintf("uitofp (%s to %s)", e.From, e.To)
}

// ~~~ [ sitofp ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ExprSIToFP is an LLVM IR sitofp expression.
type ExprSIToFP struct {
	// Value before conversion.
	From Constant
	// Type after conversion.
	To types.Type
}

// NewSIToFP returns a new sitofp expression based on the given source value and
// target type.
func NewSIToFP(from Constant, to types.Type) *ExprSIToFP {
	e := &ExprSIToFP{From: from, To: to}
	// Compute type.
	e.Type()
	return e
}

// String returns the LLVM syntax representation of the constant expression as a
// type-value pair.
func (e *ExprSIToFP) String() string {
	return fmt.Sprintf("%s %s", e.Type(), e.Ident())
}

// Type returns the type of the constant expression.
func (e *ExprSIToFP) Type() types.Type {
	return e.To
}

// Ident returns the identifier associated with the constant expression.
func (e *ExprSIToFP) Ident() string {
	// 'sitofp' '(' From=TypeConst 'to' To=Type ')'
	return fmt.Sprintf("sitofp (%s to %s)", e.From, e.To)
}

// ~~~ [ ptrtoint ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ExprPtrToInt is an LLVM IR ptrtoint expression.
type ExprPtrToInt struct {
	// Value before conversion.
	From Constant
	// Type after conversion.
	To types.Type
}

// NewPtrToInt returns a new ptrtoint expression based on the given source value
// and target type.
func NewPtrToInt(from Constant, to types.Type) *ExprPtrToInt {
	e := &ExprPtrToInt{From: from, To: to}
	// Compute type.
	e.Type()
	return e
}

// String returns the LLVM syntax representation of the constant expression as a
// type-value pair.
func (e *ExprPtrToInt) String() string {
	return fmt.Sprintf("%s %s", e.Type(), e.Ident())
}

// Type returns the type of the constant expression.
func (e *ExprPtrToInt) Type() types.Type {
	return e.To
}

// Ident returns the identifier associated with the constant expression.
func (e *ExprPtrToInt) Ident() string {
	// 'ptrtoint' '(' From=TypeConst 'to' To=Type ')'
	return fmt.Sprintf("ptrtoint (%s to %s)", e.From, e.To)
}

// ~~~ [ inttoptr ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ExprIntToPtr is an LLVM IR inttoptr expression.
type ExprIntToPtr struct {
	// Value before conversion.
	From Constant
	// Type after conversion.
	To types.Type
}

// NewIntToPtr returns a new inttoptr expression based on the given source value
// and target type.
func NewIntToPtr(from Constant, to types.Type) *ExprIntToPtr {
	e := &ExprIntToPtr{From: from, To: to}
	// Compute type.
	e.Type()
	return e
}

// String returns the LLVM syntax representation of the constant expression as a
// type-value pair.
func (e *ExprIntToPtr) String() string {
	return fmt.Sprintf("%s %s", e.Type(), e.Ident())
}

// Type returns the type of the constant expression.
func (e *ExprIntToPtr) Type() types.Type {
	return e.To
}

// Ident returns the identifier associated with the constant expression.
func (e *ExprIntToPtr) Ident() string {
	// 'inttoptr' '(' From=TypeConst 'to' To=Type ')'
	return fmt.Sprintf("inttoptr (%s to %s)", e.From, e.To)
}

// ~~~ [ bitcast ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ExprBitCast is an LLVM IR bitcast expression.
type ExprBitCast struct {
	// Value before conversion.
	From Constant
	// Type after conversion.
	To types.Type
}

// NewBitCast returns a new bitcast expression based on the given source value
// and target type.
func NewBitCast(from Constant, to types.Type) *ExprBitCast {
	e := &ExprBitCast{From: from, To: to}
	// Compute type.
	e.Type()
	return e
}

// String returns the LLVM syntax representation of the constant expression as a
// type-value pair.
func (e *ExprBitCast) String() string {
	return fmt.Sprintf("%s %s", e.Type(), e.Ident())
}

// Type returns the type of the constant expression.
func (e *ExprBitCast) Type() types.Type {
	return e.To
}

// Ident returns the identifier associated with the constant expression.
func (e *ExprBitCast) Ident() string {
	// 'bitcast' '(' From=TypeConst 'to' To=Type ')'
	return fmt.Sprintf("bitcast (%s to %s)", e.From, e.To)
}

// ~~~ [ addrspacecast ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ExprAddrSpaceCast is an LLVM IR addrspacecast expression.
type ExprAddrSpaceCast struct {
	// Value before conversion.
	From Constant
	// Type after conversion.
	To types.Type
}

// NewAddrSpaceCast returns a new addrspacecast expression based on the given
// source value and target type.
func NewAddrSpaceCast(from Constant, to types.Type) *ExprAddrSpaceCast {
	e := &ExprAddrSpaceCast{From: from, To: to}
	// Compute type.
	e.Type()
	return e
}

// String returns the LLVM syntax representation of the constant expression as a
// type-value pair.
func (e *ExprAddrSpaceCast) String() string {
	return fmt.Sprintf("%s %s", e.Type(), e.Ident())
}

// Type returns the type of the constant expression.
func (e *ExprAddrSpaceCast) Type() types.Type {
	return e.To
}

// Ident returns the identifier associated with the constant expression.
func (e *ExprAddrSpaceCast) Ident() string {
	// 'addrspacecast' '(' From=TypeConst 'to' To=Type ')'
	return fmt.Sprintf("addrspacecast (%s to %s)", e.From, e.To)
}
