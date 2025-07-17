pts-tinype: tiny hello-world Win32 PE .exe
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
pts-tinype is a set of tiny hello-world Win32 PE .exe executables for the
console (Command Prompt), with assembly source code. The smallest one,
hh2.golden.exe is just 402 bytes large, and it runs on Windows XP ...
Windows 10. The smallest one which runs on all Win32 systems (Windows NT 3.1
to Windows 10), hh6d.golden.exe, is 584 bytes.

How to run:

* Download and run hh2.golden.exe in the Command Prompt of any 32-bit (i386)
  or 64-bit (amd64, x86_64) Windows system or Wine. (It has been tested and
  it works on Windows XP, Windows 10 and Wine 1.6.2.)
* Alternatively, download and run hh5.golden.exe on Windows NT 3.1, Windows
  95, ..., Windows XP, ..., Windows 10 and Wine. It should work everywhere.
* Alternatively, if you don't have a Windows system to try it on, run it
  with Wine.
* Alternatively, if you don't have a Windows system to try it on, run it on
  a virtual machine running Windows. Example Windows XP virtual machine with
  QEMU:
  http://ptspts.blogspot.com/2017/09/how-to-run-windows-xp-on-linux-using-qemu-and-kvm.html

Size and compatibility matrix:

                     hh1   hh2   hh2d  hh3gf hh3tf hh3wf hh3tg hh3tw hh4t  hh6a  hh6b  hh6c  hh6d  hh6r  hh7
---------------------------------------------------------------------------------------------------------------
size (bytes)         268   402   408   2048  1536  3072  3072  3072  610   1536  1024  713   584   584   584
Win32s 1.25a         --    --    --    --    --    bg    yes   yes   yes   --    --    --    --    --    --
Wine 5.0, 1.6.2      yes   yes   yes   yes   yes   yes   yes   yes   yes   yes   yes   yes   yes   yes   yes
ReactOS 0.4.14       --    --    --    yes   yes   yes   yes   yes   yes   yes   yes   yes   yes   yes   yes
Windows NT 3.1       --    --    --    yes   yes   yes   yes   yes   yes   yes   yes   yes   yes   yes   yes
Windows NT 3.5       --    --    yes   yes   yes   yes   yes   yes   yes   yes   yes   yes   yes   yes   yes
Windows 95           --    --    yes   yes   yes   yes   yes   yes   yes   yes   yes   yes   yes   yes   yes
Windows NT 4.0       --    --    yes   yes   yes   yes   yes   yes   yes   yes   yes   yes   yes   yes   yes
Windows XP           --    yes   yes   yes   yes   yes   yes   yes   yes   yes   yes   yes   yes   yes   yes
Windows 7            yes   yes   yes   yes   yes   yes   yes   yes   yes   yes   yes   yes   yes   yes   yes
Windows 10 2020-07   --    yes   yes   yes   yes   yes   yes   yes   yes   yes   yes   yes   yes   yes   yes

``bg'' means that the program runs in the background, and the message it
prints is not displayed in any window.

Win32s doesn't have a console window where standard output of console
programs could be displayed. It also requires a relocation table (so that it
can load the .exe to any address), and currently on hh3wf.golden.exe contains
a relocation table.

Variants:

* hh1.golden.exe (268 bytes): Doesn't work on Windows NT 3.1, Windows 95,
  Windows XP, works on Windows 7, doesn't work on Windows 10,
  should work on Windows Vista ... Windows 7,
  contains some string constants overlapping header fields. On 32-bit
  Windows 7 the first 256 bytes would have been enough.
* hh2.golden.exe (402 bytes): Should work on Windows XP ... Windows 10,
  contains some string constants overlapping header fields.
  It doesn't work on Windows NT 3.51 (not even after changing the
  SubsystemVersion to 3.10), and it doesn't work on Windows 95 either.
* hh2d.golden.exe (408 bytes): Should work on Windows 95 ... Windows 10,
  contains some string constants overlapping header fields.
  It doesn't work on Windows NT 3.51 (not even after changing the
  SubsystemVersion to 3.10). It employs a trick so that the entire file is
  loaded to section .text, without having to align to it 512 bytes.
* hh3gf.golden.exe (2048 bytes): Works on Windows NT 3.1 ... Windows 10.
  Built with MinGW GCC from a .c source and has SubsystemVersion 3.10
  for Windows NT 3.1 compatibility.
* hh3tf.golden.exe (1536 bytes): Works on Windows NT 3.1 ... Windows 10.
  Built with TCC 0.9.26 from a .c source, and the SubsystemVersion field in
  the PE header was changed from 4.0 to 3.10 for Windows NT 3.1
  compatibility.
* hh3wf.golden.exe (3072 bytes): Works on Windows NT 3.1 ... Windows 10.
  Built with OpenWatcom V2 owcc from a .c source and has SubsystemVersion 3.10
  for Windows NT 3.1 compatibility.
* hh3tg.golden.exe (3072 bytes): Works on Windows NT 3.1 ... Windows 10 and
  Win32s. It's a GUI application, it uses MessageBox, loading it from
  USER32.DLL with LoadLibraryA. Built with MinGW GCC from a .c
  source and has SubsystemVersion 3.10 for Windows NT 3.1 compatibility.
  In addition to the .c source, a bit-by-bit identical NASM reimplementation
  is also available (hh2tgn.nasm).
* hh3tw.golden.exe (3072 bytes): Works on Windows NT 3.1 ... Windows 10 and
  Win32s. It's a GUI application, it uses MessageBox, loading it from
  USER32.DLL with LoadLibraryA. Built with OpenWatcom V2 owcc from a .c
  source and has SubsystemVersion 3.10 for Windows NT 3.1 compatibility.
* hh4t.golden.exe (610 bytes): It's an optimized NASM reimplementation of
  hh3tg.exe. It could be optimized further by directly importing user32.dll
  (instead of with LoadLibraryA. It works on Windows NT 3.1 ... Windows 10
  and Win32s. It's a GUI application, it uses MessageBox, loading it from
  USER32.DLL with LoadLibraryA.
* hh6a.golden.exe (1536 bytes); Same as hh3tf.golden.exe, but reimplmented
  in NASM.
* hh6b.golden.exe (1024 bytes): Like hh6a.golden.exe, but smaller, because
  the .data section was merged to the .text section. It works on
  Windows NT 3.1--Windows 10, tested on Windows NT 3.1, Windows 95, Windows
  XP and Wine 5.0.
* hh6c.golden.exe (688 bytes): Like hh6b.golden.exe, but contains optimized
  code for the hello-world, and the trailing 0 bytes are stripped.
  .data section was merged to the .text section. It works on
  Windows NT 3.1--Windows 10, tested on Windows NT 3.1, Windows
  95, Windows XP and Wine 5.0.
* hh6d.golden.exe (584 bytes): Like hh6c.golden.exe, but some padding bytes
  and some image data directory entried were removed, and some read-only data
  has been moved from the .text section to the header. It works on
  Windows NT 3.1--Windows 10, tested on Windows NT 3.1, Windows
  95, Windows XP and Wine 5.0. It's not possible to go below 512 bytes,
  because Windows NT 3.1 and Windows 95 don't support section
  alignment lower than 512 or section starting at file offset 0. See
  hh2.golden.exe for the `-2' hack to make it work on Windows XP and Wine.
* hh6r.golden.exe (584 bytes): Like hh6d.golden.exe, but with relocation table.
  It works on the same systems as hh6d.golden.exe, because Win32s doesn't
  support Win32 console programs.
* hh7.golden.exe (584 bytes): Like hh6d.golden.exe, but it uses NASM library
  smallpe.inc.nasm, for convenient creation of small arbitrary (i.e. not
  only hello-world) Win32 PE .exe executables using KERNEL32.DLL only.
* box1.golden.exe (268 bytes): Doesn't work on Windows XP, works on Windows
  7, should work on Windows Vista ... Windows 10,
  contains some string constants overlapping header fields. On 32-bit
  Windows 7 the first 261 bytes would have been enough.
  It is a copy of the coee at
  https://www.codejuggle.dj/creating-the-smallest-possible-windows-executable-using-assembly-language/

How to compile:

* On a Unix system (e.g. Linux) with the `nasm' and `make' tools installed,
  just run `make' (without the quotes) in the directory containing hh2.nasm.
  The minimum NASM version required is 0.98.39.
* Alternatively, on other systems, look at the beginning of the hh6d.nasm
  etc. On Windows, you may have to run `nasmw' instead of `nasm'.

Related projects and docs:

* https://www.codejuggle.dj/creating-the-smallest-possible-windows-executable-using-assembly-language/
  is a related project from 2015, and its tiny .exe is even smaller: 268
  bytes. Unfortunately it doesn't run on Windows XP (``The application
  failed to initialize properly (0xc0000007b). Click on OK to terminate the
  application.''. It works on Wine 1.6.2, Windows 7 32-bit,
  and its author claims that it runs on Windows 7 64-bit. See
  box1.nasm and box1.golden.exe for a copy of the code.
* The 268-byte PE .exe header pattern:
  http://pferrie.host22.com/misc/tiny/pehdr.htm
* 268-byte amd64 tiny PE .exe where every byte is executed:
  https://drakopensulo.wordpress.com/2017/08/06/smallest-pe-executable-x64-with-every-byte-executed/
* A longer, useful writeup on tiny PE .exe:
  http://www.phreedom.org/research/tinype/tiny.import.209/tiny.asm
  The subpage
  http://www.phreedom.org/research/tinype/tiny.import.209/tiny.asm
  contains 209-byte tiny.exe with an import. Windows XP SP3 says:
  ``Program too big to fit in memory''.
* Crinkler-related discussion of tiny PE .exe and the 268-byte minimum:
  http://www.pouet.net/topic.php?which=9565
* Crinkler (http://www.crinkler.net/), a combined linker and compressor to
  generate tiny Win32 PE .exe files. An .exe files generated by Crinkler 2.0
  (aw50cm8_by_knl__ishy.exe)
  didn't work for the author of hh2.nasm on Windows XP SP3 (even though
  the documentation of Crinkler explicitly says that Windows XP is
  supported). Crinkler 2.0 itself didn't work for the author of hh2.nasm on
  Windows XP SP3 (``The application failed to initialize properly
  (0xc0000022). Click OK to terminate the application.''.) Crinkler 2.0
  started up on Wine 1.6.2, but it failed to create an .exe file
  (``Oops! Crinkler has crashed.'', probably because the dbghelp.dll in Wine
  doesn't work.)
* https://code.google.com/archive/p/corkami/wikis/PE.wiki
  contains older documentation about PE.
* https://stackoverflow.com/questions/33247785/compile-windows-executables-with-nasm
  asks how to create Win32 PE .exe files with nasm.
* https://stackoverflow.com/questions/42022132/how-to-create-tiny-pe-win32-executables-using-mingw
  contains a C hello-world Win32 PE .exe, 2048 bytes.

Loader limitations:

* The VirtualAddress of any section must be at least SizeOfHeaders.
* The header is mapped read-only (so no write and no execute). Thus
  IMPORT_ADDRESS_TABLE can't be stored in the header (because that
  requires write access), and the program code (at _start:) can't be
  stored in the header either (because that requires execute access.
* On Windows 95, IMAGE_IMPORT_DESCRIPTORS must not be stored in the
  header. (Wine, Windows NT 3.1 and Windows XP allow it.)
* On Windows NT 3.1 and Windows XP, SectionAlignment must be 4096 (0x1000),
  and FileAlignment must be a power of 2 at least 512.
* On Windows 7, SectionAligment == FileAlignment == 4 can work. See
  hh1.nasm.
* Windows NT 3.1 requires SubsystemVersion=3.10, more recent systems work
  with =3.10 and =4.0 (and possibly others).
* Windows NT 3.1 and Windows 95 look at some PE header fields ignored by
  Windows XP and above, so hh2.nasm works on Windows XP and above, but not
  on Windows NT 3.1 and Windows 95.
* The executable code must be in a section with PointerToRawData larger than 0.
* It's OK that the file size isn't divisile by 0x200 (512), the file can be
  truncated.
* On Windows 7 64-bit, the file size must be at least 268 bytes.
* On Windows NT 3.1 (and possibly others), the PE header (IMAGE_NT_HEADERS)
  must start on a file offset divisible by 4.
* SizeOfOptionalHeader must be >= 0x78.
* SizeOfHeaders must be > 0.
* Windows 95 needs at least 5 entries in IMAGE_DATA_DIRECTORY.
* On Win32s, the PE header (ending with the last byte of the last section
  header) must fit in 0x800 (2048) bytes.
* ReactOS 0.4.14 needs section.VirtualSize >= section.SizeOfRawData.
* ReactOS 0.4.14 is picky about the low 12 bits of SizeOfImage,
  section.VirtualSize and section.SizeOfRawData being too small. Other
  systems seem to round these up to page boundary.
* Imported symbols (i.e. function names) and DLL names don't have to be
  word-aligned, even though OpenWatcom wlink(1) adds NUL bytes for word
  alignment.
* WDOSX doesn't work if there are no relocations (i.e. the
  IMAGE_FILE_RELOCS_STRIPPED must be 0 and the
  IMAGE_DIRECTORY_ENTRY_BASERELOC must be present), thus the minimum
  number of entries in the image directory for WDOSX is 6.
* The minimum number of entries in the image directory if relocations are
  not present is 5 for Windows 95 (and many others).
* The minimum number of entries in the image directory if relocations are
  present is 7 for Windows XP and 6 for many others. 6 is the bare minimum
  because IMAGE_DIRECTORY_ENTRY_BASERELOC must be present.
* There are some others, not mentioned here.

Virus declaration: According to VirusTotal (https://www.virtustotal.com/),
some antivirus software claim that these files (hh*.exe) contain virus or
malware. All of these reports are wrong. These files don't contain any
malware. I wrote each byte of them, and I haven't added any malware.

__END__
