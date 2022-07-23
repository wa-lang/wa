package llir

import (
	"fmt"
	"strings"

	enum "github.com/wa-lang/wa/internal/3rdparty/llir/llenum"
	types "github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
	value "github.com/wa-lang/wa/internal/3rdparty/llir/llvalue"
)

// --- [ Other instructions ] --------------------------------------------------

// ~~~ [ icmp ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// InstICmp is an LLVM IR icmp instruction.
type InstICmp struct {
	// Name of local variable associated with the result.
	LocalIdent
	// Integer comparison predicate.
	Pred enum.IPred
	// Integer scalar or vector operands.
	X, Y value.Value // integer scalar, pointer, integer vector or pointer vector.

	// extra.

	// Type of result produced by the instruction.
	Typ types.Type // boolean or boolean vector
	// (optional) Metadata.
	Metadata
}

// NewICmp returns a new icmp instruction based on the given integer comparison
// predicate and integer scalar or vector operands.
func NewICmp(pred enum.IPred, x, y value.Value) *InstICmp {
	inst := &InstICmp{Pred: pred, X: x, Y: y}
	// Compute type.
	inst.Type()
	return inst
}

// String returns the LLVM syntax representation of the instruction as a
// type-value pair.
func (inst *InstICmp) String() string {
	return fmt.Sprintf("%s %s", inst.Type(), inst.Ident())
}

// Type returns the type of the instruction. The result type is either boolean
// type or vector of booleans type.
func (inst *InstICmp) Type() types.Type {
	// Cache type if not present.
	if inst.Typ == nil {
		switch xType := inst.X.Type().(type) {
		case *types.IntType, *types.PointerType:
			inst.Typ = types.I1
		case *types.VectorType:
			inst.Typ = types.NewVector(xType.Len, types.I1)
		default:
			panic(fmt.Errorf("invalid icmp operand type; expected *types.IntType, *types.PointerType or *types.VectorType, got %T", xType))
		}
	}
	return inst.Typ
}

// LLString returns the LLVM syntax representation of the instruction.
//
// 'icmp' Pred=IPred X=TypeValue ',' Y=Value Metadata=(',' MetadataAttachment)+?
func (inst *InstICmp) LLString() string {
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "%s = ", inst.Ident())
	fmt.Fprintf(buf, "icmp %s %s, %s", inst.Pred, inst.X, inst.Y.Ident())
	for _, md := range inst.Metadata {
		fmt.Fprintf(buf, ", %s", md)
	}
	return buf.String()
}

// ~~~ [ fcmp ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// InstFCmp is an LLVM IR fcmp instruction.
type InstFCmp struct {
	// Name of local variable associated with the result.
	LocalIdent
	// Floating-point comparison predicate.
	Pred enum.FPred
	// Floating-point scalar or vector operands.
	X, Y value.Value // floating-point scalar or floating-point vector

	// extra.

	// Type of result produced by the instruction.
	Typ types.Type // boolean or boolean vector
	// (optional) Fast math flags.
	FastMathFlags []enum.FastMathFlag
	// (optional) Metadata.
	Metadata
}

// NewFCmp returns a new fcmp instruction based on the given floating-point
// comparison predicate and floating-point scalar or vector operands.
func NewFCmp(pred enum.FPred, x, y value.Value) *InstFCmp {
	inst := &InstFCmp{Pred: pred, X: x, Y: y}
	// Compute type.
	inst.Type()
	return inst
}

// String returns the LLVM syntax representation of the instruction as a
// type-value pair.
func (inst *InstFCmp) String() string {
	return fmt.Sprintf("%s %s", inst.Type(), inst.Ident())
}

// Type returns the type of the instruction. The result type is either boolean
// type or vector of booleans type.
func (inst *InstFCmp) Type() types.Type {
	// Cache type if not present.
	if inst.Typ == nil {
		switch xType := inst.X.Type().(type) {
		case *types.FloatType:
			inst.Typ = types.I1
		case *types.VectorType:
			inst.Typ = types.NewVector(xType.Len, types.I1)
		default:
			panic(fmt.Errorf("invalid fcmp operand type; expected *types.FloatType or *types.VectorType, got %T", xType))
		}
	}
	return inst.Typ
}

// LLString returns the LLVM syntax representation of the instruction.
//
// 'fcmp' FastMathFlags=FastMathFlag* Pred=FPred X=TypeValue ',' Y=Value Metadata=(',' MetadataAttachment)+?
func (inst *InstFCmp) LLString() string {
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "%s = ", inst.Ident())
	buf.WriteString("fcmp")
	for _, flag := range inst.FastMathFlags {
		fmt.Fprintf(buf, " %s", flag)
	}
	fmt.Fprintf(buf, " %s %s, %s", inst.Pred, inst.X, inst.Y.Ident())
	for _, md := range inst.Metadata {
		fmt.Fprintf(buf, ", %s", md)
	}
	return buf.String()
}

// ~~~ [ phi ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// InstPhi is an LLVM IR phi instruction.
type InstPhi struct {
	// Name of local variable associated with the result.
	LocalIdent
	// Incoming values.
	Incs []*Incoming

	// extra.

	// Type of result produced by the instruction.
	Typ types.Type // type of incoming value
	// (optional) Fast math flags.
	FastMathFlags []enum.FastMathFlag
	// (optional) Metadata.
	Metadata
}

// NewPhi returns a new phi instruction based on the given incoming values.
func NewPhi(incs ...*Incoming) *InstPhi {
	inst := &InstPhi{Incs: incs}
	// Compute type.
	inst.Type()
	return inst
}

// String returns the LLVM syntax representation of the instruction as a
// type-value pair.
func (inst *InstPhi) String() string {
	return fmt.Sprintf("%s %s", inst.Type(), inst.Ident())
}

// Type returns the type of the instruction. The result type is the type of the
// incoming value.
func (inst *InstPhi) Type() types.Type {
	// Cache type if not present.
	if inst.Typ == nil {
		inst.Typ = inst.Incs[0].X.Type()
	}
	return inst.Typ
}

// LLString returns the LLVM syntax representation of the instruction.
//
// 'phi' Typ=Type Incs=(Inc separator ',')+ Metadata=(',' MetadataAttachment)+?
func (inst *InstPhi) LLString() string {
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "%s = ", inst.Ident())
	buf.WriteString("phi ")
	for _, flag := range inst.FastMathFlags {
		buf.WriteString(flag.String())
		buf.WriteString(" ")
	}
	buf.WriteString(inst.Typ.String())
	buf.WriteString(" ")
	for i, inc := range inst.Incs {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(inc.String())
	}
	for _, md := range inst.Metadata {
		fmt.Fprintf(buf, ", %s", md)
	}
	return buf.String()
}

// ___ [ Incoming value ] ______________________________________________________

// Incoming is an incoming value of a phi instruction.
type Incoming struct {
	// Incoming value.
	X value.Value
	// Predecessor basic block of the incoming value.
	Pred value.Value // *ir.Block
}

// NewIncoming returns a new incoming value based on the given value and
// predecessor basic block.
func NewIncoming(x value.Value, pred *Block) *Incoming {
	return &Incoming{X: x, Pred: pred}
}

// String returns the string representation of the incoming value.
func (inc *Incoming) String() string {
	// '[' X=Value ',' Pred=LocalIdent ']'
	return fmt.Sprintf("[ %s, %s ]", inc.X.Ident(), inc.Pred.Ident())
}

// ~~~ [ select ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// InstSelect is an LLVM IR select instruction.
type InstSelect struct {
	// Name of local variable associated with the result.
	LocalIdent
	// Selection condition.
	Cond value.Value // boolean or boolean vector
	// True condition value.
	ValueTrue value.Value
	// False condition value.
	ValueFalse value.Value

	// extra.

	// Type of result produced by the instruction.
	Typ types.Type
	// (optional) Fast math flags.
	FastMathFlags []enum.FastMathFlag
	// (optional) Metadata.
	Metadata
}

// NewSelect returns a new select instruction based on the given selection
// condition and true and false condition values.
func NewSelect(cond, valueTrue, valueFalse value.Value) *InstSelect {
	inst := &InstSelect{Cond: cond, ValueTrue: valueTrue, ValueFalse: valueFalse}
	// Compute type.
	inst.Type()
	return inst
}

// String returns the LLVM syntax representation of the instruction as a
// type-value pair.
func (inst *InstSelect) String() string {
	return fmt.Sprintf("%s %s", inst.Type(), inst.Ident())
}

// Type returns the type of the instruction.
func (inst *InstSelect) Type() types.Type {
	// Cache type if not present.
	if inst.Typ == nil {
		inst.Typ = inst.ValueTrue.Type()
	}
	return inst.Typ
}

// LLString returns the LLVM syntax representation of the instruction.
//
// 'select' FastMathFlags=FastMathFlag* Cond=TypeValue ',' ValueTrue=TypeValue ',' ValueFalse=TypeValue Metadata=(',' MetadataAttachment)+?
func (inst *InstSelect) LLString() string {
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "%s = ", inst.Ident())
	buf.WriteString("select")
	for _, flag := range inst.FastMathFlags {
		fmt.Fprintf(buf, " %s", flag)
	}
	fmt.Fprintf(buf, " %s, %s, %s", inst.Cond, inst.ValueTrue, inst.ValueFalse)
	for _, md := range inst.Metadata {
		fmt.Fprintf(buf, ", %s", md)
	}
	return buf.String()
}

// ~~~ [ freeze ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// InstFreeze is an LLVM IR freeze instruction.
type InstFreeze struct {
	// Name of local variable associated with the result.
	LocalIdent
	// Operand.
	X value.Value
	// extra.

	// Type of result produced by the instruction.
	Typ types.Type
	// (optional) Metadata.
	Metadata
}

// NewInstFreeze returns a new freeze instruction based on the given
// operand.
func NewInstFreeze(x value.Value) *InstFreeze {
	inst := &InstFreeze{X: x}
	// Compute type.
	inst.Type()
	return inst
}

// String returns the LLVM syntax representation of the instruction as a
// type-value pair.
func (inst *InstFreeze) String() string {
	return fmt.Sprintf("%s %s", inst.Type(), inst.Ident())
}

// Type returns the type of the instruction.
func (inst *InstFreeze) Type() types.Type {
	// Cache type if not present.
	if inst.Typ == nil {
		inst.Typ = inst.X.Type()
	}
	return inst.Typ
}

// LLString returns the LLVM syntax representation of the instruction.
//
// 'freeze' Type Value
func (inst *InstFreeze) LLString() string {
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "%s = ", inst.Ident())
	fmt.Fprintf(buf, "freeze %s", inst.X)
	for _, md := range inst.Metadata {
		fmt.Fprintf(buf, ", %s", md)
	}
	return buf.String()
}

// ~~~ [ call ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// InstCall is an LLVM IR call instruction.
type InstCall struct {
	// Name of local variable associated with the result.
	LocalIdent
	// Callee.
	// TODO: specify the set of underlying types of Callee.
	Callee value.Value
	// Function arguments.
	//
	// Arg has one of the following underlying types:
	//    value.Value
	//    *ir.Arg
	//    TODO: add metadata value?
	Args []value.Value

	// extra.

	// Type of result produced by the instruction.
	Typ types.Type
	// (optional) Tail; zero if not present.
	Tail enum.Tail
	// (optional) Fast math flags.
	FastMathFlags []enum.FastMathFlag
	// (optional) Calling convention; zero if not present.
	CallingConv enum.CallingConv
	// (optional) Return attributes.
	ReturnAttrs []ReturnAttribute
	// (optional) Address space; zero if not present.
	AddrSpace types.AddrSpace
	// (optional) Function attributes.
	FuncAttrs []FuncAttribute
	// (optional) Operand bundles.
	OperandBundles []*OperandBundle
	// (optional) Metadata.
	Metadata
}

// NewCall returns a new call instruction based on the given callee and function
// arguments.
//
// TODO: specify the set of underlying types of callee.
func NewCall(callee value.Value, args ...value.Value) *InstCall {
	inst := &InstCall{Callee: callee, Args: args}
	// Compute type.
	inst.Type()
	return inst
}

// String returns the LLVM syntax representation of the instruction as a
// type-value pair.
func (inst *InstCall) String() string {
	return fmt.Sprintf("%s %s", inst.Type(), inst.Ident())
}

// Type returns the type of the instruction.
func (inst *InstCall) Type() types.Type {
	// Cache type if not present.
	if inst.Typ == nil {
		sig := inst.Sig()
		inst.Typ = sig.RetType
	}
	return inst.Typ
}

// LLString returns the LLVM syntax representation of the instruction.
//
// Tailopt 'call' FastMathFlags=FastMathFlag* CallingConvopt ReturnAttrs=ReturnAttribute* AddrSpaceopt Typ=Type Callee=Value '(' Args ')' FuncAttrs=FuncAttribute* OperandBundles=('[' (OperandBundle separator ',')+ ']')? Metadata=(',' MetadataAttachment)+?
func (inst *InstCall) LLString() string {
	buf := &strings.Builder{}
	if !inst.Type().Equal(types.Void) {
		fmt.Fprintf(buf, "%s = ", inst.Ident())
	}
	if inst.Tail != enum.TailNone {
		fmt.Fprintf(buf, "%s ", inst.Tail)
	}
	buf.WriteString("call")
	for _, flag := range inst.FastMathFlags {
		fmt.Fprintf(buf, " %s", flag)
	}
	if inst.CallingConv != enum.CallingConvNone {
		fmt.Fprintf(buf, " %s", callingConvString(inst.CallingConv))
	}
	for _, attr := range inst.ReturnAttrs {
		fmt.Fprintf(buf, " %s", attr)
	}
	// (optional) Address space.
	if inst.AddrSpace != 0 {
		fmt.Fprintf(buf, " %s", inst.AddrSpace)
	}
	// Use function signature instead of return type for variadic functions.
	calleeType := inst.Type()
	if sig := inst.Sig(); sig.Variadic {
		calleeType = sig
	}
	fmt.Fprintf(buf, " %s %s(", calleeType, inst.Callee.Ident())
	for i, arg := range inst.Args {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(arg.String())
	}
	buf.WriteString(")")
	for _, attr := range inst.FuncAttrs {
		fmt.Fprintf(buf, " %s", attr)
	}
	if len(inst.OperandBundles) > 0 {
		buf.WriteString(" [ ")
		for i, operandBundle := range inst.OperandBundles {
			if i != 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(operandBundle.String())
		}
		buf.WriteString(" ]")
	}
	for _, md := range inst.Metadata {
		fmt.Fprintf(buf, ", %s", md)
	}
	return buf.String()
}

// Sig returns the function signature of the callee.
func (inst *InstCall) Sig() *types.FuncType {
	t, ok := inst.Callee.Type().(*types.PointerType)
	if !ok {
		panic(fmt.Errorf("invalid callee type; expected *types.PointerType, got %T", inst.Callee.Type()))
	}
	sig, ok := t.ElemType.(*types.FuncType)
	if !ok {
		panic(fmt.Errorf("invalid callee type; expected *types.FuncType, got %T", t.ElemType))
	}
	return sig
}

// ~~~ [ va_arg ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// InstVAArg is an LLVM IR va_arg instruction.
type InstVAArg struct {
	// Name of local variable associated with the result.
	LocalIdent
	// Variable argument list.
	ArgList value.Value
	// Argument type.
	ArgType types.Type

	// extra.

	// (optional) Metadata.
	Metadata
}

// NewVAArg returns a new va_arg instruction based on the given variable
// argument list and argument type.
func NewVAArg(argList value.Value, argType types.Type) *InstVAArg {
	return &InstVAArg{ArgList: argList, ArgType: argType}
}

// String returns the LLVM syntax representation of the instruction as a
// type-value pair.
func (inst *InstVAArg) String() string {
	return fmt.Sprintf("%s %s", inst.Type(), inst.Ident())
}

// Type returns the type of the instruction.
func (inst *InstVAArg) Type() types.Type {
	return inst.ArgType
}

// LLString returns the LLVM syntax representation of the instruction.
//
// 'va_arg' ArgList=TypeValue ',' ArgType=Type Metadata=(',' MetadataAttachment)+?
func (inst *InstVAArg) LLString() string {
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "%s = ", inst.Ident())
	fmt.Fprintf(buf, "va_arg %s, %s", inst.ArgList, inst.ArgType)
	for _, md := range inst.Metadata {
		fmt.Fprintf(buf, ", %s", md)
	}
	return buf.String()
}

// ~~~ [ landingpad ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// InstLandingPad is an LLVM IR landingpad instruction.
type InstLandingPad struct {
	// Name of local variable associated with the result.
	LocalIdent
	// Result type.
	ResultType types.Type
	// (optional) Cleanup landing pad.
	Cleanup bool
	// Filter and catch clauses; zero or more if Cleanup is true, otherwise one
	// or more.
	Clauses []*Clause

	// extra.

	// (optional) Metadata.
	Metadata
}

// NewLandingPad returns a new landingpad instruction based on the given result
// type and filter/catch clauses.
func NewLandingPad(resultType types.Type, clauses ...*Clause) *InstLandingPad {
	return &InstLandingPad{ResultType: resultType, Clauses: clauses}
}

// String returns the LLVM syntax representation of the instruction as a
// type-value pair.
func (inst *InstLandingPad) String() string {
	return fmt.Sprintf("%s %s", inst.Type(), inst.Ident())
}

// Type returns the type of the instruction.
func (inst *InstLandingPad) Type() types.Type {
	return inst.ResultType
}

// LLString returns the LLVM syntax representation of the instruction.
//
// 'landingpad' ResultType=Type Cleanupopt Clauses=Clause* Metadata=(',' MetadataAttachment)+?
func (inst *InstLandingPad) LLString() string {
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "%s = ", inst.Ident())
	fmt.Fprintf(buf, "landingpad %s", inst.ResultType)
	if inst.Cleanup {
		buf.WriteString("\n\t\tcleanup")
	}
	for _, clause := range inst.Clauses {
		fmt.Fprintf(buf, "\n\t\t%s", clause)
	}
	for _, md := range inst.Metadata {
		fmt.Fprintf(buf, ", %s", md)
	}
	return buf.String()
}

// ___ [ Landingpad clause ] ___________________________________________________

// Clause is a landingpad catch or filter clause.
type Clause struct {
	// Clause type (catch or filter).
	Type enum.ClauseType
	// Operand.
	X value.Value
}

// NewClause returns a new landingpad clause based on the given clause type and
// operand.
func NewClause(clauseType enum.ClauseType, x value.Value) *Clause {
	return &Clause{Type: clauseType, X: x}
}

// String returns the string representation of the landingpad clause.
func (clause *Clause) String() string {
	return fmt.Sprintf("%s %s", clause.Type, clause.X)
}

// ~~~ [ catchpad ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// InstCatchPad is an LLVM IR catchpad instruction.
type InstCatchPad struct {
	// Name of local variable associated with the result.
	LocalIdent
	// Parent catchswitch terminator.
	CatchSwitch value.Value // *ir.TermCatchSwitch
	// Exception arguments.
	//
	// Arg has one of the following underlying types:
	//    value.Value
	//    TODO: add metadata value?
	Args []value.Value

	// extra.

	// (optional) Metadata.
	Metadata
}

// NewCatchPad returns a new catchpad instruction based on the given parent
// catchswitch terminator and exception arguments.
func NewCatchPad(catchSwitch *TermCatchSwitch, args ...value.Value) *InstCatchPad {
	return &InstCatchPad{CatchSwitch: catchSwitch, Args: args}
}

// String returns the LLVM syntax representation of the instruction as a
// type-value pair.
func (inst *InstCatchPad) String() string {
	return fmt.Sprintf("%s %s", inst.Type(), inst.Ident())
}

// Type returns the type of the instruction.
func (inst *InstCatchPad) Type() types.Type {
	return types.Token
}

// LLString returns the LLVM syntax representation of the instruction.
//
// 'catchpad' 'within' CatchSwitch=LocalIdent '[' Args=(ExceptionArg separator ',')* ']' Metadata=(',' MetadataAttachment)+?
func (inst *InstCatchPad) LLString() string {
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "%s = ", inst.Ident())
	fmt.Fprintf(buf, "catchpad within %s [", inst.CatchSwitch.Ident())
	for i, arg := range inst.Args {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(arg.String())
	}
	buf.WriteString("]")
	for _, md := range inst.Metadata {
		fmt.Fprintf(buf, ", %s", md)
	}
	return buf.String()
}

// ~~~ [ cleanuppad ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// InstCleanupPad is an LLVM IR cleanuppad instruction.
type InstCleanupPad struct {
	// Name of local variable associated with the result.
	LocalIdent
	// Parent exception pad.
	ParentPad value.Value // ir.ExceptionPad
	// Exception arguments.
	//
	// Arg has one of the following underlying types:
	//    value.Value
	//    TODO: add metadata value?
	Args []value.Value

	// extra.

	// (optional) Metadata.
	Metadata
}

// NewCleanupPad returns a new cleanuppad instruction based on the given
// parent exception pad and exception arguments.
func NewCleanupPad(parentPad ExceptionPad, args ...value.Value) *InstCleanupPad {
	return &InstCleanupPad{ParentPad: parentPad, Args: args}
}

// String returns the LLVM syntax representation of the instruction as a
// type-value pair.
func (inst *InstCleanupPad) String() string {
	return fmt.Sprintf("%s %s", inst.Type(), inst.Ident())
}

// Type returns the type of the instruction.
func (inst *InstCleanupPad) Type() types.Type {
	return types.Token
}

// LLString returns the LLVM syntax representation of the instruction.
//
// 'cleanuppad' 'within' ParentPad=ExceptionPad '[' Args=(ExceptionArg separator ',')* ']' Metadata=(',' MetadataAttachment)+?
func (inst *InstCleanupPad) LLString() string {
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "%s = ", inst.Ident())
	fmt.Fprintf(buf, "cleanuppad within %s [", inst.ParentPad.Ident())
	for i, arg := range inst.Args {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(arg.String())
	}
	buf.WriteString("]")
	for _, md := range inst.Metadata {
		fmt.Fprintf(buf, ", %s", md)
	}
	return buf.String()
}
