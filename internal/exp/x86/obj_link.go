package x86

// An As denotes an assembler opcode.
// There are some portable opcodes, declared here in package obj,
// that are common to all architectures.
// However, the majority of opcodes are arch-specific
// and are declared in their respective architecture's subpackage.
type As int16

// These are the portable opcodes.
const (
	AXXX = 0 + iota
	ACALL
	ACHECKNIL
	ADATA
	ADUFFCOPY
	ADUFFZERO
	AEND
	AFUNCDATA
	AGLOBL
	AJMP
	ANOP
	APCDATA
	ARET
	ATEXT
	ATYPE
	AUNDEF
	AUSEFIELD
	AVARDEF
	AVARKILL
	A_ARCHSPECIFIC
)

// Each architecture is allotted a distinct subspace of opcode values
// for declaring its arch-specific opcodes.
// Within this subspace, the first arch-specific opcode should be
// at offset A_ARCHSPECIFIC.
//
// Subspaces are aligned to a power of two so opcodes can be masked
// with AMask and used as compact array indices.
const (
	ABase386 = (1 + iota) << 11
	ABaseARM
	ABaseAMD64
	ABasePPC64
	ABaseARM64
	ABaseMIPS
	ABaseLoong64
	ABaseRISCV
	ABaseS390X
	ABaseWasm

	AllowedOpCodes = 1 << 11            // The number of opcodes available for any given architecture.
	AMask          = AllowedOpCodes - 1 // AND with this to use the opcode as an array index.
)

/*
	Each architecture defines a register space as a unique
	integer range.
	Here is the list of architectures and the base of their register spaces.
*/

const (
	// Because of masking operations in the encodings, each register
	// space should start at 0 modulo some power of 2.
	RBase386     = 1 * 1024
	RBaseAMD64   = 2 * 1024
	RBaseARM     = 3 * 1024
	RBasePPC64   = 4 * 1024  // range [4k, 8k)
	RBaseARM64   = 8 * 1024  // range [8k, 13k)
	RBaseMIPS    = 13 * 1024 // range [13k, 14k)
	RBaseS390X   = 14 * 1024 // range [14k, 15k)
	RBaseRISCV   = 15 * 1024 // range [15k, 16k)
	RBaseWasm    = 16 * 1024
	RBaseLOONG64 = 19 * 1024 // range [19K, 22k)
)
