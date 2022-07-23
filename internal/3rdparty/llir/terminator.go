package llir

import (
	"fmt"
	"strings"

	constant "github.com/wa-lang/wa/internal/3rdparty/llir/llconstant"
	enum "github.com/wa-lang/wa/internal/3rdparty/llir/llenum"
	types "github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
	value "github.com/wa-lang/wa/internal/3rdparty/llir/llvalue"
)

// === [ Terminators ] =========================================================

// Terminator is an LLVM IR terminator instruction (a control flow instruction).
//
// A Terminator has one of the following underlying types.
//
// Terminators
//
// https://llvm.org/docs/LangRef.html#terminator-instructions
//
//    *ir.TermRet           // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#TermRet
//    *ir.TermBr            // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#TermBr
//    *ir.TermCondBr        // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#TermCondBr
//    *ir.TermSwitch        // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#TermSwitch
//    *ir.TermIndirectBr    // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#TermIndirectBr
//    *ir.TermInvoke        // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#TermInvoke
//    *ir.TermCallBr        // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#TermCallBr
//    *ir.TermResume        // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#TermResume
//    *ir.TermCatchSwitch   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#TermCatchSwitch
//    *ir.TermCatchRet      // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#TermCatchRet
//    *ir.TermCleanupRet    // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#TermCleanupRet
//    *ir.TermUnreachable   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#TermUnreachable
type Terminator interface {
	LLStringer
	// Succs returns the successor basic blocks of the terminator.
	Succs() []*Block
}

// --- [ ret ] -----------------------------------------------------------------

// TermRet is an LLVM IR ret terminator.
type TermRet struct {
	// Return value; or nil if void return.
	X value.Value

	// extra.

	// (optional) Metadata.
	Metadata
}

// NewRet returns a new ret terminator based on the given return value. A nil
// return value indicates a void return.
func NewRet(x value.Value) *TermRet {
	return &TermRet{X: x}
}

// Succs returns the successor basic blocks of the terminator.
func (*TermRet) Succs() []*Block {
	// no successors.
	return nil
}

// LLString returns the LLVM syntax representation of the terminator.
//
// Void return instruction.
//
//    'ret' XTyp=VoidType Metadata=(',' MetadataAttachment)+?
//
// Value return instruction.
//
//    'ret' XTyp=ConcreteType X=Value Metadata=(',' MetadataAttachment)+?
func (term *TermRet) LLString() string {
	buf := &strings.Builder{}
	if term.X == nil {
		buf.WriteString("ret void")
	} else {
		fmt.Fprintf(buf, "ret %s", term.X)
	}
	for _, md := range term.Metadata {
		fmt.Fprintf(buf, ", %s", md)
	}
	return buf.String()
}

// --- [ br ] ------------------------------------------------------------------

// TermBr is an unconditional LLVM IR br terminator.
type TermBr struct {
	// Target branch.
	Target value.Value // *ir.Block

	// extra.

	// Successor basic blocks of the terminator.
	Successors []*Block
	// (optional) Metadata.
	Metadata
}

// NewBr returns a new unconditional br terminator based on the given target
// basic block.
func NewBr(target *Block) *TermBr {
	return &TermBr{Target: target}
}

// Succs returns the successor basic blocks of the terminator.
func (term *TermBr) Succs() []*Block {
	// Cache successors if not present.
	if term.Successors == nil {
		term.Successors = []*Block{term.Target.(*Block)}
	}
	return term.Successors
}

// LLString returns the LLVM syntax representation of the terminator.
//
// 'br' Target=Label Metadata=(',' MetadataAttachment)+?
func (term *TermBr) LLString() string {
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "br %s", term.Target)
	for _, md := range term.Metadata {
		fmt.Fprintf(buf, ", %s", md)
	}
	return buf.String()
}

// --- [ conditional br ] ------------------------------------------------------

// TermCondBr is a conditional LLVM IR br terminator.
type TermCondBr struct {
	// Branching condition.
	Cond value.Value
	// True condition target branch.
	TargetTrue value.Value // *ir.Block
	// False condition target branch.
	TargetFalse value.Value // *ir.Block

	// extra.

	// Successor basic blocks of the terminator.
	Successors []*Block
	// (optional) Metadata.
	Metadata
}

// NewCondBr returns a new conditional br terminator based on the given
// branching condition and conditional target basic blocks.
func NewCondBr(cond value.Value, targetTrue, targetFalse *Block) *TermCondBr {
	return &TermCondBr{Cond: cond, TargetTrue: targetTrue, TargetFalse: targetFalse}
}

// Succs returns the successor basic blocks of the terminator.
func (term *TermCondBr) Succs() []*Block {
	// Cache successors if not present.
	if term.Successors == nil {
		term.Successors = []*Block{term.TargetTrue.(*Block), term.TargetFalse.(*Block)}
	}
	return term.Successors
}

// LLString returns the LLVM syntax representation of the terminator.
//
// 'br' CondTyp=IntType Cond=Value ',' TargetTrue=Label ',' TargetFalse=Label Metadata=(',' MetadataAttachment)+?
func (term *TermCondBr) LLString() string {
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "br %s, %s, %s", term.Cond, term.TargetTrue, term.TargetFalse)
	for _, md := range term.Metadata {
		fmt.Fprintf(buf, ", %s", md)
	}
	return buf.String()
}

// --- [ switch ] --------------------------------------------------------------

// TermSwitch is an LLVM IR switch terminator.
type TermSwitch struct {
	// Control variable.
	X value.Value
	// Default target branch.
	TargetDefault value.Value // *ir.Block
	// Switch cases.
	Cases []*Case

	// extra.

	// Successor basic blocks of the terminator.
	Successors []*Block
	// (optional) Metadata.
	Metadata
}

// NewSwitch returns a new switch terminator based on the given control
// variable, default target basic block and switch cases.
func NewSwitch(x value.Value, targetDefault *Block, cases ...*Case) *TermSwitch {
	return &TermSwitch{X: x, TargetDefault: targetDefault, Cases: cases}
}

// Succs returns the successor basic blocks of the terminator.
func (term *TermSwitch) Succs() []*Block {
	// Cache successors if not present.
	if term.Successors == nil {
		succs := make([]*Block, 0, 1+len(term.Cases))
		succs = append(succs, term.TargetDefault.(*Block))
		for _, c := range term.Cases {
			succs = append(succs, c.Target.(*Block))
		}
		term.Successors = succs
	}
	return term.Successors
}

// LLString returns the LLVM syntax representation of the terminator.
//
// 'switch' X=TypeValue ',' Default=Label '[' Cases=Case* ']' Metadata=(',' MetadataAttachment)+?
func (term *TermSwitch) LLString() string {
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "switch %s, %s [\n", term.X, term.TargetDefault)
	for _, c := range term.Cases {
		fmt.Fprintf(buf, "\t\t%s\n", c)
	}
	buf.WriteString("\t]")
	for _, md := range term.Metadata {
		fmt.Fprintf(buf, ", %s", md)
	}
	return buf.String()
}

// ~~~ [ Switch case ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Case is a switch case.
type Case struct {
	// Case comparand.
	X value.Value // constant.Constant (integer constant or integer constant expression)
	// Case target branch.
	Target value.Value // *ir.Block
}

// NewCase returns a new switch case based on the given case comparand and
// target basic block.
func NewCase(x constant.Constant, target *Block) *Case {
	return &Case{X: x, Target: target}
}

// String returns the string representation of the switch case.
func (c *Case) String() string {
	// X=TypeConst ',' Target=Label
	return fmt.Sprintf("%s, %s", c.X, c.Target)
}

// --- [ indirectbr ] ----------------------------------------------------------

// TermIndirectBr is an LLVM IR indirectbr terminator.
type TermIndirectBr struct {
	// Target address.
	Addr value.Value // blockaddress
	// Set of valid target basic blocks.
	ValidTargets []value.Value // slice of *ir.Block

	// extra.

	// Successor basic blocks of the terminator.
	Successors []*Block
	// (optional) Metadata.
	Metadata
}

// NewIndirectBr returns a new indirectbr terminator based on the given target
// address (derived from a blockaddress constant) and set of valid target basic
// blocks.
func NewIndirectBr(addr constant.Constant, validTargets ...*Block) *TermIndirectBr {
	// convert validTargets slice to []value.Value.
	var targets []value.Value
	for _, target := range validTargets {
		targets = append(targets, target)
	}
	return &TermIndirectBr{Addr: addr, ValidTargets: targets}
}

// Succs returns the successor basic blocks of the terminator.
func (term *TermIndirectBr) Succs() []*Block {
	// Cache successors if not present.
	if term.Successors == nil {
		// convert ValidTargets slice to []*ir.Block.
		for _, target := range term.ValidTargets {
			term.Successors = append(term.Successors, target.(*Block))
		}
	}
	return term.Successors
}

// LLString returns the LLVM syntax representation of the terminator.
//
// 'indirectbr' Addr=TypeValue ',' '[' ValidTargets=(Label separator ',')* ']' Metadata=(',' MetadataAttachment)+?
func (term *TermIndirectBr) LLString() string {
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "indirectbr %s, [", term.Addr)
	for i, target := range term.ValidTargets {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(target.String())
	}
	buf.WriteString("]")
	for _, md := range term.Metadata {
		fmt.Fprintf(buf, ", %s", md)
	}
	return buf.String()
}

// --- [ invoke ] --------------------------------------------------------------

// TermInvoke is an LLVM IR invoke terminator.
type TermInvoke struct {
	// Name of local variable associated with the result.
	LocalIdent
	// Invokee (callee function).
	// TODO: specify the set of underlying types of Invokee.
	Invokee value.Value
	// Function arguments.
	//
	// Arg has one of the following underlying types:
	//    value.Value
	//    TODO: add metadata value?
	Args []value.Value
	// Normal control flow return point.
	NormalRetTarget value.Value // *ir.Block
	// Exception control flow return point.
	ExceptionRetTarget value.Value // *ir.Block

	// extra.

	// Type of result produced by the terminator.
	Typ types.Type
	// Successor basic blocks of the terminator.
	Successors []*Block
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

// NewInvoke returns a new invoke terminator based on the given invokee,
// function arguments and control flow return points for normal and exceptional
// execution.
//
// TODO: specify the set of underlying types of invokee.
func NewInvoke(invokee value.Value, args []value.Value, normalRetTarget, exceptionRetTarget *Block) *TermInvoke {
	term := &TermInvoke{Invokee: invokee, Args: args, NormalRetTarget: normalRetTarget, ExceptionRetTarget: exceptionRetTarget}
	// Compute type.
	term.Type()
	return term
}

// String returns the LLVM syntax representation of the terminator as a type-
// value pair.
func (term *TermInvoke) String() string {
	return fmt.Sprintf("%s %s", term.Type(), term.Ident())
}

// Type returns the type of the terminator.
func (term *TermInvoke) Type() types.Type {
	// Cache type if not present.
	if term.Typ == nil {
		sig := term.Sig()
		term.Typ = sig.RetType
	}
	return term.Typ
}

// Succs returns the successor basic blocks of the terminator.
func (term *TermInvoke) Succs() []*Block {
	// Cache successors if not present.
	if term.Successors == nil {
		term.Successors = []*Block{term.NormalRetTarget.(*Block), term.ExceptionRetTarget.(*Block)}
	}
	return term.Successors
}

// LLString returns the LLVM syntax representation of the terminator.
//
// 'invoke' CallingConvopt ReturnAttrs=ReturnAttribute* AddrSpaceopt Typ=Type Invokee=Value '(' Args ')' FuncAttrs=FuncAttribute* OperandBundles=('[' (OperandBundle separator ',')+ ']')? 'to' NormalRetTarget=Label 'unwind' ExceptionRetTarget=Label Metadata=(',' MetadataAttachment)+?
func (term *TermInvoke) LLString() string {
	buf := &strings.Builder{}
	if !term.Type().Equal(types.Void) {
		fmt.Fprintf(buf, "%s = ", term.Ident())
	}
	buf.WriteString("invoke")
	if term.CallingConv != enum.CallingConvNone {
		fmt.Fprintf(buf, " %s", callingConvString(term.CallingConv))
	}
	for _, attr := range term.ReturnAttrs {
		fmt.Fprintf(buf, " %s", attr)
	}
	// (optional) Address space.
	if term.AddrSpace != 0 {
		fmt.Fprintf(buf, " %s", term.AddrSpace)
	}
	// Use function signature instead of return type for variadic functions.
	invokeeType := term.Type()
	if sig := term.Sig(); sig.Variadic {
		invokeeType = sig
	}
	fmt.Fprintf(buf, " %s %s(", invokeeType, term.Invokee.Ident())
	for i, arg := range term.Args {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(arg.String())
	}
	buf.WriteString(")")
	for _, attr := range term.FuncAttrs {
		fmt.Fprintf(buf, " %s", attr)
	}
	if len(term.OperandBundles) > 0 {
		buf.WriteString(" [ ")
		for i, operandBundle := range term.OperandBundles {
			if i != 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(operandBundle.String())
		}
		buf.WriteString(" ]")
	}
	fmt.Fprintf(buf, "\n\t\tto %s unwind %s", term.NormalRetTarget, term.ExceptionRetTarget)
	for _, md := range term.Metadata {
		fmt.Fprintf(buf, ", %s", md)
	}
	return buf.String()
}

// Sig returns the function signature of the invokee.
func (term *TermInvoke) Sig() *types.FuncType {
	t, ok := term.Invokee.Type().(*types.PointerType)
	if !ok {
		panic(fmt.Errorf("invalid invokee type; expected *types.PointerType, got %T", term.Invokee.Type()))
	}
	sig, ok := t.ElemType.(*types.FuncType)
	if !ok {
		panic(fmt.Errorf("invalid invokee type; expected *types.FuncType, got %T", t.ElemType))
	}
	return sig
}

// --- [ callbr ] --------------------------------------------------------------

// TermCallBr is an LLVM IR callbr terminator.
type TermCallBr struct {
	// Name of local variable associated with the result.
	LocalIdent
	// Callee function.
	// TODO: specify the set of underlying types of Callee.
	Callee value.Value
	// Function arguments.
	//
	// Arg has one of the following underlying types:
	//    value.Value
	//    TODO: add metadata value?
	Args []value.Value
	// Normal control flow return point.
	NormalRetTarget value.Value // *ir.Block
	// Other control flow return points.
	OtherRetTargets []value.Value // slice of *ir.Block

	// extra.

	// Type of result produced by the terminator.
	Typ types.Type
	// Successor basic blocks of the terminator.
	Successors []*Block
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

// NewCallBr returns a new callbr terminator based on the given callee, function
// arguments and control flow return points for normal and exceptional
// execution.
//
// TODO: specify the set of underlying types of callee.
func NewCallBr(callee value.Value, args []value.Value, normalRetTarget *Block, otherRetTargets ...*Block) *TermCallBr {
	// Convert otherRetTargets slice to []value.Value.
	var otherRets []value.Value
	for _, otherRetTarget := range otherRetTargets {
		otherRets = append(otherRets, otherRetTarget)
	}
	term := &TermCallBr{Callee: callee, Args: args, NormalRetTarget: normalRetTarget, OtherRetTargets: otherRets}
	// Compute type.
	term.Type()
	return term
}

// String returns the LLVM syntax representation of the terminator as a type-
// value pair.
func (term *TermCallBr) String() string {
	return fmt.Sprintf("%s %s", term.Type(), term.Ident())
}

// Type returns the type of the terminator.
func (term *TermCallBr) Type() types.Type {
	// Cache type if not present.
	if term.Typ == nil {
		sig := term.Sig()
		term.Typ = sig.RetType
	}
	return term.Typ
}

// Succs returns the successor basic blocks of the terminator.
func (term *TermCallBr) Succs() []*Block {
	// Cache successors if not present.
	if term.Successors == nil {
		term.Successors = []*Block{term.NormalRetTarget.(*Block)}
		// Convert OtherRetTargets slice to []*ir.Block.
		for _, otherRetTarget := range term.OtherRetTargets {
			term.Successors = append(term.Successors, otherRetTarget.(*Block))
		}
	}
	return term.Successors
}

// LLString returns the LLVM syntax representation of the terminator.
//
// 'callbr' CallingConvopt ReturnAttrs=ReturnAttribute* AddrSpaceopt Typ=Type Callee=Value '(' Args ')' FuncAttrs=FuncAttribute* OperandBundles=('[' (OperandBundle separator ',')+ ']')? 'to' NormalRetTarget=Label '[' OtherRetTargets=(Label separator ',')* ']' Metadata=(',' MetadataAttachment)+?
func (term *TermCallBr) LLString() string {
	buf := &strings.Builder{}
	if !term.Type().Equal(types.Void) {
		fmt.Fprintf(buf, "%s = ", term.Ident())
	}
	buf.WriteString("callbr")
	if term.CallingConv != enum.CallingConvNone {
		fmt.Fprintf(buf, " %s", callingConvString(term.CallingConv))
	}
	for _, attr := range term.ReturnAttrs {
		fmt.Fprintf(buf, " %s", attr)
	}
	// (optional) Address space.
	if term.AddrSpace != 0 {
		fmt.Fprintf(buf, " %s", term.AddrSpace)
	}
	// Use function signature instead of return type for variadic functions.
	calleeType := term.Type()
	if sig := term.Sig(); sig.Variadic {
		calleeType = sig
	}
	fmt.Fprintf(buf, " %s %s(", calleeType, term.Callee.Ident())
	for i, arg := range term.Args {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(arg.String())
	}
	buf.WriteString(")")
	for _, attr := range term.FuncAttrs {
		fmt.Fprintf(buf, " %s", attr)
	}
	if len(term.OperandBundles) > 0 {
		buf.WriteString(" [ ")
		for i, operandBundle := range term.OperandBundles {
			if i != 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(operandBundle.String())
		}
		buf.WriteString(" ]")
	}
	fmt.Fprintf(buf, "\n\t\tto %s [", term.NormalRetTarget)
	for i, otherRetTarget := range term.OtherRetTargets {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(otherRetTarget.String())
	}
	buf.WriteString("]")
	for _, md := range term.Metadata {
		fmt.Fprintf(buf, ", %s", md)
	}
	return buf.String()
}

// Sig returns the function signature of the callee.
func (term *TermCallBr) Sig() *types.FuncType {
	t, ok := term.Callee.Type().(*types.PointerType)
	if !ok {
		panic(fmt.Errorf("invalid callee type; expected *types.PointerType, got %T", term.Callee.Type()))
	}
	sig, ok := t.ElemType.(*types.FuncType)
	if !ok {
		panic(fmt.Errorf("invalid callee type; expected *types.FuncType, got %T", t.ElemType))
	}
	return sig
}

// --- [ resume ] --------------------------------------------------------------

// TermResume is an LLVM IR resume terminator.
type TermResume struct {
	// Exception argument to propagate.
	X value.Value

	// extra.

	// (optional) Metadata.
	Metadata
}

// NewResume returns a new resume terminator based on the given exception
// argument to propagate.
func NewResume(x value.Value) *TermResume {
	return &TermResume{X: x}
}

// Succs returns the successor basic blocks of the terminator.
func (term *TermResume) Succs() []*Block {
	// no successors.
	return nil
}

// LLString returns the LLVM syntax representation of the terminator.
//
// 'resume' X=TypeValue Metadata=(',' MetadataAttachment)+?
func (term *TermResume) LLString() string {
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "resume %s", term.X)
	for _, md := range term.Metadata {
		fmt.Fprintf(buf, ", %s", md)
	}
	return buf.String()
}

// --- [ catchswitch ] ---------------------------------------------------------

// TermCatchSwitch is an LLVM IR catchswitch terminator.
type TermCatchSwitch struct {
	// Name of local variable associated with the result.
	LocalIdent
	// Parent exception pad.
	ParentPad value.Value // ir.ExceptionPad
	// Exception handlers.
	Handlers []value.Value // []*ir.Block
	// Optional default target basic block to transfer control flow to; or nil to
	// unwind to caller function.
	DefaultUnwindTarget value.Value // *ir.Block or nil

	// extra.

	// Successor basic blocks of the terminator.
	Successors []*Block
	// (optional) Metadata.
	Metadata
}

// NewCatchSwitch returns a new catchswitch terminator based on the given parent
// exception pad, exception handlers and optional default unwind target. If
// defaultUnwindTarget is nil, catchswitch unwinds to caller function.
func NewCatchSwitch(parentPad ExceptionPad, handlers []*Block, defaultUnwindTarget *Block) *TermCatchSwitch {
	// convert handlers slice to []value.Value.
	var hs []value.Value
	for _, handler := range handlers {
		hs = append(hs, handler)
	}
	term := &TermCatchSwitch{ParentPad: parentPad, Handlers: hs}
	if defaultUnwindTarget != nil {
		// Note: since DefaultUnwindTarget is an interface we have to be careful
		// with typed nil values (e.g. `(*ir.Block)(nil)`). This is to ensure that
		// DefaultUnwindTarget is nil and not `{Type: ir.Block, Value: nil}`.
		//
		// ref: https://golang.org/doc/faq#nil_error
		term.DefaultUnwindTarget = defaultUnwindTarget
	}
	return term
}

// String returns the LLVM syntax representation of the terminator as a type-
// value pair.
func (term *TermCatchSwitch) String() string {
	return fmt.Sprintf("%s %s", term.Type(), term.Ident())
}

// Type returns the type of the terminator.
func (term *TermCatchSwitch) Type() types.Type {
	return types.Token
}

// Succs returns the successor basic blocks of the terminator.
func (term *TermCatchSwitch) Succs() []*Block {
	// Cache successors if not present.
	if term.Successors == nil {
		// convert Handlers slice to []*ir.Block.
		for _, handler := range term.Handlers {
			term.Successors = append(term.Successors, handler.(*Block))
		}
		if defaultUnwindTarget, ok := term.DefaultUnwindTarget.(*Block); ok {
			term.Successors = append(term.Successors, defaultUnwindTarget)
		}
	}
	return term.Successors
}

// LLString returns the LLVM syntax representation of the terminator.
//
// 'catchswitch' 'within' ParentPad=ExceptionPad '[' Handlers=Handlers ']' 'unwind' DefaultUnwindTarget=UnwindTarget Metadata=(',' MetadataAttachment)+?
func (term *TermCatchSwitch) LLString() string {
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "%s = ", term.Ident())
	fmt.Fprintf(buf, "catchswitch within %s [", term.ParentPad.Ident())
	for i, handler := range term.Handlers {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(handler.String())
	}
	buf.WriteString("] unwind ")
	if term.DefaultUnwindTarget != nil {
		buf.WriteString(term.DefaultUnwindTarget.String())
	} else {
		buf.WriteString("to caller")
	}
	for _, md := range term.Metadata {
		fmt.Fprintf(buf, ", %s", md)
	}
	return buf.String()
}

// --- [ catchret ] ------------------------------------------------------------

// TermCatchRet is an LLVM IR catchret terminator, which catches an in-flight
// exception from CatchPad and returns control flow to normal at Target.
type TermCatchRet struct {
	// Exit catchpad.
	CatchPad value.Value // *ir.InstCatchPad
	// Target basic block to transfer control flow to.
	Target value.Value // *ir.Block

	// extra.

	// Successor basic blocks of the terminator.
	Successors []*Block
	// (optional) Metadata.
	Metadata
}

// NewCatchRet returns a new catchret terminator based on the given exit
// catchpad and target basic block.
func NewCatchRet(catchPad *InstCatchPad, target *Block) *TermCatchRet {
	return &TermCatchRet{CatchPad: catchPad, Target: target}
}

// Succs returns the successor basic blocks of the terminator.
func (term *TermCatchRet) Succs() []*Block {
	// Cache successors if not present.
	if term.Successors == nil {
		term.Successors = []*Block{term.Target.(*Block)}
	}
	return term.Successors
}

// LLString returns the LLVM syntax representation of the terminator.
//
// 'catchret' 'from' CatchPad=Value 'to' Target=Label Metadata=(',' MetadataAttachment)+?
func (term *TermCatchRet) LLString() string {
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "catchret from %s to %s", term.CatchPad.Ident(), term.Target)
	for _, md := range term.Metadata {
		fmt.Fprintf(buf, ", %s", md)
	}
	return buf.String()
}

// --- [ cleanupret ] ----------------------------------------------------------

// TermCleanupRet is an LLVM IR cleanupret terminator, which indicates that the
// personality function of a cleanuppad has finished and transfers control flow
// to an optional target basic block or unwinds to the caller function.
type TermCleanupRet struct {
	// Exit cleanuppad.
	CleanupPad value.Value // *ir.InstCleanupPad
	// Optional target basic block to transfer control flow to; or nil to unwind
	// to caller function.
	UnwindTarget value.Value // *ir.Block or nil

	// extra.

	// Successor basic blocks of the terminator.
	Successors []*Block
	// (optional) Metadata.
	Metadata
}

// NewCleanupRet returns a new cleanupret terminator based on the given exit
// cleanuppad and optional unwind target. If unwindTarget is nil, cleanupret
// unwinds to caller function.
func NewCleanupRet(cleanupPad *InstCleanupPad, unwindTarget *Block) *TermCleanupRet {
	term := &TermCleanupRet{CleanupPad: cleanupPad}
	if unwindTarget != nil {
		// Note: since UnwindTarget is an interface we have to be careful
		// with typed nil values (e.g. `(*ir.Block)(nil)`). This is to ensure that
		// UnwindTarget is nil and not `{Type: ir.Block, Value: nil}`.
		//
		// ref: https://golang.org/doc/faq#nil_error
		term.UnwindTarget = unwindTarget
	}
	return term
}

// Succs returns the successor basic blocks of the terminator.
func (term *TermCleanupRet) Succs() []*Block {
	// Cache successors if not present.
	if term.Successors == nil {
		if unwindTarget, ok := term.UnwindTarget.(*Block); ok {
			term.Successors = []*Block{unwindTarget}
		} else {
			term.Successors = []*Block{}
		}
	}
	return term.Successors
}

// LLString returns the LLVM syntax representation of the terminator.
//
// 'cleanupret' 'from' CleanupPad=Value 'unwind' UnwindTarget Metadata=(',' MetadataAttachment)+?
func (term *TermCleanupRet) LLString() string {
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "cleanupret from %s unwind ", term.CleanupPad.Ident())
	if term.UnwindTarget != nil {
		buf.WriteString(term.UnwindTarget.String())
	} else {
		buf.WriteString("to caller")
	}
	for _, md := range term.Metadata {
		fmt.Fprintf(buf, ", %s", md)
	}
	return buf.String()
}

// --- [ unreachable ] ---------------------------------------------------------

// TermUnreachable is an LLVM IR unreachable terminator.
type TermUnreachable struct {
	// extra.

	// (optional) Metadata.
	Metadata
}

// NewUnreachable returns a new unreachable terminator.
func NewUnreachable() *TermUnreachable {
	return &TermUnreachable{}
}

// Succs returns the successor basic blocks of the terminator.
func (term *TermUnreachable) Succs() []*Block {
	// no successors.
	return nil
}

// LLString returns the LLVM syntax representation of the terminator.
//
// 'unreachable' Metadata=(',' MetadataAttachment)+?
func (term *TermUnreachable) LLString() string {
	buf := &strings.Builder{}
	buf.WriteString("unreachable")
	for _, md := range term.Metadata {
		fmt.Fprintf(buf, ", %s", md)
	}
	return buf.String()
}
