// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package pe

const (
	PE32Magic = 0x10b
	PE64Magic = 0x20b

	DefaultSectionAlignment = 0x1000 // 内存中按 4KB 对齐
	DefaultFileAlignment    = 0x200  // 磁盘上按 512 字节对齐

	DefaultImageBase = 0x140000000
)

const (
	KDosHeaderSize = 64
)

// DOS文件头部
type DosHeader struct {
	Magic          [2]byte // 'M', 'Z'
	_              [58]byte
	NTHeaderOffset uint32 // NtHeader64 开始地址, 默认为 sizeof(DosHeader) == 64
}

// Exe头部
type NtHeader64 struct {
	Signature      [4]byte          // 'P', 'E', 0, 0
	FileHeader     FileHeader       // 20 字节
	OptionalHeader OptionalHeader64 // 240 字节
}

// PE文件头部
type FileHeader struct {
	Machine              uint16 // x64: 0x8664
	NumberOfSections     uint16
	TimeDateStamp        uint32
	PointerToSymbolTable uint32
	NumberOfSymbols      uint32
	SizeOfOptionalHeader uint16 // x64: 240
	Characteristics      uint16
}

// X64 PE文件可选头
// 对于 exe 文件是必须包含的
type OptionalHeader64 struct {
	Magic                       uint16
	MajorLinkerVersion          uint8
	MinorLinkerVersion          uint8
	SizeOfCode                  uint32
	SizeOfInitializedData       uint32
	SizeOfUninitializedData     uint32
	AddressOfEntryPoint         uint32
	BaseOfCode                  uint32
	ImageBase                   uint64
	SectionAlignment            uint32 // 段数据内存地址对齐, 通常 0x1000 (4KB), 不能小于文件内对齐
	FileAlignment               uint32 // 段数据文件内地址对齐, 通常 0x200 (512B)
	MajorOperatingSystemVersion uint16
	MinorOperatingSystemVersion uint16
	MajorImageVersion           uint16
	MinorImageVersion           uint16
	MajorSubsystemVersion       uint16
	MinorSubsystemVersion       uint16
	Win32VersionValue           uint32
	SizeOfImage                 uint32
	SizeOfHeaders               uint32
	CheckSum                    uint32
	Subsystem                   uint16
	DllCharacteristics          uint16
	SizeOfStackReserve          uint64
	SizeOfStackCommit           uint64
	SizeOfHeapReserve           uint64
	SizeOfHeapCommit            uint64
	LoaderFlags                 uint32
	NumberOfRvaAndSizes         uint32
	DataDirectory               [16]DataDirectory // DLL 导入库在此定义
}

// 程序的元数据
// 无外部依赖的exe文件可以忽略
type DataDirectory struct {
	VirtualAddress uint32
	Size           uint32
}

// X64 PE 文件段头
type SectionHeader64 struct {
	Name                 [8]byte // 必须是 8 字节，如 ".text\0\0\0"
	VirtualSize          uint32  // 段数据大小
	VirtualAddress       uint32  // 内存偏移 (RVA)
	Size                 uint32  // 文件中对齐后的大小, 需要满足 OptionalHeader64.SectionAlignment 对齐
	Offset               uint32  // 文件中对齐后的偏移, 需要满足 OptionalHeader64.FileAlignment 对齐
	PointerToRelocations uint32  // 填 0
	PointerToLineNumbers uint32  // 填 0
	NumberOfRelocations  uint16  // 填 0
	NumberOfLineNumbers  uint16  // 填 0
	Characteristics      uint32  // 权限 (代码段: 0x60000020, 数据段: 0xC0000040)
}

// 导入符号字典
type ImportDirectory struct {
	// 指向 导入名称表的 RVA(相对虚拟地址)
	// 指向一个数组, 数组里记录了要调用的所有函数的名字(比如 ExitProcess)
	OriginalFirstThunk uint32

	// 时间戳
	// 如果该值为 0xFFFFFFFF, 表示该 DLL 已被绑定, 函数地址已提前写死在文件中.
	// 但在现代开发和你的手搓编译器中, 通常直接填 0 即可
	TimeDateStamp uint32

	// 转发链 (通常被设置为 0)
	// 用于处理 API 转发.
	// 例如调用 DLL A 里的函数, 但 DLL A 实际上把这个请求转交给了 DLL B
	// 这个字段在大多数情况下不需要手动处理
	ForwarderChain uint32

	// 指向 DLL 名称字符串 的 RVA
	// 例如 "KERNEL32.dll"
	Name uint32

	// 指向 导入地址表 (Import Address Table, 简称 IAT) 的 RVA
	// 和 OriginalFirstThunk 指向的内容几乎一样, 也是函数名
	// 加载运行后, Windows 系统加载器会根据 OriginalFirstThunk 里的名字,
	// 找到函数在内存里的真实地址，并覆盖写进 FirstThunk 指向的这个位置.
	// 这是在汇编代码里 call [ExitProcess_IAT] 能跳到正确地方的原因.
	FirstThunk uint32

	// 非物理结构, 后来拼接
	// ExitProcess:KERNEL32.dll
	dll string
}

// 机器码类型
const (
	IMAGE_FILE_MACHINE_UNKNOWN = 0x0
	IMAGE_FILE_MACHINE_AMD64   = 0x8664
	IMAGE_FILE_MACHINE_I386    = 0x14c
)

// FileHeader.Characteristics 标志位组合
const (
	IMAGE_FILE_EXECUTABLE_IMAGE    = 0x0002 // 可执行程序
	IMAGE_FILE_LARGE_ADDRESS_AWARE = 0x0020 // 64 位地址空间
)

// OptionalHeader64.Subsystem 标志位组合
const (
	IMAGE_SUBSYSTEM_UNKNOWN     = 0
	IMAGE_SUBSYSTEM_WINDOWS_GUI = 2 // GUI 程序
	IMAGE_SUBSYSTEM_WINDOWS_CUI = 3 // 命令行程序
)

// 段数据标志位
const (
	IMAGE_SCN_CNT_CODE               = 0x00000020
	IMAGE_SCN_CNT_INITIALIZED_DATA   = 0x00000040
	IMAGE_SCN_CNT_UNINITIALIZED_DATA = 0x00000080
	IMAGE_SCN_LNK_COMDAT             = 0x00001000
	IMAGE_SCN_MEM_DISCARDABLE        = 0x02000000
	IMAGE_SCN_MEM_EXECUTE            = 0x20000000
	IMAGE_SCN_MEM_READ               = 0x40000000
	IMAGE_SCN_MEM_WRITE              = 0x80000000
)

// 导入表信息
const (
	IMAGE_DIRECTORY_ENTRY_EXPORT         = 0
	IMAGE_DIRECTORY_ENTRY_IMPORT         = 1
	IMAGE_DIRECTORY_ENTRY_RESOURCE       = 2
	IMAGE_DIRECTORY_ENTRY_EXCEPTION      = 3
	IMAGE_DIRECTORY_ENTRY_SECURITY       = 4
	IMAGE_DIRECTORY_ENTRY_BASERELOC      = 5
	IMAGE_DIRECTORY_ENTRY_DEBUG          = 6
	IMAGE_DIRECTORY_ENTRY_ARCHITECTURE   = 7
	IMAGE_DIRECTORY_ENTRY_GLOBALPTR      = 8
	IMAGE_DIRECTORY_ENTRY_TLS            = 9
	IMAGE_DIRECTORY_ENTRY_LOAD_CONFIG    = 10
	IMAGE_DIRECTORY_ENTRY_BOUND_IMPORT   = 11
	IMAGE_DIRECTORY_ENTRY_IAT            = 12
	IMAGE_DIRECTORY_ENTRY_DELAY_IMPORT   = 13
	IMAGE_DIRECTORY_ENTRY_COM_DESCRIPTOR = 14
)
