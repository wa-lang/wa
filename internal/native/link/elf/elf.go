// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package elf

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

// Data is found in Header.Ident[EI_DATA] and Header.Data.
type Data byte

const (
	ELFDATANONE Data = 0 // Unknown data format.
	ELFDATA2LSB Data = 1 // 2's complement little-endian.
	ELFDATA2MSB Data = 2 // 2's complement big-endian.
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
	EM_NONE    Machine = 0   // Unknown machine.
	EM_X86_64  Machine = 62  // Advanced Micro Devices x86-64
	EM_AARCH64 Machine = 183 // ARM 64-bit Architecture (AArch64)
	EM_RISCV   Machine = 243 // RISC-V
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
