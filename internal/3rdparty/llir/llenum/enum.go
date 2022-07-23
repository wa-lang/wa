// Package enum defines enumerate types of LLVM IR.
package llenum

//go:generate stringer -linecomment -type AtomicOp

// AtomicOp is an AtomicRMW binary operation.
type AtomicOp uint8

// AtomicRMW binary operations.
const (
	AtomicOpAdd  AtomicOp = iota + 1 // add
	AtomicOpAnd                      // and
	AtomicOpFAdd                     // fadd
	AtomicOpFSub                     // fsub
	AtomicOpMax                      // max
	AtomicOpMin                      // min
	AtomicOpNAnd                     // nand
	AtomicOpOr                       // or
	AtomicOpSub                      // sub
	AtomicOpUMax                     // umax
	AtomicOpUMin                     // umin
	AtomicOpXChg                     // xchg
	AtomicOpXor                      // xor
)

//go:generate stringer -linecomment -type AtomicOrdering

// AtomicOrdering is an atomic ordering attribute.
type AtomicOrdering uint8

// Atomic ordering attributes.
const (
	AtomicOrderingNone      AtomicOrdering = iota // none
	AtomicOrderingAcqRel                          // acq_rel
	AtomicOrderingAcquire                         // acquire
	AtomicOrderingMonotonic                       // monotonic
	AtomicOrderingRelease                         // release
	AtomicOrderingSeqCst                          // seq_cst
	AtomicOrderingUnordered                       // unordered
)

//go:generate stringer -linecomment -type CallingConv

// CallingConv is a calling convention.
type CallingConv uint16

// Calling conventions.
//
// From include/llvm/IR/CallingConv.h
const (
	CallingConvNone CallingConv = 0 // none
	// Note, C calling convention is defined as 0 in LLVM. To have the zero-value
	// calling convention mean no calling convention, re-define C calling
	// convention as 1, and use 0 for none.
	CallingConvC            CallingConv = 1  // ccc
	CallingConvFast         CallingConv = 8  // fastcc
	CallingConvCold         CallingConv = 9  // coldcc
	CallingConvGHC          CallingConv = 10 // ghccc
	CallingConvHiPE         CallingConv = 11 // cc 11
	CallingConvWebKitJS     CallingConv = 12 // webkit_jscc
	CallingConvAnyReg       CallingConv = 13 // anyregcc
	CallingConvPreserveMost CallingConv = 14 // preserve_mostcc
	CallingConvPreserveAll  CallingConv = 15 // preserve_allcc
	CallingConvSwift        CallingConv = 16 // swiftcc
	CallingConvCXXFastTLS   CallingConv = 17 // cxx_fast_tlscc
	CallingConvTail         CallingConv = 18 // tailcc
	CallingConvCFGuardCheck CallingConv = 19 // cfguard_checkcc

	// Start of target-specific calling conventions.
	CallingConvFirstTarget = CallingConvX86StdCall

	CallingConvX86StdCall    CallingConv = 64 // x86_stdcallcc
	CallingConvX86FastCall   CallingConv = 65 // x86_fastcallcc
	CallingConvARM_APCS      CallingConv = 66 // arm_apcscc
	CallingConvARM_AAPCS     CallingConv = 67 // arm_aapcscc
	CallingConvARM_AAPCS_VFP CallingConv = 68 // arm_aapcs_vfpcc
	CallingConvMSP430Intr    CallingConv = 69 // msp430_intrcc
	CallingConvX86ThisCall   CallingConv = 70 // x86_thiscallcc
	CallingConvPTXKernel     CallingConv = 71 // ptx_kernel
	CallingConvPTXDevice     CallingConv = 72 // ptx_device

	CallingConvSPIRFunc             CallingConv = 75  // spir_func
	CallingConvSPIRKernel           CallingConv = 76  // spir_kernel
	CallingConvIntelOCL_BI          CallingConv = 77  // intel_ocl_bicc
	CallingConvX86_64SysV           CallingConv = 78  // x86_64_sysvcc
	CallingConvWin64                CallingConv = 79  // win64cc
	CallingConvX86VectorCall        CallingConv = 80  // x86_vectorcallcc
	CallingConvHHVM                 CallingConv = 81  // hhvmcc
	CallingConvHHVM_C               CallingConv = 82  // hhvm_ccc
	CallingConvX86Intr              CallingConv = 83  // x86_intrcc
	CallingConvAVRIntr              CallingConv = 84  // avr_intrcc
	CallingConvAVRSignal            CallingConv = 85  // avr_signalcc
	CallingConvAVRBuiltin           CallingConv = 86  // cc 86
	CallingConvAMDGPU_VS            CallingConv = 87  // amdgpu_vs
	CallingConvAMDGPU_GS            CallingConv = 88  // amdgpu_gs
	CallingConvAMDGPU_PS            CallingConv = 89  // amdgpu_ps
	CallingConvAMDGPU_CS            CallingConv = 90  // amdgpu_cs
	CallingConvAMDGPUKernel         CallingConv = 91  // amdgpu_kernel
	CallingConvX86RegCall           CallingConv = 92  // x86_regcallcc
	CallingConvAMDGPU_HS            CallingConv = 93  // amdgpu_hs
	CallingConvMSP430Builtin        CallingConv = 94  // cc 94
	CallingConvAMDGPU_LS            CallingConv = 95  // amdgpu_ls
	CallingConvAMDGPU_ES            CallingConv = 96  // amdgpu_es
	CallingConvAArch64VectorCall    CallingConv = 97  // aarch64_vector_pcs
	CallingConvAArch64SVEVectorCall CallingConv = 98  // aarch64_sve_vector_pcs
	CallingConvAMDGPUGfx            CallingConv = 100 // amdgpu_gfx
)

//go:generate stringer -linecomment -type ChecksumKind

// ChecksumKind is a checksum algorithm.
type ChecksumKind uint8

// Checksum algorithms.
//
// From include/llvm/IR/DebugInfoMetadata.h
const (
	ChecksumKindMD5  ChecksumKind = 1 // CSK_MD5
	ChecksumKindSHA1 ChecksumKind = 2 // CSK_SHA1
)

//go:generate stringer -linecomment -type ClauseType

// ClauseType specifies the clause type of a landingpad clause.
type ClauseType uint8

// Clause types.
const (
	ClauseTypeCatch  ClauseType = iota + 1 // catch
	ClauseTypeFilter                       // filter
)

//go:generate stringer -type DIFlag

// DIFlag is a debug info flag bitfield.
type DIFlag uint64

// Debug info flags.
//
// From include/llvm/IR/DebugInfoFlags.def (LLVM 9.0)
const (
	DIFlagZero                DIFlag = 0
	DIFlagPrivate             DIFlag = 1
	DIFlagProtected           DIFlag = 2
	DIFlagPublic              DIFlag = 3
	DIFlagFwdDecl             DIFlag = 1 << 2
	DIFlagAppleBlock          DIFlag = 1 << 3
	DIFlagBlockByrefStruct    DIFlag = 1 << 4
	DIFlagVirtual             DIFlag = 1 << 5
	DIFlagArtificial          DIFlag = 1 << 6
	DIFlagExplicit            DIFlag = 1 << 7
	DIFlagPrototyped          DIFlag = 1 << 8
	DIFlagObjcClassComplete   DIFlag = 1 << 9
	DIFlagObjectPointer       DIFlag = 1 << 10
	DIFlagVector              DIFlag = 1 << 11
	DIFlagStaticMember        DIFlag = 1 << 12
	DIFlagLValueReference     DIFlag = 1 << 13
	DIFlagRValueReference     DIFlag = 1 << 14
	DIFlagReserved            DIFlag = 1 << 15
	DIFlagSingleInheritance   DIFlag = 1 << 16
	DIFlagMultipleInheritance DIFlag = 2 << 16
	DIFlagVirtualInheritance  DIFlag = 3 << 16
	DIFlagIntroducedVirtual   DIFlag = 1 << 18
	DIFlagBitField            DIFlag = 1 << 19
	DIFlagNoReturn            DIFlag = 1 << 20
	DIFlagArgumentNotModified DIFlag = 1 << 21
	DIFlagTypePassByValue     DIFlag = 1 << 22
	DIFlagTypePassByReference DIFlag = 1 << 23
	DIFlagEnumClass           DIFlag = 1 << 24
	DIFlagThunk               DIFlag = 1 << 25
	DIFlagNonTrivial          DIFlag = 1 << 26
	DIFlagBigEndian           DIFlag = 1 << 27
	DIFlagLittleEndian        DIFlag = 1 << 28
	DIFlagAllCallsDescribed   DIFlag = 1 << 29

	DIFlagIndirectVirtualBase DIFlag = DIFlagFwdDecl | DIFlagVirtual
	// Mask for accessibility.
	DIFlagAccessibility DIFlag = DIFlagPrivate | DIFlagProtected | DIFlagPublic
	// Mask for inheritance.
	DIFlagPtrToMemberRep DIFlag = DIFlagSingleInheritance | DIFlagMultipleInheritance | DIFlagVirtualInheritance

	// Track first and last debug info flag, used by diFlagsString in
	// ir/metadata/helper.go.
	DIFlagFirst = DIFlagFwdDecl
	DIFlagLast  = DIFlagAllCallsDescribed
)

//go:generate stringer -type DISPFlag

// DISPFlag is a subprogram specific flag bitfield.
type DISPFlag uint64

// Subprogram specific flags.
//
// From include/llvm/IR/DebugInfoFlags.def (LLVM 10.0)
const (
	DISPFlagZero           DISPFlag = 0
	DISPFlagVirtual        DISPFlag = 1
	DISPFlagPureVirtual    DISPFlag = 2
	DISPFlagLocalToUnit    DISPFlag = 1 << 2
	DISPFlagDefinition     DISPFlag = 1 << 3
	DISPFlagOptimized      DISPFlag = 1 << 4
	DISPFlagPure           DISPFlag = 1 << 5
	DISPFlagElemental      DISPFlag = 1 << 6
	DISPFlagRecursive      DISPFlag = 1 << 7
	DISPFlagMainSubprogram DISPFlag = 1 << 8
	// LLVM 10.x flags
	DISPFlagDeleted    DISPFlag = 1 << 9
	DISPFlagObjCDirect DISPFlag = 1 << 11

	// Virtuality and non-virtuality.
	DISPFlagNonvirtual DISPFlag = DISPFlagZero
	DISPFlagVirtuality DISPFlag = DISPFlagVirtual | DISPFlagPureVirtual

	// Track first and last subprogram specific flag, used by diSPFlagsString in
	// ir/metadata/helper.go.
	DISPFlagFirst = DISPFlagVirtual
	DISPFlagLast  = DISPFlagMainSubprogram
)

//go:generate stringer -linecomment -type DLLStorageClass

// DLLStorageClass specifies the DLL storage class of a global identifier.
type DLLStorageClass uint8

// DLL storage classes.
const (
	DLLStorageClassNone      DLLStorageClass = iota // none
	DLLStorageClassDLLExport                        // dllexport
	DLLStorageClassDLLImport                        // dllimport
)

//go:generate stringer -linecomment -type DwarfAttEncoding

// DwarfAttEncoding is a DWARF attribute type encoding.
type DwarfAttEncoding int64

// DWARF attribute type encodings.
//
// From include/llvm/BinaryFormat/Dwarf.def
const (
	// DWARF v2.
	DwarfAttEncodingAddress      DwarfAttEncoding = 0x01 // DW_ATE_address
	DwarfAttEncodingBoolean      DwarfAttEncoding = 0x02 // DW_ATE_boolean
	DwarfAttEncodingComplexFloat DwarfAttEncoding = 0x03 // DW_ATE_complex_float
	DwarfAttEncodingFloat        DwarfAttEncoding = 0x04 // DW_ATE_float
	DwarfAttEncodingSigned       DwarfAttEncoding = 0x05 // DW_ATE_signed
	DwarfAttEncodingSignedChar   DwarfAttEncoding = 0x06 // DW_ATE_signed_char
	DwarfAttEncodingUnsigned     DwarfAttEncoding = 0x07 // DW_ATE_unsigned
	DwarfAttEncodingUnsignedChar DwarfAttEncoding = 0x08 // DW_ATE_unsigned_char
	// DWARF v3.
	DwarfAttEncodingImaginaryFloat DwarfAttEncoding = 0x09 // DW_ATE_imaginary_float
	DwarfAttEncodingPackedDecimal  DwarfAttEncoding = 0x0A // DW_ATE_packed_decimal
	DwarfAttEncodingNumericString  DwarfAttEncoding = 0x0B // DW_ATE_numeric_string
	DwarfAttEncodingEdited         DwarfAttEncoding = 0x0C // DW_ATE_edited
	DwarfAttEncodingSignedFixed    DwarfAttEncoding = 0x0D // DW_ATE_signed_fixed
	DwarfAttEncodingUnsignedFixed  DwarfAttEncoding = 0x0E // DW_ATE_unsigned_fixed
	DwarfAttEncodingDecimalFloat   DwarfAttEncoding = 0x0F // DW_ATE_decimal_float
	// DWARF v4.
	DwarfAttEncodingUTF DwarfAttEncoding = 0x10 // DW_ATE_UTF
	// DWARF v5.
	DwarfAttEncodingUCS   DwarfAttEncoding = 0x11 // DW_ATE_UCS
	DwarfAttEncodingASCII DwarfAttEncoding = 0x12 // DW_ATE_ASCII
)

//go:generate stringer -linecomment -type DwarfCC

// DwarfCC is a DWARF calling convention.
type DwarfCC int64

// DWARF calling conventions.
const (
	DwarfCCNormal  DwarfCC = 0x01 // DW_CC_normal
	DwarfCCProgram DwarfCC = 0x02 // DW_CC_program
	DwarfCCNoCall  DwarfCC = 0x03 // DW_CC_nocall
	// DWARF v5.
	DwarfCCPassByReference DwarfCC = 0x04 // DW_CC_pass_by_reference
	DwarfCCPassByValue     DwarfCC = 0x05 // DW_CC_pass_by_value
	// Vendor extensions.
	DwarfCCGNUBorlandFastcallI386 DwarfCC = 0x41 // DW_CC_GNU_borland_fastcall_i386
	DwarfCCBORLANDSafecall        DwarfCC = 0xB0 // DW_CC_BORLAND_safecall
	DwarfCCBORLANDStdcall         DwarfCC = 0xB1 // DW_CC_BORLAND_stdcall
	DwarfCCBORLANDPascal          DwarfCC = 0xB2 // DW_CC_BORLAND_pascal
	DwarfCCBORLANDMSFastcall      DwarfCC = 0xB3 // DW_CC_BORLAND_msfastcall
	DwarfCCBORLANDMSReturn        DwarfCC = 0xB4 // DW_CC_BORLAND_msreturn
	DwarfCCBORLANDThiscall        DwarfCC = 0xB5 // DW_CC_BORLAND_thiscall
	DwarfCCBORLANDFastcall        DwarfCC = 0xB6 // DW_CC_BORLAND_fastcall
	DwarfCCLLVMVectorcall         DwarfCC = 0xC0 // DW_CC_LLVM_vectorcall
)

//go:generate stringer -linecomment -type DwarfLang

// DwarfLang is a DWARF language.
type DwarfLang int64

// DWARF languages.
//
// From include/llvm/BinaryFormat/Dwarf.def
const (
	// DWARF v2.
	DwarfLangC89       DwarfLang = 0x0001 // DW_LANG_C89
	DwarfLangC         DwarfLang = 0x0002 // DW_LANG_C
	DwarfLangAda83     DwarfLang = 0x0003 // DW_LANG_Ada83
	DwarfLangCPlusPlus DwarfLang = 0x0004 // DW_LANG_C_plus_plus
	DwarfLangCobol74   DwarfLang = 0x0005 // DW_LANG_Cobol74
	DwarfLangCobol85   DwarfLang = 0x0006 // DW_LANG_Cobol85
	DwarfLangFortran77 DwarfLang = 0x0007 // DW_LANG_Fortran77
	DwarfLangFortran90 DwarfLang = 0x0008 // DW_LANG_Fortran90
	DwarfLangPascal83  DwarfLang = 0x0009 // DW_LANG_Pascal83
	DwarfLangModula2   DwarfLang = 0x000A // DW_LANG_Modula2
	// DWARF v3.
	DwarfLangJava         DwarfLang = 0x000B // DW_LANG_Java
	DwarfLangC99          DwarfLang = 0x000C // DW_LANG_C99
	DwarfLangAda95        DwarfLang = 0x000D // DW_LANG_Ada95
	DwarfLangFortran95    DwarfLang = 0x000E // DW_LANG_Fortran95
	DwarfLangPLI          DwarfLang = 0x000F // DW_LANG_PLI
	DwarfLangObjC         DwarfLang = 0x0010 // DW_LANG_ObjC
	DwarfLangObjCPlusPlus DwarfLang = 0x0011 // DW_LANG_ObjC_plus_plus
	DwarfLangUPC          DwarfLang = 0x0012 // DW_LANG_UPC
	DwarfLangD            DwarfLang = 0x0013 // DW_LANG_D
	// DWARF v4.
	DwarfLangPython DwarfLang = 0x0014 // DW_LANG_Python
	// DWARF v5.
	DwarfLangOpenCL       DwarfLang = 0x0015 // DW_LANG_OpenCL
	DwarfLangGo           DwarfLang = 0x0016 // DW_LANG_Go
	DwarfLangModula3      DwarfLang = 0x0017 // DW_LANG_Modula3
	DwarfLangHaskell      DwarfLang = 0x0018 // DW_LANG_Haskell
	DwarfLangCPlusPlus03  DwarfLang = 0x0019 // DW_LANG_C_plus_plus_03
	DwarfLangCPlusPlus11  DwarfLang = 0x001A // DW_LANG_C_plus_plus_11
	DwarfLangOCaml        DwarfLang = 0x001B // DW_LANG_OCaml
	DwarfLangRust         DwarfLang = 0x001C // DW_LANG_Rust
	DwarfLangC11          DwarfLang = 0x001D // DW_LANG_C11
	DwarfLangSwift        DwarfLang = 0x001E // DW_LANG_Swift
	DwarfLangJulia        DwarfLang = 0x001F // DW_LANG_Julia
	DwarfLangDylan        DwarfLang = 0x0020 // DW_LANG_Dylan
	DwarfLangCPlusPlus14  DwarfLang = 0x0021 // DW_LANG_C_plus_plus_14
	DwarfLangFortran03    DwarfLang = 0x0022 // DW_LANG_Fortran03
	DwarfLangFortran08    DwarfLang = 0x0023 // DW_LANG_Fortran08
	DwarfLangRenderScript DwarfLang = 0x0024 // DW_LANG_RenderScript
	DwarfLangBLISS        DwarfLang = 0x0025 // DW_LANG_BLISS
	// Vendor extensions.
	DwarfLangMipsAssembler      DwarfLang = 0x8001 // DW_LANG_Mips_Assembler
	DwarfLangGoogleRenderScript DwarfLang = 0x8E57 // DW_LANG_GOOGLE_RenderScript
	DwarfLangBorlandDelphi      DwarfLang = 0xB000 // DW_LANG_BORLAND_Delphi
)

//go:generate stringer -linecomment -type DwarfMacinfo

// DwarfMacinfo is a macinfo type encoding.
type DwarfMacinfo int64

// Macinfo type encodings.
//
// From llvm/BinaryFormat/Dwarf.h
const (
	DwarfMacinfoDefine    DwarfMacinfo = 0x01 // DW_MACINFO_define
	DwarfMacinfoUndef     DwarfMacinfo = 0x02 // DW_MACINFO_undef
	DwarfMacinfoStartFile DwarfMacinfo = 0x03 // DW_MACINFO_start_file
	DwarfMacinfoEndFile   DwarfMacinfo = 0x04 // DW_MACINFO_end_file
	DwarfMacinfoVendorExt DwarfMacinfo = 0xFF // DW_MACINFO_vendor_ext
)

//go:generate stringer -linecomment -type DwarfOp

// DwarfOp is a DWARF expression operator.
type DwarfOp int64

// DWARF expression operators.
//
// From include/llvm/BinaryFormat/Dwarf.def
const (
	// DWARF v2.
	DwarfOpAddr       DwarfOp = 0x03 // DW_OP_addr
	DwarfOpDeref      DwarfOp = 0x06 // DW_OP_deref
	DwarfOpConst1u    DwarfOp = 0x08 // DW_OP_const1u
	DwarfOpConst1s    DwarfOp = 0x09 // DW_OP_const1s
	DwarfOpConst2u    DwarfOp = 0x0A // DW_OP_const2u
	DwarfOpConst2s    DwarfOp = 0x0B // DW_OP_const2s
	DwarfOpConst4u    DwarfOp = 0x0C // DW_OP_const4u
	DwarfOpConst4s    DwarfOp = 0x0D // DW_OP_const4s
	DwarfOpConst8u    DwarfOp = 0x0E // DW_OP_const8u
	DwarfOpConst8s    DwarfOp = 0x0F // DW_OP_const8s
	DwarfOpConstu     DwarfOp = 0x10 // DW_OP_constu
	DwarfOpConsts     DwarfOp = 0x11 // DW_OP_consts
	DwarfOpDup        DwarfOp = 0x12 // DW_OP_dup
	DwarfOpDrop       DwarfOp = 0x13 // DW_OP_drop
	DwarfOpOver       DwarfOp = 0x14 // DW_OP_over
	DwarfOpPick       DwarfOp = 0x15 // DW_OP_pick
	DwarfOpSwap       DwarfOp = 0x16 // DW_OP_swap
	DwarfOpRot        DwarfOp = 0x17 // DW_OP_rot
	DwarfOpXderef     DwarfOp = 0x18 // DW_OP_xderef
	DwarfOpAbs        DwarfOp = 0x19 // DW_OP_abs
	DwarfOpAnd        DwarfOp = 0x1A // DW_OP_and
	DwarfOpDiv        DwarfOp = 0x1B // DW_OP_div
	DwarfOpMinus      DwarfOp = 0x1C // DW_OP_minus
	DwarfOpMod        DwarfOp = 0x1D // DW_OP_mod
	DwarfOpMul        DwarfOp = 0x1E // DW_OP_mul
	DwarfOpNeg        DwarfOp = 0x1F // DW_OP_neg
	DwarfOpNot        DwarfOp = 0x20 // DW_OP_not
	DwarfOpOr         DwarfOp = 0x21 // DW_OP_or
	DwarfOpPlus       DwarfOp = 0x22 // DW_OP_plus
	DwarfOpPlusUconst DwarfOp = 0x23 // DW_OP_plus_uconst
	DwarfOpShl        DwarfOp = 0x24 // DW_OP_shl
	DwarfOpShr        DwarfOp = 0x25 // DW_OP_shr
	DwarfOpShra       DwarfOp = 0x26 // DW_OP_shra
	DwarfOpXor        DwarfOp = 0x27 // DW_OP_xor
	DwarfOpBra        DwarfOp = 0x28 // DW_OP_bra
	DwarfOpEq         DwarfOp = 0x29 // DW_OP_eq
	DwarfOpGe         DwarfOp = 0x2A // DW_OP_ge
	DwarfOpGt         DwarfOp = 0x2B // DW_OP_gt
	DwarfOpLe         DwarfOp = 0x2C // DW_OP_le
	DwarfOpLt         DwarfOp = 0x2D // DW_OP_lt
	DwarfOpNe         DwarfOp = 0x2E // DW_OP_ne
	DwarfOpSkip       DwarfOp = 0x2F // DW_OP_skip
	DwarfOpLit0       DwarfOp = 0x30 // DW_OP_lit0
	DwarfOpLit1       DwarfOp = 0x31 // DW_OP_lit1
	DwarfOpLit2       DwarfOp = 0x32 // DW_OP_lit2
	DwarfOpLit3       DwarfOp = 0x33 // DW_OP_lit3
	DwarfOpLit4       DwarfOp = 0x34 // DW_OP_lit4
	DwarfOpLit5       DwarfOp = 0x35 // DW_OP_lit5
	DwarfOpLit6       DwarfOp = 0x36 // DW_OP_lit6
	DwarfOpLit7       DwarfOp = 0x37 // DW_OP_lit7
	DwarfOpLit8       DwarfOp = 0x38 // DW_OP_lit8
	DwarfOpLit9       DwarfOp = 0x39 // DW_OP_lit9
	DwarfOpLit10      DwarfOp = 0x3A // DW_OP_lit10
	DwarfOpLit11      DwarfOp = 0x3B // DW_OP_lit11
	DwarfOpLit12      DwarfOp = 0x3C // DW_OP_lit12
	DwarfOpLit13      DwarfOp = 0x3D // DW_OP_lit13
	DwarfOpLit14      DwarfOp = 0x3E // DW_OP_lit14
	DwarfOpLit15      DwarfOp = 0x3F // DW_OP_lit15
	DwarfOpLit16      DwarfOp = 0x40 // DW_OP_lit16
	DwarfOpLit17      DwarfOp = 0x41 // DW_OP_lit17
	DwarfOpLit18      DwarfOp = 0x42 // DW_OP_lit18
	DwarfOpLit19      DwarfOp = 0x43 // DW_OP_lit19
	DwarfOpLit20      DwarfOp = 0x44 // DW_OP_lit20
	DwarfOpLit21      DwarfOp = 0x45 // DW_OP_lit21
	DwarfOpLit22      DwarfOp = 0x46 // DW_OP_lit22
	DwarfOpLit23      DwarfOp = 0x47 // DW_OP_lit23
	DwarfOpLit24      DwarfOp = 0x48 // DW_OP_lit24
	DwarfOpLit25      DwarfOp = 0x49 // DW_OP_lit25
	DwarfOpLit26      DwarfOp = 0x4A // DW_OP_lit26
	DwarfOpLit27      DwarfOp = 0x4B // DW_OP_lit27
	DwarfOpLit28      DwarfOp = 0x4C // DW_OP_lit28
	DwarfOpLit29      DwarfOp = 0x4D // DW_OP_lit29
	DwarfOpLit30      DwarfOp = 0x4E // DW_OP_lit30
	DwarfOpLit31      DwarfOp = 0x4F // DW_OP_lit31
	DwarfOpReg0       DwarfOp = 0x50 // DW_OP_reg0
	DwarfOpReg1       DwarfOp = 0x51 // DW_OP_reg1
	DwarfOpReg2       DwarfOp = 0x52 // DW_OP_reg2
	DwarfOpReg3       DwarfOp = 0x53 // DW_OP_reg3
	DwarfOpReg4       DwarfOp = 0x54 // DW_OP_reg4
	DwarfOpReg5       DwarfOp = 0x55 // DW_OP_reg5
	DwarfOpReg6       DwarfOp = 0x56 // DW_OP_reg6
	DwarfOpReg7       DwarfOp = 0x57 // DW_OP_reg7
	DwarfOpReg8       DwarfOp = 0x58 // DW_OP_reg8
	DwarfOpReg9       DwarfOp = 0x59 // DW_OP_reg9
	DwarfOpReg10      DwarfOp = 0x5A // DW_OP_reg10
	DwarfOpReg11      DwarfOp = 0x5B // DW_OP_reg11
	DwarfOpReg12      DwarfOp = 0x5C // DW_OP_reg12
	DwarfOpReg13      DwarfOp = 0x5D // DW_OP_reg13
	DwarfOpReg14      DwarfOp = 0x5E // DW_OP_reg14
	DwarfOpReg15      DwarfOp = 0x5F // DW_OP_reg15
	DwarfOpReg16      DwarfOp = 0x60 // DW_OP_reg16
	DwarfOpReg17      DwarfOp = 0x61 // DW_OP_reg17
	DwarfOpReg18      DwarfOp = 0x62 // DW_OP_reg18
	DwarfOpReg19      DwarfOp = 0x63 // DW_OP_reg19
	DwarfOpReg20      DwarfOp = 0x64 // DW_OP_reg20
	DwarfOpReg21      DwarfOp = 0x65 // DW_OP_reg21
	DwarfOpReg22      DwarfOp = 0x66 // DW_OP_reg22
	DwarfOpReg23      DwarfOp = 0x67 // DW_OP_reg23
	DwarfOpReg24      DwarfOp = 0x68 // DW_OP_reg24
	DwarfOpReg25      DwarfOp = 0x69 // DW_OP_reg25
	DwarfOpReg26      DwarfOp = 0x6A // DW_OP_reg26
	DwarfOpReg27      DwarfOp = 0x6B // DW_OP_reg27
	DwarfOpReg28      DwarfOp = 0x6C // DW_OP_reg28
	DwarfOpReg29      DwarfOp = 0x6D // DW_OP_reg29
	DwarfOpReg30      DwarfOp = 0x6E // DW_OP_reg30
	DwarfOpReg31      DwarfOp = 0x6F // DW_OP_reg31
	DwarfOpBreg0      DwarfOp = 0x70 // DW_OP_breg0
	DwarfOpBreg1      DwarfOp = 0x71 // DW_OP_breg1
	DwarfOpBreg2      DwarfOp = 0x72 // DW_OP_breg2
	DwarfOpBreg3      DwarfOp = 0x73 // DW_OP_breg3
	DwarfOpBreg4      DwarfOp = 0x74 // DW_OP_breg4
	DwarfOpBreg5      DwarfOp = 0x75 // DW_OP_breg5
	DwarfOpBreg6      DwarfOp = 0x76 // DW_OP_breg6
	DwarfOpBreg7      DwarfOp = 0x77 // DW_OP_breg7
	DwarfOpBreg8      DwarfOp = 0x78 // DW_OP_breg8
	DwarfOpBreg9      DwarfOp = 0x79 // DW_OP_breg9
	DwarfOpBreg10     DwarfOp = 0x7A // DW_OP_breg10
	DwarfOpBreg11     DwarfOp = 0x7B // DW_OP_breg11
	DwarfOpBreg12     DwarfOp = 0x7C // DW_OP_breg12
	DwarfOpBreg13     DwarfOp = 0x7D // DW_OP_breg13
	DwarfOpBreg14     DwarfOp = 0x7E // DW_OP_breg14
	DwarfOpBreg15     DwarfOp = 0x7F // DW_OP_breg15
	DwarfOpBreg16     DwarfOp = 0x80 // DW_OP_breg16
	DwarfOpBreg17     DwarfOp = 0x81 // DW_OP_breg17
	DwarfOpBreg18     DwarfOp = 0x82 // DW_OP_breg18
	DwarfOpBreg19     DwarfOp = 0x83 // DW_OP_breg19
	DwarfOpBreg20     DwarfOp = 0x84 // DW_OP_breg20
	DwarfOpBreg21     DwarfOp = 0x85 // DW_OP_breg21
	DwarfOpBreg22     DwarfOp = 0x86 // DW_OP_breg22
	DwarfOpBreg23     DwarfOp = 0x87 // DW_OP_breg23
	DwarfOpBreg24     DwarfOp = 0x88 // DW_OP_breg24
	DwarfOpBreg25     DwarfOp = 0x89 // DW_OP_breg25
	DwarfOpBreg26     DwarfOp = 0x8A // DW_OP_breg26
	DwarfOpBreg27     DwarfOp = 0x8B // DW_OP_breg27
	DwarfOpBreg28     DwarfOp = 0x8C // DW_OP_breg28
	DwarfOpBreg29     DwarfOp = 0x8D // DW_OP_breg29
	DwarfOpBreg30     DwarfOp = 0x8E // DW_OP_breg30
	DwarfOpBreg31     DwarfOp = 0x8F // DW_OP_breg31
	DwarfOpRegx       DwarfOp = 0x90 // DW_OP_regx
	DwarfOpFbreg      DwarfOp = 0x91 // DW_OP_fbreg
	DwarfOpBregx      DwarfOp = 0x92 // DW_OP_bregx
	DwarfOpPiece      DwarfOp = 0x93 // DW_OP_piece
	DwarfOpDerefSize  DwarfOp = 0x94 // DW_OP_deref_size
	DwarfOpXderefSize DwarfOp = 0x95 // DW_OP_xderef_size
	DwarfOpNop        DwarfOp = 0x96 // DW_OP_nop
	// DWARF v3.
	DwarfOpPushObjectAddress DwarfOp = 0x97 // DW_OP_push_object_address
	DwarfOpCall2             DwarfOp = 0x98 // DW_OP_call2
	DwarfOpCall4             DwarfOp = 0x99 // DW_OP_call4
	DwarfOpCallRef           DwarfOp = 0x9A // DW_OP_call_ref
	DwarfOpFormTLSAddress    DwarfOp = 0x9B // DW_OP_form_tls_address
	DwarfOpCallFrameCFA      DwarfOp = 0x9C // DW_OP_call_frame_cfa
	DwarfOpBitPiece          DwarfOp = 0x9D // DW_OP_bit_piece
	// DWARF v4.
	DwarfOpImplicitValue DwarfOp = 0x9E // DW_OP_implicit_value
	DwarfOpStackValue    DwarfOp = 0x9F // DW_OP_stack_value
	// DWARF v5.
	DwarfOpImplicitPointer DwarfOp = 0xA0 // DW_OP_implicit_pointer
	DwarfOpAddrx           DwarfOp = 0xA1 // DW_OP_addrx
	DwarfOpConstx          DwarfOp = 0xA2 // DW_OP_constx
	DwarfOpEntryValue      DwarfOp = 0xA3 // DW_OP_entry_value
	DwarfOpConstType       DwarfOp = 0xA4 // DW_OP_const_type
	DwarfOpRegvalType      DwarfOp = 0xA5 // DW_OP_regval_type
	DwarfOpDerefType       DwarfOp = 0xA6 // DW_OP_deref_type
	DwarfOpXderefType      DwarfOp = 0xA7 // DW_OP_xderef_type
	DwarfOpConvert         DwarfOp = 0xA8 // DW_OP_convert
	DwarfOpReinterpret     DwarfOp = 0xA9 // DW_OP_reinterpret
	// Vendor extensions.
	DwarfOpGNUPushTLSAddress DwarfOp = 0xE0 // DW_OP_GNU_push_tls_address
	DwarfOpGNUEntryValue     DwarfOp = 0xF3 // DW_OP_GNU_entry_value
	DwarfOpGNUAddrIndex      DwarfOp = 0xFB // DW_OP_GNU_addr_index
	DwarfOpGNUConstIndex     DwarfOp = 0xFC // DW_OP_GNU_const_index
	// Only used in LLVM metadata.
	DwarfOpLLVMFragment  DwarfOp = 0x1000 // DW_OP_LLVM_fragment
	DwarfOpLLVMConvert   DwarfOp = 0x1001 // DW_OP_LLVM_convert
	DwarfOpLLVMTagOffset DwarfOp = 0x1002 // DW_OP_LLVM_tag_offset
)

//go:generate stringer -linecomment -type DwarfTag

// DwarfTag is a DWARF tag.
type DwarfTag int64

// DWARF tags.
//
// From include/llvm/BinaryFormat/Dwarf.def
const (
	// DWARF v2.
	DwarfTagNull                   DwarfTag = 0x0000 // DW_TAG_null
	DwarfTagArrayType              DwarfTag = 0x0001 // DW_TAG_array_type
	DwarfTagClassType              DwarfTag = 0x0002 // DW_TAG_class_type
	DwarfTagEntryPoint             DwarfTag = 0x0003 // DW_TAG_entry_point
	DwarfTagEnumerationType        DwarfTag = 0x0004 // DW_TAG_enumeration_type
	DwarfTagFormalParameter        DwarfTag = 0x0005 // DW_TAG_formal_parameter
	DwarfTagImportedDeclaration    DwarfTag = 0x0008 // DW_TAG_imported_declaration
	DwarfTagLabel                  DwarfTag = 0x000A // DW_TAG_label
	DwarfTagLexicalBlock           DwarfTag = 0x000B // DW_TAG_lexical_block
	DwarfTagMember                 DwarfTag = 0x000D // DW_TAG_member
	DwarfTagPointerType            DwarfTag = 0x000F // DW_TAG_pointer_type
	DwarfTagReferenceType          DwarfTag = 0x0010 // DW_TAG_reference_type
	DwarfTagCompileUnit            DwarfTag = 0x0011 // DW_TAG_compile_unit
	DwarfTagStringType             DwarfTag = 0x0012 // DW_TAG_string_type
	DwarfTagStructureType          DwarfTag = 0x0013 // DW_TAG_structure_type
	DwarfTagSubroutineType         DwarfTag = 0x0015 // DW_TAG_subroutine_type
	DwarfTagTypedef                DwarfTag = 0x0016 // DW_TAG_typedef
	DwarfTagUnionType              DwarfTag = 0x0017 // DW_TAG_union_type
	DwarfTagUnspecifiedParameters  DwarfTag = 0x0018 // DW_TAG_unspecified_parameters
	DwarfTagVariant                DwarfTag = 0x0019 // DW_TAG_variant
	DwarfTagCommonBlock            DwarfTag = 0x001A // DW_TAG_common_block
	DwarfTagCommonInclusion        DwarfTag = 0x001B // DW_TAG_common_inclusion
	DwarfTagInheritance            DwarfTag = 0x001C // DW_TAG_inheritance
	DwarfTagInlinedSubroutine      DwarfTag = 0x001D // DW_TAG_inlined_subroutine
	DwarfTagModule                 DwarfTag = 0x001E // DW_TAG_module
	DwarfTagPtrToMemberType        DwarfTag = 0x001F // DW_TAG_ptr_to_member_type
	DwarfTagSetType                DwarfTag = 0x0020 // DW_TAG_set_type
	DwarfTagSubrangeType           DwarfTag = 0x0021 // DW_TAG_subrange_type
	DwarfTagWithStmt               DwarfTag = 0x0022 // DW_TAG_with_stmt
	DwarfTagAccessDeclaration      DwarfTag = 0x0023 // DW_TAG_access_declaration
	DwarfTagBaseType               DwarfTag = 0x0024 // DW_TAG_base_type
	DwarfTagCatchBlock             DwarfTag = 0x0025 // DW_TAG_catch_block
	DwarfTagConstType              DwarfTag = 0x0026 // DW_TAG_const_type
	DwarfTagConstant               DwarfTag = 0x0027 // DW_TAG_constant
	DwarfTagEnumerator             DwarfTag = 0x0028 // DW_TAG_enumerator
	DwarfTagFileType               DwarfTag = 0x0029 // DW_TAG_file_type
	DwarfTagFriend                 DwarfTag = 0x002A // DW_TAG_friend
	DwarfTagNamelist               DwarfTag = 0x002B // DW_TAG_namelist
	DwarfTagNamelistItem           DwarfTag = 0x002C // DW_TAG_namelist_item
	DwarfTagPackedType             DwarfTag = 0x002D // DW_TAG_packed_type
	DwarfTagSubprogram             DwarfTag = 0x002E // DW_TAG_subprogram
	DwarfTagTemplateTypeParameter  DwarfTag = 0x002F // DW_TAG_template_type_parameter
	DwarfTagTemplateValueParameter DwarfTag = 0x0030 // DW_TAG_template_value_parameter
	DwarfTagThrownType             DwarfTag = 0x0031 // DW_TAG_thrown_type
	DwarfTagTryBlock               DwarfTag = 0x0032 // DW_TAG_try_block
	DwarfTagVariantPart            DwarfTag = 0x0033 // DW_TAG_variant_part
	DwarfTagVariable               DwarfTag = 0x0034 // DW_TAG_variable
	DwarfTagVolatileType           DwarfTag = 0x0035 // DW_TAG_volatile_type
	// DWARF v3.
	DwarfTagDwarfProcedure  DwarfTag = 0x0036 // DW_TAG_dwarf_procedure
	DwarfTagRestrictType    DwarfTag = 0x0037 // DW_TAG_restrict_type
	DwarfTagInterfaceType   DwarfTag = 0x0038 // DW_TAG_interface_type
	DwarfTagNamespace       DwarfTag = 0x0039 // DW_TAG_namespace
	DwarfTagImportedModule  DwarfTag = 0x003A // DW_TAG_imported_module
	DwarfTagUnspecifiedType DwarfTag = 0x003B // DW_TAG_unspecified_type
	DwarfTagPartialUnit     DwarfTag = 0x003C // DW_TAG_partial_unit
	DwarfTagImportedUnit    DwarfTag = 0x003D // DW_TAG_imported_unit
	DwarfTagCondition       DwarfTag = 0x003F // DW_TAG_condition
	DwarfTagSharedType      DwarfTag = 0x0040 // DW_TAG_shared_type
	// DWARF v4.
	DwarfTagTypeUnit            DwarfTag = 0x0041 // DW_TAG_type_unit
	DwarfTagRvalueReferenceType DwarfTag = 0x0042 // DW_TAG_rvalue_reference_type
	DwarfTagTemplateAlias       DwarfTag = 0x0043 // DW_TAG_template_alias
	// DWARF v5.
	DwarfTagCoarrayType       DwarfTag = 0x0044 // DW_TAG_coarray_type
	DwarfTagGenericSubrange   DwarfTag = 0x0045 // DW_TAG_generic_subrange
	DwarfTagDynamicType       DwarfTag = 0x0046 // DW_TAG_dynamic_type
	DwarfTagAtomicType        DwarfTag = 0x0047 // DW_TAG_atomic_type
	DwarfTagCallSite          DwarfTag = 0x0048 // DW_TAG_call_site
	DwarfTagCallSiteParameter DwarfTag = 0x0049 // DW_TAG_call_site_parameter
	DwarfTagSkeletonUnit      DwarfTag = 0x004A // DW_TAG_skeleton_unit
	DwarfTagImmutableType     DwarfTag = 0x004B // DW_TAG_immutable_type
	// Vendor extensions.
	DwarfTagMIPSLoop                  DwarfTag = 0x4081 // DW_TAG_MIPS_loop
	DwarfTagFormatLabel               DwarfTag = 0x4101 // DW_TAG_format_label
	DwarfTagFunctionTemplate          DwarfTag = 0x4102 // DW_TAG_function_template
	DwarfTagClassTemplate             DwarfTag = 0x4103 // DW_TAG_class_template
	DwarfTagGNUTemplateTemplateParam  DwarfTag = 0x4106 // DW_TAG_GNU_template_template_param
	DwarfTagGNUTemplateParameterPack  DwarfTag = 0x4107 // DW_TAG_GNU_template_parameter_pack
	DwarfTagGNUFormalParameterPack    DwarfTag = 0x4108 // DW_TAG_GNU_formal_parameter_pack
	DwarfTagGNUCallSite               DwarfTag = 0x4109 // DW_TAG_GNU_call_site
	DwarfTagGNUCallSiteParameter      DwarfTag = 0x410A // DW_TAG_GNU_call_site_parameter
	DwarfTagAPPLEProperty             DwarfTag = 0x4200 // DW_TAG_APPLE_property
	DwarfTagBORLANDProperty           DwarfTag = 0xB000 // DW_TAG_BORLAND_property
	DwarfTagBORLANDDelphiString       DwarfTag = 0xB001 // DW_TAG_BORLAND_Delphi_string
	DwarfTagBORLANDDelphiDynamicArray DwarfTag = 0xB002 // DW_TAG_BORLAND_Delphi_dynamic_array
	DwarfTagBORLANDDelphiSet          DwarfTag = 0xB003 // DW_TAG_BORLAND_Delphi_set
	DwarfTagBORLANDDelphiVariant      DwarfTag = 0xB004 // DW_TAG_BORLAND_Delphi_variant
)

//go:generate stringer -linecomment -type DwarfVirtuality

// DwarfVirtuality is a DWARF virtuality code.
type DwarfVirtuality int64

// DWARF virtuality codes.
const (
	DwarfVirtualityNone        DwarfVirtuality = 0x00 // DW_VIRTUALITY_none
	DwarfVirtualityVirtual     DwarfVirtuality = 0x01 // DW_VIRTUALITY_virtual
	DwarfVirtualityPureVirtual DwarfVirtuality = 0x02 // DW_VIRTUALITY_pure_virtual
)

//go:generate stringer -linecomment -type EmissionKind

// EmissionKind specifies the debug emission kind.
type EmissionKind int64

// ref: include/llvm/IR/DebugInfoMetadata.h (LLVM 9.0).

// Debug emission kinds.
const (
	EmissionKindNoDebug             EmissionKind = 0 // NoDebug
	EmissionKindFullDebug           EmissionKind = 1 // FullDebug
	EmissionKindLineTablesOnly      EmissionKind = 2 // LineTablesOnly
	EmissionKindDebugDirectivesOnly EmissionKind = 3 // DebugDirectivesOnly
)

//go:generate stringer -linecomment -type FastMathFlag

// FastMathFlag is a fast-math flag.
type FastMathFlag uint8

// Fast-math flags.
const (
	FastMathFlagAFn      FastMathFlag = iota // afn
	FastMathFlagARcp                         // arcp
	FastMathFlagContract                     // contract
	FastMathFlagFast                         // fast
	FastMathFlagNInf                         // ninf
	FastMathFlagNNaN                         // nnan
	FastMathFlagNSZ                          // nsz
	FastMathFlagReassoc                      // reassoc
)

//go:generate stringer -linecomment -type FPred

// FPred is a floating-point comparison predicate.
type FPred uint8

// Floating-point predicates.
const (
	FPredFalse FPred = iota // false
	FPredOEQ                // oeq
	FPredOGE                // oge
	FPredOGT                // ogt
	FPredOLE                // ole
	FPredOLT                // olt
	FPredONE                // one
	FPredORD                // ord
	FPredTrue               // true
	FPredUEQ                // ueq
	FPredUGE                // uge
	FPredUGT                // ugt
	FPredULE                // ule
	FPredULT                // ult
	FPredUNE                // une
	FPredUNO                // uno
)

//go:generate stringer -linecomment -type FuncAttr

// FuncAttr is a function attribute.
type FuncAttr uint8

// Function attributes.
const (
	FuncAttrAlwaysInline                FuncAttr = iota // alwaysinline
	FuncAttrArgMemOnly                                  // argmemonly
	FuncAttrBuiltin                                     // builtin
	FuncAttrCold                                        // cold
	FuncAttrConvergent                                  // convergent
	FuncAttrInaccessibleMemOnly                         // inaccessiblememonly
	FuncAttrInaccessibleMemOrArgMemOnly                 // inaccessiblemem_or_argmemonly
	FuncAttrInlineHint                                  // inlinehint
	FuncAttrJumpTable                                   // jumptable
	FuncAttrMinSize                                     // minsize
	FuncAttrNaked                                       // naked
	FuncAttrNoBuiltin                                   // nobuiltin
	FuncAttrNoCFCheck                                   // nocf_check
	FuncAttrNoDuplicate                                 // noduplicate
	FuncAttrNoFree                                      // nofree
	FuncAttrNoImplicitFloat                             // noimplicitfloat
	FuncAttrNoInline                                    // noinline
	FuncAttrNoMerge                                     // nomerge
	FuncAttrNonLazyBind                                 // nonlazybind
	FuncAttrNoRecurse                                   // norecurse
	FuncAttrNoRedZone                                   // noredzone
	FuncAttrNoReturn                                    // noreturn
	FuncAttrNoSync                                      // nosync
	FuncAttrNoUnwind                                    // nounwind
	FuncAttrNullPointerIsValid                          // null_pointer_is_valid
	FuncAttrOptForFuzzing                               // optforfuzzing
	FuncAttrOptNone                                     // optnone
	FuncAttrOptSize                                     // optsize
	FuncAttrReadNone                                    // readnone
	FuncAttrReadOnly                                    // readonly
	FuncAttrReturnsTwice                                // returns_twice
	FuncAttrSafeStack                                   // safestack
	FuncAttrSanitizeAddress                             // sanitize_address
	FuncAttrSanitizeHWAddress                           // sanitize_hwaddress
	FuncAttrSanitizeMemory                              // sanitize_memory
	FuncAttrSanitizeMemTag                              // sanitize_memtag
	FuncAttrSanitizeThread                              // sanitize_thread
	FuncAttrShadowCallStack                             // shadowcallstack
	FuncAttrSpeculatable                                // speculatable
	FuncAttrSpeculativeLoadHardening                    // speculative_load_hardening
	FuncAttrSSP                                         // ssp
	FuncAttrSSPReq                                      // sspreq
	FuncAttrSSPStrong                                   // sspstrong
	FuncAttrStrictFP                                    // strictfp
	FuncAttrUwtable                                     // uwtable
	FuncAttrWillReturn                                  // willreturn
	FuncAttrWriteOnly                                   // writeonly
)

//go:generate stringer -linecomment -type IPred

// IPred is an integer comparison predicate.
type IPred uint8

// Integer predicates.
const (
	IPredEQ  IPred = iota // eq
	IPredNE               // ne
	IPredSGE              // sge
	IPredSGT              // sgt
	IPredSLE              // sle
	IPredSLT              // slt
	IPredUGE              // uge
	IPredUGT              // ugt
	IPredULE              // ule
	IPredULT              // ult
)

//go:generate stringer -linecomment -type Linkage

// Linkage specifies the linkage of a global identifier.
type Linkage uint8

// Linkage kinds.
const (
	LinkageNone                Linkage = iota // none
	LinkageAppending                          // appending
	LinkageAvailableExternally                // available_externally
	LinkageCommon                             // common
	LinkageInternal                           // internal
	LinkageLinkOnce                           // linkonce
	LinkageLinkOnceODR                        // linkonce_odr
	LinkagePrivate                            // private
	LinkageWeak                               // weak
	LinkageWeakODR                            // weak_odr
	// External linkage.
	LinkageExternal   // external
	LinkageExternWeak // extern_weak
)

//go:generate stringer -linecomment -type NameTableKind

// NameTableKind is a name table specifier.
type NameTableKind uint8

// Name table kinds.
//
// From include/llvm/IR/DebugInfoMetadata.h
const (
	NameTableKindDefault NameTableKind = 0 // Default
	NameTableKindGNU     NameTableKind = 1 // GNU
	NameTableKindNone    NameTableKind = 2 // None
)

//go:generate stringer -linecomment -type OverflowFlag

// OverflowFlag is an integer overflow flag.
type OverflowFlag uint8

// Overflow flags.
const (
	OverflowFlagNSW OverflowFlag = iota // nsw
	OverflowFlagNUW                     // nuw
)

//go:generate stringer -linecomment -type ParamAttr

// ParamAttr is a parameter attribute.
type ParamAttr uint8

// Parameter attributes.
const (
	ParamAttrImmArg             ParamAttr = iota // immarg
	ParamAttrInAlloca                            // inalloca
	ParamAttrInReg                               // inreg
	ParamAttrNest                                // nest
	ParamAttrNoAlias                             // noalias
	ParamAttrNoCapture                           // nocapture
	ParamAttrNoFree                              // nofree
	ParamAttrNoMerge                             // nomerge
	ParamAttrNonNull                             // nonnull
	ParamAttrNoUndef                             // noundef
	ParamAttrNullPointerIsValid                  // null_pointer_is_valid
	ParamAttrReadNone                            // readnone
	ParamAttrReadOnly                            // readonly
	ParamAttrReturned                            // returned
	ParamAttrSignExt                             // signext
	ParamAttrSwiftError                          // swifterror
	ParamAttrSwiftSelf                           // swiftself
	ParamAttrWriteOnly                           // writeonly
	ParamAttrZeroExt                             // zeroext
)

//go:generate stringer -linecomment -type Preemption

// Preemption specifies the preemtion of a global identifier.
type Preemption uint8

// Preemption kinds.
const (
	PreemptionNone               Preemption = iota // none
	PreemptionDSOLocal                             // dso_local
	PreemptionDSOLocalEquivalent                   // dso_local_equivalent
	PreemptionDSOPreemptable                       // dso_preemptable
)

//go:generate stringer -linecomment -type ReturnAttr

// ReturnAttr is a return argument attribute.
type ReturnAttr uint8

// Return argument attributes.
const (
	ReturnAttrInReg              ReturnAttr = iota // inreg
	ReturnAttrNoAlias                              // noalias
	ReturnAttrNoMerge                              // nomerge
	ReturnAttrNonNull                              // nonnull
	ReturnAttrNoUndef                              // noundef
	ReturnAttrNullPointerIsValid                   // null_pointer_is_valid
	ReturnAttrSignExt                              // signext
	ReturnAttrZeroExt                              // zeroext
)

//go:generate stringer -linecomment -type SelectionKind

// SelectionKind is a Comdat selection kind.
type SelectionKind uint8

// Comdat selection kinds.
const (
	SelectionKindAny          SelectionKind = iota // any
	SelectionKindExactMatch                        // exactmatch
	SelectionKindLargest                           // largest
	SelectionKindNoDuplicates                      // noduplicates
	SelectionKindSameSize                          // samesize
)

//go:generate stringer -linecomment -type Tail

// Tail is a tail call attribute.
type Tail uint8

// Tail call attributes.
const (
	TailNone     Tail = iota // none
	TailMustTail             // musttail
	TailNoTail               // notail
	TailTail                 // tail
)

//go:generate stringer -linecomment -type TLSModel

// TLSModel is a thread local storage model.
type TLSModel uint8

// Thread local storage models.
const (
	TLSModelNone TLSModel = iota // none
	// If no explicit model is given, the "general dynamic" model is used.
	TLSModelGeneric      // generic
	TLSModelInitialExec  // initialexec
	TLSModelLocalDynamic // localdynamic
	TLSModelLocalExec    // localexec
)

//go:generate stringer -linecomment -type UnnamedAddr

// UnnamedAddr specifies whether the address is significant.
type UnnamedAddr uint8

// Unnamed address specifiers.
const (
	UnnamedAddrNone             UnnamedAddr = iota // none
	UnnamedAddrLocalUnnamedAddr                    // local_unnamed_addr
	UnnamedAddrUnnamedAddr                         // unnamed_addr
)

//go:generate stringer -linecomment -type Visibility

// Visibility specifies the visibility of a global identifier.
type Visibility uint8

// Visibility kinds.
const (
	VisibilityNone      Visibility = iota // none
	VisibilityDefault                     // default
	VisibilityHidden                      // hidden
	VisibilityProtected                   // protected
)
