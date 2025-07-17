;
; hh2d.nasm: 408-byte, 0-based tiny hello-world Win32 PE .exe
; by pts@fazekas.hu at Sun Jul 26 22:59:33 CEST 2020
;
; Compile: nasm -O0 -f bin -o hh2d.exe hh2d.nasm
;
; It doesn't work on Windows NT 3.1 (empty output, even with
; SubsystemVersion set to 3.10). It works on Wine 5.0, Windows 95, Windows
; NT 4.0 and Windows XP.
;
; This file is based on hh6d.nasm. The file offset of the .text section was
; adjusted so that the entire file gets loaded.
;
; This program doesn't work on ReactOS 0.4.14, even if we reduce
; SKIP_PREFIX_SIZE to 1 and we make the .exe 1 byte longer.
;

; Asserts that we are at offset %1 from the beginning of the input file
%macro aa 1
times $-(%1) times 0 nop
times (%1)-$ times 0 nop
%endmacro

bits 32
cpu 386

FILE_HEADER:
IMAGE_DOS_HEADER:
aa $$+0x0000
; https://github.com/pts/pts-nasm-fullprog/blob/master/pe_stub1.nasm
.mz_signature: dw 'MZ'
.image_size_lo: dw IMAGE_NT_HEADERS
.image_size_hi: dw 1
dw 0, 1, 0x0fff, -1, 1, -1, 0, 8, 0
aa $$+24
.stub_start:
push ss
pop ds
mov ah, 9  ; WRITE_DOLLAR_STDOUT.
db 0xba  ; 16-bit mov dx, ...
dw 6  ; .stub_msg, based on .stub_data.
int 0x21
db 0xb8  ; 16-bit mov ax, ...
dw 0x4c01  ; EXIT(1).
int 0x21
.stub_msg: db 'Not a DOS program.', 13, 10, '$'
times 60-($-$$) db 0
dd IMAGE_NT_HEADERS
aa $$+0x40

IMAGE_BASE equ 0x00400000  ; Variable.
VADDR_TEXT equ 0x1000
SKIP_PREFIX_SIZE equ 0x40  ; Don't load this much of file prefix. Can be 2, 0x10 or 0x40, it doesn't matter.
BSS_SIZE EQU 0

IMAGE_NT_HEADERS:
db 'PE', 0, 0

IMAGE_FILE_HEADER:
Machine: dw 0x14c  ; IMAGE_FILE_MACHINE_I386
NumberOfSections: dw (IMAGE_SECTION_HEADER_end-IMAGE_SECTION_HEADER)/40
TimeDateStamp: dd 0
PointerToSymbolTable: dd 0
NumberOfSymbols: dd 0
SizeOfOptionalHeader: dw IMAGE_OPTIONAL_HEADER32_end-IMAGE_OPTIONAL_HEADER32
IMAGE_FILE_RELOCS_STRIPPED equ 1
IMAGE_FILE_EXECUTABLE_IMAGE equ 2
IMAGE_FILE_LINE_NUMS_STRIPPED equ 4
IMAGE_FILE_LOCAL_SYMS_STRIPPED equ 8
IMAGE_FILE_BYTES_REVERSED_LO equ 0x80  ; Deprecated, shouldn't be specified.
IMAGE_FILE_32BIT_MACHINE equ 0x100
IMAGE_FILE_DEBUG_STRIPPED equ 0x200
IMAGE_FILE_DLL equ 0x2000  ; Shouldn't be specified for .exe.
Characteristics: dw IMAGE_FILE_RELOCS_STRIPPED|IMAGE_FILE_EXECUTABLE_IMAGE|IMAGE_FILE_LINE_NUMS_STRIPPED|IMAGE_FILE_LOCAL_SYMS_STRIPPED|IMAGE_FILE_32BIT_MACHINE|IMAGE_FILE_DEBUG_STRIPPED

IMAGE_OPTIONAL_HEADER32:
Magic: dw 0x10b  ; IMAGE_NT_OPTIONAL_HDR32_MAGIC
MajorLinkerVersion: db 6
MinorLinkerVersion: db 0
SizeOfCode: dd 0x00000000
SizeOfInitializedData: dd 0x00000000
SizeOfUninitializedData: dd 0x00000000
AddressOfEntryPoint: dd _start+(VADDR_TEXT)  ; Also called starting address.
BaseOfCode: dd VADDR_TEXT
BaseOfData: dd VADDR_TEXT
ImageBase: dd IMAGE_BASE
SectionAlignment: dd 0x1000  ; Single allowed value for Windows XP.
FileAlignment: dd 0x200  ; Minimum value for Windows NT 3.1.
MajorOperatingSystemVersion: dw 4
MinorOperatingSystemVersion: dw 0
MajorImageVersion: dw 0
MinorImageVersion: dw 0
;MajorSubsystemVersion: dw 3   ; Windows NT 3.1.
;MinorSubsystemVersion: dw 10  ; Windows NT 3.1.
MajorSubsystemVersion: dw 4   ; Windows 95 and above.
MinorSubsystemVersion: dw 0   ; Windows 95 and above.
Win32VersionValue: dd 0
SizeOfImage: dd VADDR_TEXT+FILE_end-FILE_HEADER+EXTRA_BSS_SIZE+BSS_SIZE
SizeOfHeaders: dd FILE_end-FILE_HEADER
CheckSum: dd 0
Subsystem: dw 3  ; IMAGE_SUBSYSTEM_WINDOWS_CUI; gcc -mconsole
DllCharacteristics: dw 0
SizeOfStackReserve: dd 0x00100000
SizeOfStackCommit: dd 0x00001000
SizeOfHeapReserve: dd 0x100000  ; Why not 0?
SizeOfHeapCommit: dd 0x1000  ; Why not 0?
LoaderFlags: dd 0
NumberOfRvaAndSizes: dd (IMAGE_DATA_DIRECTORY_end-IMAGE_DATA_DIRECTORY)/8
IMAGE_DATA_DIRECTORY:
IMAGE_DIRECTORY_ENTRY_EXPORT:
.VirtualAddress: dd 0x00000000
.Size: dd 0x00000000
IMAGE_DIRECTORY_ENTRY_IMPORT:
.VirtualAddress: dd IMAGE_IMPORT_DESCRIPTORS+(VADDR_TEXT)
.Size: dd IMAGE_IMPORT_DESCRIPTORS_end-IMAGE_IMPORT_DESCRIPTORS

; Overlapping 5+5 dwords of IMAGE_IMPORT_DESCRIPTORS with 5 ignored entries of
; IMAGE_DATA_DIRECTORY.
IMAGE_IMPORT_DESCRIPTORS:
IMAGE_IMPORT_DESCRIPTOR_0:
.OriginalFirstThunk: dd IMPORT_ADDRESS_TABLE+(VADDR_TEXT)
.TimeDateStamp: dd 0
.ForwarderChain: dd 0
.Name: dd NAME_KERNEL32_DLL+(VADDR_TEXT)
.FirstThunk: dd IMPORT_ADDRESS_TABLE+(VADDR_TEXT)
IMAGE_IMPORT_DESCRIPTOR_1:  ; Last Import directory table, marks end-of-list.
;IID_BSS_SIZE equ 4*5
;IMAGE_IMPORT_DESCRIPTORS_end equ $+IID_BSS_SIZE
dd 0, 0, 0, 0, 0  ; Same fields as above, filled with 0s.
IMAGE_IMPORT_DESCRIPTORS_end:

; Overlapping NAME_KERNEL32_DLL with 2 ignored entries of IMAGE_DATA_DIRECTORY.
dd 0
NAME_KERNEL32_DLL: db 'kernel32.dll'

dd 0, 0  ; 3 more ignored entries of IMAGE_DATA_DIRECTORY. Needed by Windows 95, which needs >= 10 entries.

IMAGE_DATA_DIRECTORY_end:
IMAGE_OPTIONAL_HEADER32_end:

IMAGE_SECTION_HEADER:
IMAGE_SECTION_HEADER__0:
.Name: db '.text'
times 8-($-.Name) db 0
.VirtualSize: dd FILE_end-FILE_HEADER+EXTRA_BSS_SIZE+BSS_SIZE
.VirtualAddress: dd VADDR_TEXT  ; Adding +SKIP_PREFIX_SIZE breaks it on Wine 5.0, Windows NT 4.0 etc.
; The trick is to add SKIP_PREFIX_SIZE below. The actual value is ignored by
; Windows NT 4.0, Windows 95 and Windows XP, so the entire file will be loaded to
; IMAGE_BASE+VADDR_TEXT. This is exactly what we want. Unfortunately, the trick
; doesn't work on Windows NT 3.1, there PointerToRawData must be divisible by 0x200.
.SizeOfRawData: dd FILE_end-FILE_HEADER-SKIP_PREFIX_SIZE
.PointerToRawData: dd SKIP_PREFIX_SIZE
.PointerToRelocations: dd 0
.PointerToLineNumbers: dd 0
.NumberOfRelocations: dw 0
.NumberOfLineNumbers: dw 0
IMAGE_SCN_CNT_CODE equ 0x20
IMAGE_SCN_MEM_EXECUTE equ 0x20000000
IMAGE_SCN_MEM_READ equ 0x40000000
IMAGE_SCN_CNT_INITIALIZED_DATA equ 0x40
IMAGE_SCN_MEM_WRITE equ 0x80000000
.Characteristics: dd IMAGE_SCN_CNT_CODE|IMAGE_SCN_MEM_EXECUTE|IMAGE_SCN_MEM_READ|IMAGE_SCN_CNT_INITIALIZED_DATA|IMAGE_SCN_MEM_WRITE
IMAGE_SECTION_HEADER_end:

; The `0, 0, ' is the .Hint.
NAME_GetStdHandle: db 0, 0, 'GetStdHandle', 0
NAME_WriteFile: db 0, 0, 'WriteFile', 0
NAME_ExitProcess: db 0, 0, 'ExitProcess', 0

_start:
push strict byte -11  ; STD_OUTPUT_HANDLE.
call [__imp__GetStdHandle@4+(IMAGE_BASE+VADDR_TEXT)]
;push eax  ; Save stdout handle.
push eax  ; Value is arbitrary, we allocate an output variable on the stack.
mov ebx, esp
push strict byte 0  ; Argument 5: lpOverlapped = 0.
push ebx  ; Argument 4: Address of the output variable.
push strict byte message_end-message  ; Argument 3: message size.
push strict dword message+(IMAGE_BASE+VADDR_TEXT)  ; Argument 2: message.
push eax  ; if it was saved, strict dword [esp+20]  ; Argument 1: Stdout handle.
call [__imp__WriteFile@20+(IMAGE_BASE+VADDR_TEXT)]
push byte 0  ; EXIT_SUCCESS == 0.
call [__imp__ExitProcess@4+(IMAGE_BASE+VADDR_TEXT)]
;add esp, 8  ; Too late, we've already exited.
;ret  ; Too late, we've already exited.

;SECTION_DATA:

message:
;db 'Hello, World! MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM', 13, 10, 'MSG', 13, 10
db 'Hello, World!', 13, 10
message_end:

; Because of the modification, this must start after SECTION_TEXT.
IMPORT_ADDRESS_TABLE:  ; Import address table. Modified by the PE loader before jumping to _entry.
__imp__GetStdHandle@4: dd NAME_GetStdHandle+(VADDR_TEXT)
__imp__WriteFile@20: dd NAME_WriteFile+(VADDR_TEXT)
__imp__ExitProcess@4 dd NAME_ExitProcess+(VADDR_TEXT)
;dd 0  ; Marks end-of-list.
EXTRA_BSS_SIZE equ 4  ; For the end-of-list marker above. Must be at the end.
IMPORT_ADDRESS_TABLE_end:

FILE_end:
