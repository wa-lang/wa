package types

import (
	"fmt"
	"math/big"

	"github.com/wa-lang/wa/internal/3rdparty/llir"
	"github.com/wa-lang/wa/internal/3rdparty/llir/llconstant"
	"github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
	"github.com/wa-lang/wa/internal/3rdparty/llir/llvalue"
	"github.com/wa-lang/wa/internal/backends/compiler_ll/llutil"
)

type Type interface {
	LLVM() lltypes.Type
	Name() string

	// Size of type in bytes
	Size() int64

	AddMethod(string, *Method)
	GetMethod(string) (*Method, bool)

	Zero(*llir.Block, llvalue.Value)

	IsSigned() bool
}

type backingType struct {
	methods map[string]*Method
}

func (b *backingType) AddMethod(name string, method *Method) {
	if b.methods == nil {
		b.methods = make(map[string]*Method)
	}
	b.methods[name] = method
}

func (b *backingType) GetMethod(name string) (*Method, bool) {
	m, ok := b.methods[name]
	return m, ok
}

func (backingType) Size() int64 {
	panic("Type does not have size set")
}

func (backingType) Zero(*llir.Block, llvalue.Value) {
	// NOOP
}

func (backingType) IsSigned() bool {
	return false
}

type Struct struct {
	backingType

	Members       map[string]Type
	MemberIndexes map[string]int

	IsHeapAllocated bool

	SourceName string
	Type       lltypes.Type
}

func (s Struct) LLVM() lltypes.Type {
	return s.Type
}

func (s Struct) Name() string {
	return fmt.Sprintf("struct(%s)", s.SourceName)
}

func (s Struct) Zero(block *llir.Block, alloca llvalue.Value) {
	for key, valType := range s.Members {
		ptr := block.NewGetElementPtr(llutil.PtrElemType(alloca), alloca,
			llconstant.NewInt(lltypes.I32, 0),
			llconstant.NewInt(lltypes.I32, int64(s.MemberIndexes[key])),
		)
		valType.Zero(block, ptr)
	}
}

func (s Struct) Size() int64 {
	var sum int64
	for _, valType := range s.Members {
		sum += valType.Size()
	}
	return sum
}

type Method struct {
	backingType

	Function        *Function
	LlvmFunction    llvalue.Named
	PointerReceiver bool
	MethodName      string
}

func (m Method) LLVM() lltypes.Type {
	return m.Function.LLVM()
}

func (m Method) Name() string {
	return m.MethodName
}

type Function struct {
	backingType

	// LlvmFunction llvalue.Named
	FuncType lltypes.Type

	// The return type of the LLVM function (is always 1)
	LlvmReturnType Type
	// Return types of the Tre function
	ReturnTypes []Type

	IsVariadic    bool
	ArgumentTypes []Type
	IsExternal    bool

	// Is used when calling an interface method
	JumpFunction *llir.Func
}

func (f Function) LLVM() lltypes.Type {
	return f.FuncType
}

func (f Function) Name() string {
	return "func"
}

type BoolType struct {
	backingType
}

func (BoolType) LLVM() lltypes.Type {
	return lltypes.I1
}

func (BoolType) Name() string {
	return "bool"
}

func (BoolType) Size() int64 {
	return 1
}

func (b BoolType) Zero(block *llir.Block, alloca llvalue.Value) {
	block.NewStore(llconstant.NewInt(lltypes.I1, 0), alloca)
}

type VoidType struct {
	backingType
}

func (VoidType) LLVM() lltypes.Type {
	return lltypes.Void
}

func (VoidType) Name() string {
	return "void"
}

func (VoidType) Size() int64 {
	return 0
}

type Int struct {
	backingType

	Type     *lltypes.IntType
	TypeName string
	TypeSize int64
	Signed   bool
}

func (i Int) LLVM() lltypes.Type {
	return i.Type
}

func (i Int) Name() string {
	return i.TypeName
}

func (i Int) Size() int64 {
	return i.TypeSize
}

func (i Int) Zero(block *llir.Block, alloca llvalue.Value) {
	b := big.NewInt(0)
	if !i.IsSigned() {
		b.SetUint64(0)
	}

	c := &llconstant.Int{
		Typ: i.Type,
		X:   b,
	}

	block.NewStore(c, alloca)
}

func (i Int) IsSigned() bool {
	return i.Signed
}

type StringType struct {
	backingType
	Type lltypes.Type
}

// Populated by compiler.go
var ModuleStringType lltypes.Type
var EmptyStringConstant *llir.Global

func (StringType) LLVM() lltypes.Type {
	return ModuleStringType
}

func (StringType) Name() string {
	return "string"
}

func (StringType) Size() int64 {
	return 16
}

func (s StringType) Zero(block *llir.Block, alloca llvalue.Value) {
	lenPtr := block.NewGetElementPtr(llutil.PtrElemType(alloca), alloca, llconstant.NewInt(lltypes.I32, 0), llconstant.NewInt(lltypes.I32, 0))
	backingDataPtr := block.NewGetElementPtr(llutil.PtrElemType(alloca), alloca, llconstant.NewInt(lltypes.I32, 0), llconstant.NewInt(lltypes.I32, 1))
	block.NewStore(llconstant.NewInt(lltypes.I64, 0), lenPtr)
	block.NewStore(llutil.StrToi8Ptr(block, EmptyStringConstant), backingDataPtr)
}

type Array struct {
	backingType
	Type     Type
	Len      uint64
	LlvmType lltypes.Type
}

func (a Array) LLVM() lltypes.Type {
	return a.LlvmType
}

func (a Array) Name() string {
	return "array"
}

func (a Array) Zero(block *llir.Block, alloca llvalue.Value) {
	for i := uint64(0); i < a.Len; i++ {
		ptr := block.NewGetElementPtr(llutil.PtrElemType(alloca), alloca, llconstant.NewInt(lltypes.I64, 0), llconstant.NewInt(lltypes.I64, int64(i)))
		a.Type.Zero(block, ptr)
	}
}

type Slice struct {
	backingType
	Type     Type // type of the items in the slice []int => int
	LlvmType lltypes.Type
}

func (s Slice) LLVM() lltypes.Type {
	return s.LlvmType
}

func (Slice) Name() string {
	return "slice"
}

func (Slice) Size() int64 {
	return 3*4 + 8 // 3 int32s and a pointer
}

func (s Slice) SliceZero(block *llir.Block, mallocFunc llvalue.Named, initCap int, emptySlice llvalue.Value) {
	// The cap must always be larger than 0
	// Use 2 as the default value
	if initCap < 2 {
		initCap = 2
	}

	len := block.NewGetElementPtr(llutil.PtrElemType(emptySlice), emptySlice, llconstant.NewInt(lltypes.I32, 0), llconstant.NewInt(lltypes.I32, 0))
	len.SetName("len")
	cap := block.NewGetElementPtr(llutil.PtrElemType(emptySlice), emptySlice, llconstant.NewInt(lltypes.I32, 0), llconstant.NewInt(lltypes.I32, 1))
	cap.SetName("cap")
	offset := block.NewGetElementPtr(llutil.PtrElemType(emptySlice), emptySlice, llconstant.NewInt(lltypes.I32, 0), llconstant.NewInt(lltypes.I32, 2))
	offset.SetName("offset")
	backingArray := block.NewGetElementPtr(llutil.PtrElemType(emptySlice), emptySlice, llconstant.NewInt(lltypes.I32, 0), llconstant.NewInt(lltypes.I32, 3))
	backingArray.SetName("backing")

	block.NewStore(llconstant.NewInt(lltypes.I32, 0), len)
	block.NewStore(llconstant.NewInt(lltypes.I32, int64(initCap)), cap)
	block.NewStore(llconstant.NewInt(lltypes.I32, 0), offset)

	mallocatedSpaceRaw := block.NewCall(mallocFunc, llconstant.NewInt(lltypes.I64, int64(initCap)*s.Type.Size()))
	mallocatedSpaceRaw.SetName("slicezero")
	bitcasted := block.NewBitCast(mallocatedSpaceRaw, lltypes.NewPointer(s.Type.LLVM()))
	block.NewStore(bitcasted, backingArray)
}

type Pointer struct {
	backingType

	Type                  Type
	IsNonAllocDereference bool

	LlvmType lltypes.Type
}

func (p Pointer) LLVM() lltypes.Type {
	return lltypes.NewPointer(p.Type.LLVM())
}

func (p Pointer) Name() string {
	return fmt.Sprintf("pointer(%s)", p.Type.Name())
}

func (p Pointer) Size() int64 {
	return 8
}

// MultiValue is used when returning multiple values from a function
type MultiValue struct {
	backingType
	Types []Type
}

func (m MultiValue) Name() string {
	return "multivalue"
}

func (m MultiValue) LLVM() lltypes.Type {
	panic("MutliValue has no LLVM type")
}

type UntypedConstantNumber struct {
	backingType
}

func (m UntypedConstantNumber) Name() string {
	return "UntypedConstantNumber"
}

func (m UntypedConstantNumber) LLVM() lltypes.Type {
	panic("UntypedConstantNumber has no LLVM type")
}
