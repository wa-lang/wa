;
; hh2.nasm: 402-byte, 0-based tiny hello-world Win32 PE .exe
; by pts@fazekas.hu at Sat Jan 13 11:53:58 CET 2018
;
; By 0-based we mean that Windows loads the entire .exe file to memory, even
; the headers; and the first section starts at file offset 0. This lets us
; get rid of lots of padding compared to solutions which are not 0-based.
;
; How to compile hh2.exe:
;
;   $ nasm -f bin -o hh2.exe hh2.nasm
;   $ chmod 755 hh2.exe  # For QEMU Samba server.
;   $ ndisasm -b 32 -e 0x114 -o 0x403000 hh2.exe
;
; hh2.asm was inspired by the 268-byte .exe on
; https://www.codejuggle.dj/creating-the-smallest-possible-windows-executable-using-assembly-language/
; . The fundamental difference is that hh2.exe works on Windows XP ... Windows
; 10 (but not on Windows NT 3.1 and Windows 95), while the program above
; doesn't work on Windows XP.
;
; The generated hh2.exe works on:
;
; * Wine 1.6.2 on Linux.
; * Windows XP SP3, 32-bit: Microsoft Windows XP [Version 5.1.2600]
; * Windows 10 64-bit: Microsoft Windows [Version 10.0.16299.192]
;
; It doesn't work on Windows NT 3.51 (not even after changing the
; SubsystemVersion to 3.10), and it doesn't work on Windows 95 either.
;
; Output .exe file size in bytes (approximately):
;
;   len(text_bytes) + len(data_bytes) + len(rodata_bytes) +
;   + 320
;   + sum(len(name) for name in imported_names) + 2 * len(imported_names) - 1
;   + 8 * len(imported_names) + 6
;   + sum(len(name) for name in library_names) + len(library_names) - (len('kernel32') + 1) - ...
;     # Example library_name: 'kernel32'.
;     # About the ...: some more strings can be inlined to the 'x' paddings.
;   + 20 * len(library_names)
;
; Assumptions:
;
; * len(imported_names) >= 1: ['ExitProcess']
; * len(library_names) >= 1: ['kernel32']
;
;
; hh2.nasm was created like this:
;
; * The .nasm source in
;   https://www.codejuggle.dj/creating-the-smallest-possible-windows-executable-using-assembly-language/
;   was studied.
; * The .exe created from the hello-world .c program in
;   https://stackoverflow.com/questions/42022132/how-to-create-tiny-pe-win32-executables-using-mingw
;   was manually converted back to .nasm, pointers changed to symbols and
;   address computations added one-by-one.
; * The 2nd .nasm file was gradually changed to resemble the 1st .nasm file,
;   while making sure that the generated .exe still works on Windows XP.
;
bits 32
imagebase equ 0x400000  ; Default base since Windows 95.
textbase equ imagebase + 0x3000
file_alignment equ 0x200
bits 32
org 0  ; Can be anything, this file doesn't depend on it.

_filestart:
_text:

IMAGE_DOS_HEADER:  ; Truncated, breaks file(1) etc.
db 'MZ'
_KERNEL32_str: db 'kernel32', 0  ; 'KERNEL32' and 'KERNEL32.dll' also work.
times 0xc - ($-$$) db 0

IMAGE_NT_HEADERS:
Signature: dw 'PE', 0

IMAGE_FILE_HEADER:
Machine: dw 0x14c  ; IMAGE_FILE_MACHINE_I386
NumberOfSections: dw (_headers_end - _sechead) / 40  ; Windows XP needs >= 3.
TimeDateStamp: dd 0x00000000
PointerToSymbolTable: dd 0x00000000
NumberOfSymbols: dd 0x00000000
SizeOfOptionalHeader: dw _datadir_end - _opthd  ; Windows XP needs >= 0x78.
IMAGE_FILE_RELOCS_STRIPPED equ 1
IMAGE_FILE_EXECUTABLE_IMAGE equ 2
IMAGE_FILE_LINE_NUMS_STRIPPED equ 4
IMAGE_FILE_LOCAL_SYMS_STRIPPED equ 8
IMAGE_FILE_BYTES_REVERSED_LO equ 0x80  ; Deprecated, shouldn't be specified.
IMAGE_FILE_32BIT_MACHINE equ 0x100
IMAGE_FILE_DEBUG_STRIPPED equ 0x200
IMAGE_FILE_DLL equ 0x2000  ; Shouldn't be specified for .exe.
Characteristics: dw IMAGE_FILE_RELOCS_STRIPPED|IMAGE_FILE_EXECUTABLE_IMAGE|IMAGE_FILE_LINE_NUMS_STRIPPED|IMAGE_FILE_LOCAL_SYMS_STRIPPED|IMAGE_FILE_32BIT_MACHINE|IMAGE_FILE_DEBUG_STRIPPED

_opthd:
IMAGE_OPTIONAL_HEADER32:
Magic: dw 0x10b  ; IMAGE_NT_OPTIONAL_HDR32_MAGIC
MajorLinkerVersion: db 0
MinorLinkerVersion: db 0
SizeOfCode: dd 0x00000000
SizeOfInitializedData: dd 0x00000000
SizeOfUninitializedData: dd 0x00000000
AddressOfEntryPoint: dd (textbase - imagebase) + (_entry - _text)
BaseOfCode: dd 0x00000000
BaseOfData: dd (IMAGE_NT_HEADERS - _filestart)  ; Overlaps with: IMAGE_DOS_HEADER.e_lfanew.
ImageBase: dd imagebase
SectionAlignment: dd 0x1000  ; Minimum value for Windows XP.
%if file_alignment == 0 || file_alignment & (file_alignment - 1)
%error Invalid file_alignment, must be a power of 2.
times -1 nop  ; Force error even in NASM 0.98.39.
%endif
%if file_alignment < 0x200
%error Windows XP needs file_alignment >= 0x200
times -1 nop  ; Force error even in NASM 0.98.39.
%endif
FileAlignment: dd file_alignment  ; Minimum value for Windows XP.
MajorOperatingSystemVersion: dw 4
MinorOperatingSystemVersion: dw 0
MajorImageVersion: dw 1
MinorImageVersion: dw 0
MajorSubsystemVersion: dw 4
MinorSubsystemVersion: dw 0
Win32VersionValue: dd 0
SizeOfImage: dd (textbase - imagebase) + (_eof + bss_size - _text)  ; Wine rounds it up to a multiple of 0x1000, and loads and maps that much.
SizeOfHeaders: dd _headers_end - _filestart  ; Windows XP needs > 0.
CheckSum: dd 0
Subsystem: dw 3  ; IMAGE_SUBSYSTEM_WINDOWS_CUI; gcc -mconsole
DllCharacteristics: dw 0
SizeOfStackReserve: dd 0x00100000
SizeOfStackCommit: dd 0x00001000
SizeOfHeapReserve: dd 0
SizeOfHeapCommit: dd 0
LoaderFlags: dd 0
; If we hardcode 2 here, on Windows XP we can put arbitrary bytes to
; IMAGE_DIRECTORY_ENTRY_RESOURCE.VirtualAddress and .Size. If we put
; 3 here (autogenerated), then the values must be 0.
;NumberOfRvaAndSizes: dd (_datadir_end - _datadir) / 8  ; Number of IMAGE_DATA_DIRECTORY entries below.
NumberOfRvaAndSizes: dd 2

_datadir:
DataDirectory:
IMAGE_DIRECTORY_ENTRY_EXPORT:
.VirtualAddress: dd 0x00000000
.Size: dd 0x00000000
IMAGE_DIRECTORY_ENTRY_IMPORT:
.VirtualAddress: dd (textbase - imagebase) + (_idescs - _text)
.Size: dd _idata_data_end - _idata
IMAGE_DIRECTORY_ENTRY_RESOURCE:
;.VirtualAddress_AndSize: db 'xxxxxxxx'  ; Arbitrary string OK here, continues below in IMAGE_SECTION_HEADER__0.name. db 'tiny.exe'
IMAGE_IMPORT_BY_NAME_GetStdHandle:
.Hint: dw 0
;.Name: db 'GetStdHandle'  ; Terminated by below.
.VirtualAddress_AndSize: db 'GetStd'
%if 0
; Changing all 0x78787878 to 0 below may fix startup errors.
IMAGE_DIRECTORY_ENTRY_EXCEPTION:
.VirtualAddress: dd 0x78787878
.Size: dd 0x78787878
IMAGE_DIRECTORY_ENTRY_SECURITY:
.VirtualAddress: dd 0x78787878
.Size: dd 0x78787878
IMAGE_DIRECTORY_ENTRY_BASERELOC:
.VirtualAddress: dd 0x78787878
.Size: dd 0x78787878
IMAGE_DIRECTORY_ENTRY_DEBUG:
.VirtualAddress: dd 0x78787878
.Size: dd 0x00000000
IMAGE_DIRECTORY_ENTRY_ARCHITECTURE:
.VirtualAddress: dd 0x00000000
.Size: dd 0x00000000
IMAGE_DIRECTORY_ENTRY_GLOBALPTR:
.VirtualAddress: dd 0x00000000
.Size: dd 0x78787878
IMAGE_DIRECTORY_ENTRY_TLS:
.VirtualAddress: dd 0x78787878
.Size: dd 0x78787878
IMAGE_DIRECTORY_ENTRY_LOAD_CONFIG:
.VirtualAddress: dd 0x78787878
.Size: dd 0x78787878
IMAGE_DIRECTORY_ENTRY_BOUND_IMPORT:
.VirtualAddress: dd 0x78787878
.Size: dd 0x78787878
IMAGE_DIRECTORY_ENTRY_IAT:
.VirtualAddress: dd 0x78787878
.Size: dd 0x78787878
IMAGE_DIRECTORY_ENTRY_DELAY_IMPORT:
.VirtualAddress: dd 0x78787878
.Size: dd 0x78787878
 Missing:
IMAGE_DIRECTORY_ENTRY_COM_DESCRIPTOR:
.VirtualAddress: dd 0x78787878
.Size: dd 0x78787878
IMAGE_DIRECTORY_ENTRY_RESERVED:
.VirtualAddress: dd 0x78787878
.Size: dd 0x78787878
%endif
_datadir_end:

_sechead:

IMAGE_SCN_CNT_CODE equ 0x20
IMAGE_SCN_MEM_EXECUTE equ 0x20000000
IMAGE_SCN_MEM_READ equ 0x40000000
IMAGE_SCN_CNT_INITIALIZED_DATA equ 0x40
IMAGE_SCN_MEM_WRITE equ 0x80000000
IMAGE_SCN_ALIGN_4BYTES equ 0x00300000

IMAGE_SECTION_HEADER__0:
;.Name: db 'xxxxxxx', 0  ; Arbitrary ASCIIZ string OK here.  ; db '.dummy1', 0
.Name: db 'Handle', 0, 0  ; Pad to 16 bytes, to the end of .Name.
.VirtualSize: dd 1  ; Must be positive for Windows XP.
.VirtualAddress: dd 0x1000  ; Must be positive and divisible by 0x1000 for Windows XP.
.SizeOfRawData: dd 0
.PointerToRawData: dd 0
.PointerToRelocations: dd 0
.PointerToLineNumbers: dd 0
.NumberOfRelocations: dw 0
.NumberOfLineNumbers: dw 0
.Characteristics: dd IMAGE_SCN_ALIGN_4BYTES|IMAGE_SCN_MEM_READ|IMAGE_SCN_CNT_INITIALIZED_DATA|IMAGE_SCN_MEM_WRITE

IMAGE_SECTION_HEADER__1:
.Name: db 'xxxxxxx', 0  ;  Arbitaray ASCIIZ string OK here.  ; db '.dummy2', 0
.VirtualSize: dd 1  ; Must be positive for Windows XP.
.VirtualAddress: dd 0x2000  ; Must be positive, divisible by 0x1000, and larger then the prev .VirtualAddress for Windows XP.
.SizeOfRawData: dd 0
.PointerToRawData: dd 0
.PointerToRelocations: dd 0
.PointerToLineNumbers: dd 0
.NumberOfRelocations: dw 0
.NumberOfLineNumbers: dw 0
.Characteristics: dd IMAGE_SCN_ALIGN_4BYTES|IMAGE_SCN_MEM_WRITE|IMAGE_SCN_MEM_READ|IMAGE_SCN_CNT_INITIALIZED_DATA

IMAGE_SECTION_HEADER__2:
.Name: db 'xxxxxxx', 0  ;  Arbitrary ASCIIZ string OK here.  ; db '.text', 0, 0, 0
.VirtualSize: dd (_eof - _text) + bss_size
%if (textbase - imagebase) & 0xfff
%error _text doesn't start at page boundary, needed by Windows XP.
times -1 nop  ; Force error even in NASM 0.98.39.
%endif
%if (textbase - imagebase) <= 0x2000
%error _text doesn't start later than the previous sections, needed by Windows XP.
times -1 nop  ; Force error even in NASM 0.98.39.
%endif
.VirtualAddress: dd textbase - imagebase
.SizeOfRawData: dd _eof - _text - 2 ;* !(_text - _filestart)
; It's not true that PointerToRawData must be a muliple of file_alignment.
; What is true is that Windows XP ignores the low 9 bits
; (log2(file_alignment)) of PointerToRawData (except when adding it to
; SizeOfRawData), and also loads the preceding few bytes. The reason why we
; don't set PointerToRawData here is that if it's 0, then Windows XP doesn't
; load anything for this section.
.PointerToRawData: dd _text - _filestart + 2 ;* !(_text - _filestart)
.PointerToRelocations: dd 0
.PointerToLineNumbers: dd 0
.NumberOfRelocations: dw 0
.NumberOfLineNumbers: dw 0
.Characteristics: dd IMAGE_SCN_ALIGN_4BYTES|IMAGE_SCN_CNT_CODE|IMAGE_SCN_MEM_EXECUTE|IMAGE_SCN_MEM_READ|IMAGE_SCN_MEM_WRITE

_headers_end:
; We can check it only this late, when _headers_end is defined.
%if (_headers_end - _sechead) % 40 != 0
%error Multiples of IMAGE_SECTION_HEADER needed.
times -1 nop  ; Force error even in NASM 0.98.39.
%endif
%if (_headers_end - _sechead) / 40 < 3
; Please note that hh3t.golden.exe has only 2 sections, and it still works
; on Windows XP (Tiny Windows XP in QEMU 4.2.0). So it's unclear why this
; doesn't work with only 2 sections.
%error Windows XP in this setting needs at least 3 sections.
times -1 nop  ; Force error even in NASM 0.98.39.
%endif

;times 0x200 - ($-$$) db 'x'
;times 0x100 db 'y'  ; Doesn't work, _text is not aligned properly.
;times 0x200 db 'y'  ; Works, making the .exe larger.
;_text:

_entry:
; Arguments pushed in reverse order, popped by the callee.
; WINBASEAPI HANDLE WINAPI GetStdHandle (DWORD nStdHandle);
; HANDLE hfile = GetStdHandle(STD_OUTPUT_HANDLE);
push byte -11                ; STD_OUTPUT_HANDLE
call [textbase + (__imp__GetStdHandle@4 - _text)]
; Arguments pushed in reverse order, popped by the callee.
; WINBASEAPI WINBOOL WINAPI WriteFile (HANDLE hFile, LPCVOID lpBuffer, DWORD nNumberOfBytesToWrite, LPDWORD lpNumberOfBytesWritten, LPOVERLAPPED lpOverlapped);
; DWORD bw;
push eax                     ; Value does't matter.
mov ecx, esp
push byte 0                  ; lpOverlapped
push ecx                     ; lpNumberOfBytesWritten = &dw
push byte (_msg_end - _msg)  ; nNumberOfBytesToWrite
push textbase + (_msg - _text)  ; lpBuffer
push eax                     ; hFile = hfile
call [textbase + (__imp__WriteFile@20 - _text)]
;pop eax                     ; This would pop dw. Needed for cleanup.
; Arguments pushed in reverse order, popped by the callee.
; WINBASEAPI DECLSPEC_NORETURN VOID WINAPI ExitProcess(UINT uExitCode);
push byte 0                  ; uExitCode
call [textbase + (__imp__ExitProcess@4 - _text)]

_data:
_msg:
db 'Hello, World!', 13, 10
_msg_end:

; This can be before of after _entry, it doesn't matter.
_idata:  ; Relocations, IMAGE_DIRECTORY_ENTRY_IMPORT data.
_hintnames:
dd (textbase - imagebase) + (IMAGE_IMPORT_BY_NAME_ExitProcess - _text)
dd (textbase - imagebase) + (IMAGE_IMPORT_BY_NAME_GetStdHandle - _text)
dd (textbase - imagebase) + (IMAGE_IMPORT_BY_NAME_WriteFile - _text)
dd 0  ; Marks end-of-list.
_iat:  ; Modified by the PE loader before jumping to _entry.
__imp__ExitProcess@4:  dd (textbase - imagebase) + (IMAGE_IMPORT_BY_NAME_ExitProcess - _text)
__imp__GetStdHandle@4: dd (textbase - imagebase) + (IMAGE_IMPORT_BY_NAME_GetStdHandle - _text)
__imp__WriteFile@20:   dd (textbase - imagebase) + (IMAGE_IMPORT_BY_NAME_WriteFile - _text)
dw 0  ; Marks end-of-list, 2nd half of the dd is the dw below.
IMAGE_IMPORT_BY_NAME_ExitProcess:
.Hint: dw 0
.Name: db 'ExitProcess'  ; Terminated by the subsequent .Hint.
IMAGE_IMPORT_BY_NAME_WriteFile:
.Hint: dw 0
.Name: db 'WriteFile'  ; Terminated below.
db 0  ; Terminates last .Name.

_idescs:
IMAGE_IMPORT_DESCRIPTOR__0:
.OriginalFirstThunk: dd (textbase - imagebase) + (_hintnames - _text)
.TimeDateStamp: dd 0
.ForwarderChain: dd 0
.Name: dd (textbase - imagebase) + (_KERNEL32_str - _text)
.FirstThunk: dd (textbase - imagebase) + (_iat - _text)

_idata_data_end:
_eof:
;bss_size equ 0
;IMAGE_IMPORT_DESCRIPTOR__1:  ; Empty, marks end-of-list.
;.OriginalFirstThunk: dd 0
;.TimeDateStamp: dd 0
;.ForwarderChain: dd 0
;.Name: dd 0
;.FirstThunk: dd 0
;_idata_end:
bss_size equ 20  ; _idata_end - _eof

%if (_text - _filestart) & (file_alignment - 1)
%error _text is not aligned to file_alignment, needed by Windows XP.
times -1 nop  ; Force error even in NASM 0.98.39.
%endif
