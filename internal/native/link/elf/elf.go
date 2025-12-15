// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package elf

import "fmt"

// Initial magic number for ELF files.
const ELFMAG = "\177ELF"

// Version is found in Header.Ident[EI_VERSION] and Header.Version.
type Version byte

const (
	EV_NONE    Version = 0
	EV_CURRENT Version = 1
)

const (
	ELF64HDRSIZE = 64 // ELF64 文件头大小
	ELF32HDRSIZE = 52 // ELF32 文件头大小

	ELF64PHDRSIZE = 56 // ELF64 程序头大小
	ELF32PHDRSIZE = 32 // ELF32 程序头大小

	ELF64SHDRSIZE = 64 // ELF64 段头大小
	ELF32SHDRSIZE = 40 // ELF32 段头大小
)

// Indexes into the Header.Ident array.
const (
	EI_CLASS      = 4  // Class of machine.
	EI_DATA       = 5  // Data format.
	EI_VERSION    = 6  // ELF format version.
	EI_OSABI      = 7  // Operating system / ABI identification
	EI_ABIVERSION = 8  // ABI version
	EI_PAD        = 9  // Start of padding (per SVR4 ABI).
	EI_NIDENT     = 16 // Size of e_ident array.
)

// Class is found in Header.Ident[EI_CLASS] and Header.Class.
type Class byte

const (
	ELFCLASSNONE Class = 0 // Unknown class.
	ELFCLASS32   Class = 1 // 32-bit architecture.
	ELFCLASS64   Class = 2 // 64-bit architecture.
)

var classStrings = []string{
	ELFCLASSNONE: "ELFCLASSNONE",
	ELFCLASS32:   "ELFCLASS32",
	ELFCLASS64:   "ELFCLASS64",
}

func (i Class) String() string {
	var s string
	if int(i) < len(classStrings) {
		s = classStrings[int(i)]
	}
	if s == "" {
		s = fmt.Sprintf("elf.Class(%d)", int(i))
	}
	return s
}

// Data is found in Header.Ident[EI_DATA] and Header.Data.
type Data byte

const (
	ELFDATANONE Data = 0 // Unknown data format.
	ELFDATA2LSB Data = 1 // 2's complement little-endian.
	ELFDATA2MSB Data = 2 // 2's complement big-endian.
)

// OSABI is found in Header.Ident[EI_OSABI] and Header.OSABI.
type OSABI byte

const (
	ELFOSABI_NONE  OSABI = 0 // UNIX System V ABI
	ELFOSABI_LINUX OSABI = 3 // Linux
)

// Type is found in Header.Type.
type Type uint16

const (
	ET_NONE   Type = 0      // Unknown type.
	ET_REL    Type = 1      // Relocatable.
	ET_EXEC   Type = 2      // Executable.
	ET_DYN    Type = 3      // Shared object.
	ET_CORE   Type = 4      // Core file.
	ET_LOOS   Type = 0xfe00 // First operating system specific.
	ET_HIOS   Type = 0xfeff // Last operating system-specific.
	ET_LOPROC Type = 0xff00 // First processor-specific.
	ET_HIPROC Type = 0xffff // Last processor-specific.
)

// Machine is found in Header.Machine.
type Machine uint16

const (
	EM_NONE      Machine = 0   // Unknown machine.
	EM_X86_64    Machine = 62  // Advanced Micro Devices x86-64
	EM_AARCH64   Machine = 183 // ARM 64-bit Architecture (AArch64)
	EM_RISCV     Machine = 243 // RISC-V
	EM_BPF       Machine = 247 // Linux BPF – in-kernel virtual machine
	EM_LOONGARCH Machine = 258 // LoongArch
)

var machineStrings = []string{
	EM_NONE:      "EM_NONE",
	EM_X86_64:    "EM_X86_64",
	EM_AARCH64:   "EM_AARCH64",
	EM_RISCV:     "EM_RISCV",
	EM_BPF:       "EM_BPF",
	EM_LOONGARCH: "EM_LOONGARCH",
}

func (i Machine) String() string {
	var s string
	if int(i) < len(machineStrings) {
		s = machineStrings[int(i)]
	}
	if s == "" {
		s = fmt.Sprintf("elf.Machine(%d)", int(i))
	}
	return s
}

// Special section indices.
type SectionIndex int

const (
	SHN_UNDEF     SectionIndex = 0      // Undefined, missing, irrelevant.
	SHN_LORESERVE SectionIndex = 0xff00 // First of reserved range.
	SHN_LOPROC    SectionIndex = 0xff00 // First processor-specific.
	SHN_HIPROC    SectionIndex = 0xff1f // Last processor-specific.
	SHN_LOOS      SectionIndex = 0xff20 // First operating system-specific.
	SHN_HIOS      SectionIndex = 0xff3f // Last operating system-specific.
	SHN_ABS       SectionIndex = 0xfff1 // Absolute values.
	SHN_COMMON    SectionIndex = 0xfff2 // Common data.
	SHN_XINDEX    SectionIndex = 0xffff // Escape; index stored elsewhere.
	SHN_HIRESERVE SectionIndex = 0xffff // Last of reserved range.
)

// Prog.Type
type ProgType uint32

const (
	PT_NULL    ProgType = 0 // Unused entry.
	PT_LOAD    ProgType = 1 // Loadable segment.
	PT_DYNAMIC ProgType = 2 // Dynamic linking information segment.
	PT_INTERP  ProgType = 3 // Pathname of interpreter.
	PT_NOTE    ProgType = 4 // Auxiliary information.
	PT_SHLIB   ProgType = 5 // Reserved (not used).
	PT_PHDR    ProgType = 6 // Location of program header itself.
	PT_TLS     ProgType = 7 // Thread local storage segment
)

// Prog.Flag
type ProgFlag uint32

const (
	PF_X        ProgFlag = 0x1        // Executable.
	PF_W        ProgFlag = 0x2        // Writable.
	PF_R        ProgFlag = 0x4        // Readable.
	PF_MASKOS   ProgFlag = 0x0ff00000 // Operating system-specific.
	PF_MASKPROC ProgFlag = 0xf0000000 // Processor-specific.
)

// ELF32 头
type ElfHeader32 struct {
	Ident     [EI_NIDENT]byte // File identification.
	Type      uint16          // File type.
	Machine   uint16          // Machine architecture.
	Version   uint32          // ELF format version.
	Entry     uint32          // Entry point.
	Phoff     uint32          // Program header file offset.
	Shoff     uint32          // Section header file offset.
	Flags     uint32          // Architecture-specific flags.
	Ehsize    uint16          // Size of ELF header in bytes.
	Phentsize uint16          // Size of program header entry.
	Phnum     uint16          // Number of program header entries.
	Shentsize uint16          // Size of section header entry.
	Shnum     uint16          // Number of section header entries.
	Shstrndx  uint16          // Section name strings section.
}

// ELF64 头
type ElfHeader64 struct {
	Ident     [EI_NIDENT]byte // File identification.
	Type      uint16          // File type.
	Machine   uint16          // Machine architecture.
	Version   uint32          // ELF format version.
	Entry     uint64          // Entry point.
	Phoff     uint64          // Program header file offset.
	Shoff     uint64          // Section header file offset.
	Flags     uint32          // Architecture-specific flags.
	Ehsize    uint16          // Size of ELF header in bytes.
	Phentsize uint16          // Size of program header entry.
	Phnum     uint16          // Number of program header entries.
	Shentsize uint16          // Size of section header entry.
	Shnum     uint16          // Number of section header entries.
	Shstrndx  uint16          // Section name strings section.
}

// ELF32 段数据头
type ElfSection32 struct {
	Name      uint32 // Section name (index into the section header string table).
	Type      uint32 // Section type.
	Flags     uint32 // Section flags.
	Addr      uint32 // Address in memory image.
	Off       uint32 // Offset in file.
	Size      uint32 // Size in bytes.
	Link      uint32 // Index of a related section.
	Info      uint32 // Depends on section type.
	Addralign uint32 // Alignment in bytes.
	Entsize   uint32 // Size of each entry in section.
}

// ELF64 段数据头
type ElfSection64 struct {
	Name      uint32 // Section name (index into the section header string table).
	Type      uint32 // Section type.
	Flags     uint64 // Section flags.
	Addr      uint64 // Address in memory image.
	Off       uint64 // Offset in file.
	Size      uint64 // Size in bytes.
	Link      uint32 // Index of a related section.
	Info      uint32 // Depends on section type.
	Addralign uint64 // Alignment in bytes.
	Entsize   uint64 // Size of each entry in section.
}

// ELF32 程序头
// 注意: ELF32 和 ELF64 成员顺序不同
type ElfProgHeader32 struct {
	Type   ProgType // Entry type.
	Off    uint32   // File offset of contents.
	Vaddr  uint32   // Virtual address in memory image.
	Paddr  uint32   // Physical address (not used).
	Filesz uint32   // Size of contents in file.
	Memsz  uint32   // Size of contents in memory.
	Flags  ProgFlag // Access permission flags.
	Align  uint32   // Alignment in memory and file.
}

// ELF64 程序头
// 注意: ELF32 和 ELF64 成员顺序不同
type ElfProgHeader64 struct {
	Type   ProgType
	Flags  ProgFlag
	Off    uint64
	Vaddr  uint64
	Paddr  uint64
	Filesz uint64
	Memsz  uint64
	Align  uint64
}

// ELF Section 头
type ElfShdr struct {
	Name      uint32 // Section name (index into the section header string table).
	Type      uint32 // Section type.
	Flags     uint64 // Section flags.
	Addr      uint64 // Address in memory image.
	Off       uint64 // Offset in file.
	Size      uint64 // Size in bytes.
	Link      uint32 // Index of a related section.
	Info      uint32 // Depends on section type.
	Addralign uint64 // Alignment in bytes.
	Entsize   uint64 // Size of each entry in section.
}

// Section type.
type SectionType uint32

const (
	SHT_NULL           SectionType = 0          // inactive
	SHT_PROGBITS       SectionType = 1          // program defined information
	SHT_SYMTAB         SectionType = 2          // symbol table section
	SHT_STRTAB         SectionType = 3          // string table section
	SHT_RELA           SectionType = 4          // relocation section with addends
	SHT_HASH           SectionType = 5          // symbol hash table section
	SHT_DYNAMIC        SectionType = 6          // dynamic section
	SHT_NOTE           SectionType = 7          // note section
	SHT_NOBITS         SectionType = 8          // no space section
	SHT_REL            SectionType = 9          // relocation section - no addends
	SHT_SHLIB          SectionType = 10         // reserved - purpose unknown
	SHT_DYNSYM         SectionType = 11         // dynamic symbol table section
	SHT_INIT_ARRAY     SectionType = 14         // Initialization function pointers.
	SHT_FINI_ARRAY     SectionType = 15         // Termination function pointers.
	SHT_PREINIT_ARRAY  SectionType = 16         // Pre-initialization function ptrs.
	SHT_GROUP          SectionType = 17         // Section group.
	SHT_SYMTAB_SHNDX   SectionType = 18         // Section indexes (see SHN_XINDEX).
	SHT_LOOS           SectionType = 0x60000000 // First of OS specific semantics
	SHT_GNU_ATTRIBUTES SectionType = 0x6ffffff5 // GNU object attributes
	SHT_GNU_HASH       SectionType = 0x6ffffff6 // GNU hash table
	SHT_GNU_LIBLIST    SectionType = 0x6ffffff7 // GNU prelink library list
	SHT_GNU_VERDEF     SectionType = 0x6ffffffd // GNU version definition section
	SHT_GNU_VERNEED    SectionType = 0x6ffffffe // GNU version needs section
	SHT_GNU_VERSYM     SectionType = 0x6fffffff // GNU version symbol table
	SHT_HIOS           SectionType = 0x6fffffff // Last of OS specific semantics
	SHT_LOPROC         SectionType = 0x70000000 // reserved range for processor
	SHT_MIPS_ABIFLAGS  SectionType = 0x7000002a // .MIPS.abiflags
	SHT_HIPROC         SectionType = 0x7fffffff // specific section header types
	SHT_LOUSER         SectionType = 0x80000000 // reserved range for application
	SHT_HIUSER         SectionType = 0xffffffff // specific indexes
)

// Section flags.
type SectionFlag uint32

const (
	SHF_WRITE            SectionFlag = 0x1        // Section contains writable data.
	SHF_ALLOC            SectionFlag = 0x2        // Section occupies memory.
	SHF_EXECINSTR        SectionFlag = 0x4        // Section contains instructions.
	SHF_MERGE            SectionFlag = 0x10       // Section may be merged.
	SHF_STRINGS          SectionFlag = 0x20       // Section contains strings.
	SHF_INFO_LINK        SectionFlag = 0x40       // sh_info holds section index.
	SHF_LINK_ORDER       SectionFlag = 0x80       // Special ordering requirements.
	SHF_OS_NONCONFORMING SectionFlag = 0x100      // OS-specific processing required.
	SHF_GROUP            SectionFlag = 0x200      // Member of section group.
	SHF_TLS              SectionFlag = 0x400      // Section contains TLS data.
	SHF_COMPRESSED       SectionFlag = 0x800      // Section is compressed.
	SHF_MASKOS           SectionFlag = 0x0ff00000 // OS-specific semantics.
	SHF_MASKPROC         SectionFlag = 0xf0000000 // Processor-specific semantics.
)

// Section compression type.
type CompressionType int

const (
	COMPRESS_ZLIB   CompressionType = 1          // ZLIB compression.
	COMPRESS_ZSTD   CompressionType = 2          // ZSTD compression.
	COMPRESS_LOOS   CompressionType = 0x60000000 // First OS-specific.
	COMPRESS_HIOS   CompressionType = 0x6fffffff // Last OS-specific.
	COMPRESS_LOPROC CompressionType = 0x70000000 // First processor-specific type.
	COMPRESS_HIPROC CompressionType = 0x7fffffff // Last processor-specific type.
)

// ELF32 Compression header.
type CompressionHeader32 struct {
	Type      uint32
	Size      uint32
	Addralign uint32
}

// ELF64 Compression header.
type CompressionHeader64 struct {
	Type      uint32
	_         uint32 /* Reserved. */
	Size      uint64
	Addralign uint64
}
