// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ld

/*
 * Derived from:
 * $FreeBSD: src/sys/sys/elf32.h,v 1.8.14.1 2005/12/30 22:13:58 marcel Exp $
 * $FreeBSD: src/sys/sys/elf64.h,v 1.10.14.1 2005/12/30 22:13:58 marcel Exp $
 * $FreeBSD: src/sys/sys/elf_common.h,v 1.15.8.1 2005/12/30 22:13:58 marcel Exp $
 * $FreeBSD: src/sys/alpha/include/elf.h,v 1.14 2003/09/25 01:10:22 peter Exp $
 * $FreeBSD: src/sys/amd64/include/elf.h,v 1.18 2004/08/03 08:21:48 dfr Exp $
 * $FreeBSD: src/sys/arm/include/elf.h,v 1.5.2.1 2006/06/30 21:42:52 cognet Exp $
 * $FreeBSD: src/sys/i386/include/elf.h,v 1.16 2004/08/02 19:12:17 dfr Exp $
 * $FreeBSD: src/sys/powerpc/include/elf.h,v 1.7 2004/11/02 09:47:01 ssouhlal Exp $
 * $FreeBSD: src/sys/sparc64/include/elf.h,v 1.12 2003/09/25 01:10:26 peter Exp $
 *
 * Copyright (c) 1996-1998 John D. Polstra.  All rights reserved.
 * Copyright (c) 2001 David E. O'Brien
 * Portions Copyright 2009 The Go Authors.  All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 * 1. Redistributions of source code must retain the above copyright
 *    notice, this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright
 *    notice, this list of conditions and the following disclaimer in the
 *    documentation and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE AUTHOR AND CONTRIBUTORS ``AS IS'' AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED.  IN NO EVENT SHALL THE AUTHOR OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
 * OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
 * LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
 * OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
 * SUCH DAMAGE.
 *
 */

/*
 * ELF definitions that are independent of architecture or word size.
 */

/*
 * Note header.  The ".note" section contains an array of notes.  Each
 * begins with this header, aligned to a word boundary.  Immediately
 * following the note header is n_namesz bytes of name, padded to the
 * next word boundary.  Then comes n_descsz bytes of descriptor, again
 * padded to a word boundary.  The values of n_namesz and n_descsz do
 * not include the padding.
 */
type Elf_Note struct {
	n_namesz uint32
	n_descsz uint32
	n_type   uint32
}

const (
	EI_MAG0              = 0
	EI_MAG1              = 1
	EI_MAG2              = 2
	EI_MAG3              = 3
	EI_CLASS             = 4
	EI_DATA              = 5
	EI_VERSION           = 6
	EI_OSABI             = 7
	EI_ABIVERSION        = 8
	OLD_EI_BRAND         = 8
	EI_PAD               = 9
	EI_NIDENT            = 16
	ELFMAG0              = 0x7f
	ELFMAG1              = 'E'
	ELFMAG2              = 'L'
	ELFMAG3              = 'F'
	SELFMAG              = 4
	EV_NONE              = 0
	EV_CURRENT           = 1
	ELFCLASSNONE         = 0
	ELFCLASS32           = 1
	ELFCLASS64           = 2
	ELFDATANONE          = 0
	ELFDATA2LSB          = 1
	ELFDATA2MSB          = 2
	ELFOSABI_NONE        = 0
	ELFOSABI_HPUX        = 1
	ELFOSABI_NETBSD      = 2
	ELFOSABI_LINUX       = 3
	ELFOSABI_HURD        = 4
	ELFOSABI_86OPEN      = 5
	ELFOSABI_SOLARIS     = 6
	ELFOSABI_AIX         = 7
	ELFOSABI_IRIX        = 8
	ELFOSABI_FREEBSD     = 9
	ELFOSABI_TRU64       = 10
	ELFOSABI_MODESTO     = 11
	ELFOSABI_OPENBSD     = 12
	ELFOSABI_OPENVMS     = 13
	ELFOSABI_NSK         = 14
	ELFOSABI_ARM         = 97
	ELFOSABI_STANDALONE  = 255
	ELFOSABI_SYSV        = ELFOSABI_NONE
	ELFOSABI_MONTEREY    = ELFOSABI_AIX
	ET_NONE              = 0
	ET_REL               = 1
	ET_EXEC              = 2
	ET_DYN               = 3
	ET_CORE              = 4
	ET_LOOS              = 0xfe00
	ET_HIOS              = 0xfeff
	ET_LOPROC            = 0xff00
	ET_HIPROC            = 0xffff
	EM_NONE              = 0
	EM_M32               = 1
	EM_SPARC             = 2
	EM_386               = 3
	EM_68K               = 4
	EM_88K               = 5
	EM_860               = 7
	EM_MIPS              = 8
	EM_S370              = 9
	EM_MIPS_RS3_LE       = 10
	EM_PARISC            = 15
	EM_VPP500            = 17
	EM_SPARC32PLUS       = 18
	EM_960               = 19
	EM_PPC               = 20
	EM_PPC64             = 21
	EM_S390              = 22
	EM_V800              = 36
	EM_FR20              = 37
	EM_RH32              = 38
	EM_RCE               = 39
	EM_ARM               = 40
	EM_SH                = 42
	EM_SPARCV9           = 43
	EM_TRICORE           = 44
	EM_ARC               = 45
	EM_H8_300            = 46
	EM_H8_300H           = 47
	EM_H8S               = 48
	EM_H8_500            = 49
	EM_IA_64             = 50
	EM_MIPS_X            = 51
	EM_COLDFIRE          = 52
	EM_68HC12            = 53
	EM_MMA               = 54
	EM_PCP               = 55
	EM_NCPU              = 56
	EM_NDR1              = 57
	EM_STARCORE          = 58
	EM_ME16              = 59
	EM_ST100             = 60
	EM_TINYJ             = 61
	EM_X86_64            = 62
	EM_AARCH64           = 183
	EM_486               = 6
	EM_MIPS_RS4_BE       = 10
	EM_ALPHA_STD         = 41
	EM_ALPHA             = 0x9026
	SHN_UNDEF            = 0
	SHN_LORESERVE        = 0xff00
	SHN_LOPROC           = 0xff00
	SHN_HIPROC           = 0xff1f
	SHN_LOOS             = 0xff20
	SHN_HIOS             = 0xff3f
	SHN_ABS              = 0xfff1
	SHN_COMMON           = 0xfff2
	SHN_XINDEX           = 0xffff
	SHN_HIRESERVE        = 0xffff
	SHT_NULL             = 0
	SHT_PROGBITS         = 1
	SHT_SYMTAB           = 2
	SHT_STRTAB           = 3
	SHT_RELA             = 4
	SHT_HASH             = 5
	SHT_DYNAMIC          = 6
	SHT_NOTE             = 7
	SHT_NOBITS           = 8
	SHT_REL              = 9
	SHT_SHLIB            = 10
	SHT_DYNSYM           = 11
	SHT_INIT_ARRAY       = 14
	SHT_FINI_ARRAY       = 15
	SHT_PREINIT_ARRAY    = 16
	SHT_GROUP            = 17
	SHT_SYMTAB_SHNDX     = 18
	SHT_LOOS             = 0x60000000
	SHT_HIOS             = 0x6fffffff
	SHT_GNU_VERDEF       = 0x6ffffffd
	SHT_GNU_VERNEED      = 0x6ffffffe
	SHT_GNU_VERSYM       = 0x6fffffff
	SHT_LOPROC           = 0x70000000
	SHT_HIPROC           = 0x7fffffff
	SHT_LOUSER           = 0x80000000
	SHT_HIUSER           = 0xffffffff
	SHF_WRITE            = 0x1
	SHF_ALLOC            = 0x2
	SHF_EXECINSTR        = 0x4
	SHF_MERGE            = 0x10
	SHF_STRINGS          = 0x20
	SHF_INFO_LINK        = 0x40
	SHF_LINK_ORDER       = 0x80
	SHF_OS_NONCONFORMING = 0x100
	SHF_GROUP            = 0x200
	SHF_TLS              = 0x400
	SHF_MASKOS           = 0x0ff00000
	SHF_MASKPROC         = 0xf0000000
	PT_NULL              = 0
	PT_LOAD              = 1
	PT_DYNAMIC           = 2
	PT_INTERP            = 3
	PT_NOTE              = 4
	PT_SHLIB             = 5
	PT_PHDR              = 6
	PT_TLS               = 7
	PT_LOOS              = 0x60000000
	PT_HIOS              = 0x6fffffff
	PT_LOPROC            = 0x70000000
	PT_HIPROC            = 0x7fffffff
	PT_GNU_STACK         = 0x6474e551
	PT_PAX_FLAGS         = 0x65041580
	PF_X                 = 0x1
	PF_W                 = 0x2
	PF_R                 = 0x4
	PF_MASKOS            = 0x0ff00000
	PF_MASKPROC          = 0xf0000000
	DT_NULL              = 0
	DT_NEEDED            = 1
	DT_PLTRELSZ          = 2
	DT_PLTGOT            = 3
	DT_HASH              = 4
	DT_STRTAB            = 5
	DT_SYMTAB            = 6
	DT_RELA              = 7
	DT_RELASZ            = 8
	DT_RELAENT           = 9
	DT_STRSZ             = 10
	DT_SYMENT            = 11
	DT_INIT              = 12
	DT_FINI              = 13
	DT_SONAME            = 14
	DT_RPATH             = 15
	DT_SYMBOLIC          = 16
	DT_REL               = 17
	DT_RELSZ             = 18
	DT_RELENT            = 19
	DT_PLTREL            = 20
	DT_DEBUG             = 21
	DT_TEXTREL           = 22
	DT_JMPREL            = 23
	DT_BIND_NOW          = 24
	DT_INIT_ARRAY        = 25
	DT_FINI_ARRAY        = 26
	DT_INIT_ARRAYSZ      = 27
	DT_FINI_ARRAYSZ      = 28
	DT_RUNPATH           = 29
	DT_FLAGS             = 30
	DT_ENCODING          = 32
	DT_PREINIT_ARRAY     = 32
	DT_PREINIT_ARRAYSZ   = 33
	DT_LOOS              = 0x6000000d
	DT_HIOS              = 0x6ffff000
	DT_LOPROC            = 0x70000000
	DT_HIPROC            = 0x7fffffff
	DT_VERNEED           = 0x6ffffffe
	DT_VERNEEDNUM        = 0x6fffffff
	DT_VERSYM            = 0x6ffffff0
	DT_PPC64_GLINK       = DT_LOPROC + 0
	DT_PPC64_OPT         = DT_LOPROC + 3
	DF_ORIGIN            = 0x0001
	DF_SYMBOLIC          = 0x0002
	DF_TEXTREL           = 0x0004
	DF_BIND_NOW          = 0x0008
	DF_STATIC_TLS        = 0x0010
	NT_PRSTATUS          = 1
	NT_FPREGSET          = 2
	NT_PRPSINFO          = 3
	STB_LOCAL            = 0
	STB_GLOBAL           = 1
	STB_WEAK             = 2
	STB_LOOS             = 10
	STB_HIOS             = 12
	STB_LOPROC           = 13
	STB_HIPROC           = 15
	STT_NOTYPE           = 0
	STT_OBJECT           = 1
	STT_FUNC             = 2
	STT_SECTION          = 3
	STT_FILE             = 4
	STT_COMMON           = 5
	STT_TLS              = 6
	STT_LOOS             = 10
	STT_HIOS             = 12
	STT_LOPROC           = 13
	STT_HIPROC           = 15
	STV_DEFAULT          = 0x0
	STV_INTERNAL         = 0x1
	STV_HIDDEN           = 0x2
	STV_PROTECTED        = 0x3
	STN_UNDEF            = 0
)

/* For accessing the fields of r_info. */

/* For constructing r_info from field values. */

/*
 * Relocation types.
 */
const (
	R_X86_64_NONE              = 0
	R_X86_64_64                = 1
	R_X86_64_PC32              = 2
	R_X86_64_GOT32             = 3
	R_X86_64_PLT32             = 4
	R_X86_64_COPY              = 5
	R_X86_64_GLOB_DAT          = 6
	R_X86_64_JMP_SLOT          = 7
	R_X86_64_RELATIVE          = 8
	R_X86_64_GOTPCREL          = 9
	R_X86_64_32                = 10
	R_X86_64_32S               = 11
	R_X86_64_16                = 12
	R_X86_64_PC16              = 13
	R_X86_64_8                 = 14
	R_X86_64_PC8               = 15
	R_X86_64_DTPMOD64          = 16
	R_X86_64_DTPOFF64          = 17
	R_X86_64_TPOFF64           = 18
	R_X86_64_TLSGD             = 19
	R_X86_64_TLSLD             = 20
	R_X86_64_DTPOFF32          = 21
	R_X86_64_GOTTPOFF          = 22
	R_X86_64_TPOFF32           = 23
	R_X86_64_COUNT             = 24
	R_AARCH64_ABS64            = 257
	R_AARCH64_ABS32            = 258
	R_AARCH64_CALL26           = 283
	R_AARCH64_ADR_PREL_PG_HI21 = 275
	R_AARCH64_ADD_ABS_LO12_NC  = 277
	R_ALPHA_NONE               = 0
	R_ALPHA_REFLONG            = 1
	R_ALPHA_REFQUAD            = 2
	R_ALPHA_GPREL32            = 3
	R_ALPHA_LITERAL            = 4
	R_ALPHA_LITUSE             = 5
	R_ALPHA_GPDISP             = 6
	R_ALPHA_BRADDR             = 7
	R_ALPHA_HINT               = 8
	R_ALPHA_SREL16             = 9
	R_ALPHA_SREL32             = 10
	R_ALPHA_SREL64             = 11
	R_ALPHA_OP_PUSH            = 12
	R_ALPHA_OP_STORE           = 13
	R_ALPHA_OP_PSUB            = 14
	R_ALPHA_OP_PRSHIFT         = 15
	R_ALPHA_GPVALUE            = 16
	R_ALPHA_GPRELHIGH          = 17
	R_ALPHA_GPRELLOW           = 18
	R_ALPHA_IMMED_GP_16        = 19
	R_ALPHA_IMMED_GP_HI32      = 20
	R_ALPHA_IMMED_SCN_HI32     = 21
	R_ALPHA_IMMED_BR_HI32      = 22
	R_ALPHA_IMMED_LO32         = 23
	R_ALPHA_COPY               = 24
	R_ALPHA_GLOB_DAT           = 25
	R_ALPHA_JMP_SLOT           = 26
	R_ALPHA_RELATIVE           = 27
	R_ALPHA_COUNT              = 28
	R_ARM_NONE                 = 0
	R_ARM_PC24                 = 1
	R_ARM_ABS32                = 2
	R_ARM_REL32                = 3
	R_ARM_PC13                 = 4
	R_ARM_ABS16                = 5
	R_ARM_ABS12                = 6
	R_ARM_THM_ABS5             = 7
	R_ARM_ABS8                 = 8
	R_ARM_SBREL32              = 9
	R_ARM_THM_PC22             = 10
	R_ARM_THM_PC8              = 11
	R_ARM_AMP_VCALL9           = 12
	R_ARM_SWI24                = 13
	R_ARM_THM_SWI8             = 14
	R_ARM_XPC25                = 15
	R_ARM_THM_XPC22            = 16
	R_ARM_COPY                 = 20
	R_ARM_GLOB_DAT             = 21
	R_ARM_JUMP_SLOT            = 22
	R_ARM_RELATIVE             = 23
	R_ARM_GOTOFF               = 24
	R_ARM_GOTPC                = 25
	R_ARM_GOT32                = 26
	R_ARM_PLT32                = 27
	R_ARM_CALL                 = 28
	R_ARM_JUMP24               = 29
	R_ARM_V4BX                 = 40
	R_ARM_GOT_PREL             = 96
	R_ARM_GNU_VTENTRY          = 100
	R_ARM_GNU_VTINHERIT        = 101
	R_ARM_TLS_IE32             = 107
	R_ARM_TLS_LE32             = 108
	R_ARM_RSBREL32             = 250
	R_ARM_THM_RPC22            = 251
	R_ARM_RREL32               = 252
	R_ARM_RABS32               = 253
	R_ARM_RPC24                = 254
	R_ARM_RBASE                = 255
	R_ARM_COUNT                = 38
	R_386_NONE                 = 0
	R_386_32                   = 1
	R_386_PC32                 = 2
	R_386_GOT32                = 3
	R_386_PLT32                = 4
	R_386_COPY                 = 5
	R_386_GLOB_DAT             = 6
	R_386_JMP_SLOT             = 7
	R_386_RELATIVE             = 8
	R_386_GOTOFF               = 9
	R_386_GOTPC                = 10
	R_386_TLS_TPOFF            = 14
	R_386_TLS_IE               = 15
	R_386_TLS_GOTIE            = 16
	R_386_TLS_LE               = 17
	R_386_TLS_GD               = 18
	R_386_TLS_LDM              = 19
	R_386_TLS_GD_32            = 24
	R_386_TLS_GD_PUSH          = 25
	R_386_TLS_GD_CALL          = 26
	R_386_TLS_GD_POP           = 27
	R_386_TLS_LDM_32           = 28
	R_386_TLS_LDM_PUSH         = 29
	R_386_TLS_LDM_CALL         = 30
	R_386_TLS_LDM_POP          = 31
	R_386_TLS_LDO_32           = 32
	R_386_TLS_IE_32            = 33
	R_386_TLS_LE_32            = 34
	R_386_TLS_DTPMOD32         = 35
	R_386_TLS_DTPOFF32         = 36
	R_386_TLS_TPOFF32          = 37
	R_386_COUNT                = 38
	R_PPC_NONE                 = 0
	R_PPC_ADDR32               = 1
	R_PPC_ADDR24               = 2
	R_PPC_ADDR16               = 3
	R_PPC_ADDR16_LO            = 4
	R_PPC_ADDR16_HI            = 5
	R_PPC_ADDR16_HA            = 6
	R_PPC_ADDR14               = 7
	R_PPC_ADDR14_BRTAKEN       = 8
	R_PPC_ADDR14_BRNTAKEN      = 9
	R_PPC_REL24                = 10
	R_PPC_REL14                = 11
	R_PPC_REL14_BRTAKEN        = 12
	R_PPC_REL14_BRNTAKEN       = 13
	R_PPC_GOT16                = 14
	R_PPC_GOT16_LO             = 15
	R_PPC_GOT16_HI             = 16
	R_PPC_GOT16_HA             = 17
	R_PPC_PLTREL24             = 18
	R_PPC_COPY                 = 19
	R_PPC_GLOB_DAT             = 20
	R_PPC_JMP_SLOT             = 21
	R_PPC_RELATIVE             = 22
	R_PPC_LOCAL24PC            = 23
	R_PPC_UADDR32              = 24
	R_PPC_UADDR16              = 25
	R_PPC_REL32                = 26
	R_PPC_PLT32                = 27
	R_PPC_PLTREL32             = 28
	R_PPC_PLT16_LO             = 29
	R_PPC_PLT16_HI             = 30
	R_PPC_PLT16_HA             = 31
	R_PPC_SDAREL16             = 32
	R_PPC_SECTOFF              = 33
	R_PPC_SECTOFF_LO           = 34
	R_PPC_SECTOFF_HI           = 35
	R_PPC_SECTOFF_HA           = 36
	R_PPC_COUNT                = 37
	R_PPC_TLS                  = 67
	R_PPC_DTPMOD32             = 68
	R_PPC_TPREL16              = 69
	R_PPC_TPREL16_LO           = 70
	R_PPC_TPREL16_HI           = 71
	R_PPC_TPREL16_HA           = 72
	R_PPC_TPREL32              = 73
	R_PPC_DTPREL16             = 74
	R_PPC_DTPREL16_LO          = 75
	R_PPC_DTPREL16_HI          = 76
	R_PPC_DTPREL16_HA          = 77
	R_PPC_DTPREL32             = 78
	R_PPC_GOT_TLSGD16          = 79
	R_PPC_GOT_TLSGD16_LO       = 80
	R_PPC_GOT_TLSGD16_HI       = 81
	R_PPC_GOT_TLSGD16_HA       = 82
	R_PPC_GOT_TLSLD16          = 83
	R_PPC_GOT_TLSLD16_LO       = 84
	R_PPC_GOT_TLSLD16_HI       = 85
	R_PPC_GOT_TLSLD16_HA       = 86
	R_PPC_GOT_TPREL16          = 87
	R_PPC_GOT_TPREL16_LO       = 88
	R_PPC_GOT_TPREL16_HI       = 89
	R_PPC_GOT_TPREL16_HA       = 90
	R_PPC_EMB_NADDR32          = 101
	R_PPC_EMB_NADDR16          = 102
	R_PPC_EMB_NADDR16_LO       = 103
	R_PPC_EMB_NADDR16_HI       = 104
	R_PPC_EMB_NADDR16_HA       = 105
	R_PPC_EMB_SDAI16           = 106
	R_PPC_EMB_SDA2I16          = 107
	R_PPC_EMB_SDA2REL          = 108
	R_PPC_EMB_SDA21            = 109
	R_PPC_EMB_MRKREF           = 110
	R_PPC_EMB_RELSEC16         = 111
	R_PPC_EMB_RELST_LO         = 112
	R_PPC_EMB_RELST_HI         = 113
	R_PPC_EMB_RELST_HA         = 114
	R_PPC_EMB_BIT_FLD          = 115
	R_PPC_EMB_RELSDA           = 116
	R_PPC_EMB_COUNT            = R_PPC_EMB_RELSDA - R_PPC_EMB_NADDR32 + 1
	R_PPC64_REL24              = R_PPC_REL24
	R_PPC64_JMP_SLOT           = R_PPC_JMP_SLOT
	R_PPC64_ADDR64             = 38
	R_PPC64_TOC16              = 47
	R_PPC64_TOC16_LO           = 48
	R_PPC64_TOC16_HI           = 49
	R_PPC64_TOC16_HA           = 50
	R_PPC64_TOC16_DS           = 63
	R_PPC64_TOC16_LO_DS        = 64
	R_PPC64_REL16_LO           = 250
	R_PPC64_REL16_HI           = 251
	R_PPC64_REL16_HA           = 252
	R_SPARC_NONE               = 0
	R_SPARC_8                  = 1
	R_SPARC_16                 = 2
	R_SPARC_32                 = 3
	R_SPARC_DISP8              = 4
	R_SPARC_DISP16             = 5
	R_SPARC_DISP32             = 6
	R_SPARC_WDISP30            = 7
	R_SPARC_WDISP22            = 8
	R_SPARC_HI22               = 9
	R_SPARC_22                 = 10
	R_SPARC_13                 = 11
	R_SPARC_LO10               = 12
	R_SPARC_GOT10              = 13
	R_SPARC_GOT13              = 14
	R_SPARC_GOT22              = 15
	R_SPARC_PC10               = 16
	R_SPARC_PC22               = 17
	R_SPARC_WPLT30             = 18
	R_SPARC_COPY               = 19
	R_SPARC_GLOB_DAT           = 20
	R_SPARC_JMP_SLOT           = 21
	R_SPARC_RELATIVE           = 22
	R_SPARC_UA32               = 23
	R_SPARC_PLT32              = 24
	R_SPARC_HIPLT22            = 25
	R_SPARC_LOPLT10            = 26
	R_SPARC_PCPLT32            = 27
	R_SPARC_PCPLT22            = 28
	R_SPARC_PCPLT10            = 29
	R_SPARC_10                 = 30
	R_SPARC_11                 = 31
	R_SPARC_64                 = 32
	R_SPARC_OLO10              = 33
	R_SPARC_HH22               = 34
	R_SPARC_HM10               = 35
	R_SPARC_LM22               = 36
	R_SPARC_PC_HH22            = 37
	R_SPARC_PC_HM10            = 38
	R_SPARC_PC_LM22            = 39
	R_SPARC_WDISP16            = 40
	R_SPARC_WDISP19            = 41
	R_SPARC_GLOB_JMP           = 42
	R_SPARC_7                  = 43
	R_SPARC_5                  = 44
	R_SPARC_6                  = 45
	R_SPARC_DISP64             = 46
	R_SPARC_PLT64              = 47
	R_SPARC_HIX22              = 48
	R_SPARC_LOX10              = 49
	R_SPARC_H44                = 50
	R_SPARC_M44                = 51
	R_SPARC_L44                = 52
	R_SPARC_REGISTER           = 53
	R_SPARC_UA64               = 54
	R_SPARC_UA16               = 55
	ARM_MAGIC_TRAMP_NUMBER     = 0x5c000003
)

/*
 * Symbol table entries.
 */

/* For accessing the fields of st_info. */

/* For constructing st_info from field values. */

/* For accessing the fields of st_other. */

/*
 * ELF header.
 */
type ElfEhdr struct {
	ident     [EI_NIDENT]uint8
	type_     uint16
	machine   uint16
	version   uint32
	entry     uint64
	phoff     uint64
	shoff     uint64
	flags     uint32
	ehsize    uint16
	phentsize uint16
	phnum     uint16
	shentsize uint16
	shnum     uint16
	shstrndx  uint16
}

/*
 * Section header.
 */
type ElfShdr struct {
	name      uint32
	type_     uint32
	flags     uint64
	addr      uint64
	off       uint64
	size      uint64
	link      uint32
	info      uint32
	addralign uint64
	entsize   uint64
	shnum     int
	secsym    *LSym
}

/*
 * Program header.
 */
type ElfPhdr struct {
	type_  uint32
	flags  uint32
	off    uint64
	vaddr  uint64
	paddr  uint64
	filesz uint64
	memsz  uint64
	align  uint64
}

/* For accessing the fields of r_info. */

/* For constructing r_info from field values. */

/*
 * Symbol table entries.
 */

/* For accessing the fields of st_info. */

/* For constructing st_info from field values. */

/* For accessing the fields of st_other. */

/*
 * Go linker interface
 */
const (
	ELF64HDRSIZE  = 64
	ELF64PHDRSIZE = 56
	ELF64SHDRSIZE = 64
	ELF64RELSIZE  = 16
	ELF64RELASIZE = 24
	ELF64SYMSIZE  = 24
	ELF32HDRSIZE  = 52
	ELF32PHDRSIZE = 32
	ELF32SHDRSIZE = 40
	ELF32SYMSIZE  = 16
	ELF32RELSIZE  = 8
)

/*
 * The interface uses the 64-bit structures always,
 * to avoid code duplication.  The writers know how to
 * marshal a 32-bit representation from the 64-bit structure.
 */

var Elfstrdat []byte

/*
 * Total amount of space to reserve at the start of the file
 * for Header, PHeaders, SHeaders, and interp.
 * May waste some.
 * On FreeBSD, cannot be larger than a page.
 */
const (
	ELFRESERVE = 4096
)

/*
 * We use the 64-bit data structures on both 32- and 64-bit machines
 * in order to write the code just once.  The 64-bit data structure is
 * written in the 32-bit format on the 32-bit machines.
 */
const (
	NSECT = 48
)

var Iself bool

var Nelfsym int = 1

var elf64 bool

var ehdr ElfEhdr

var phdr [NSECT]*ElfPhdr

var shdr [NSECT]*ElfShdr

var interp string

type Elfstring struct {
	s   string
	off int
}

var elfstr [100]Elfstring

var nelfstr int

var buildinfo []byte

func getElfEhdr() *ElfEhdr {
	return &ehdr
}

/* Taken directly from the definition document for ELF64 */
func elfhash(name []byte) uint32 {
	var h uint32 = 0
	var g uint32
	for len(name) != 0 {
		h = (h << 4) + uint32(name[0])
		name = name[1:]
		g = h & 0xf0000000
		if g != 0 {
			h ^= g >> 24
		}
		h &= 0x0fffffff
	}

	return h
}

func elfinterp(sh *ElfShdr, startva uint64, resoff uint64, p string) int {
	interp = p
	n := len(interp) + 1
	sh.addr = startva + resoff - uint64(n)
	sh.off = resoff - uint64(n)
	sh.size = uint64(n)

	return n
}

func elfnote(sh *ElfShdr, startva uint64, resoff uint64, sz int, alloc bool) int {
	n := 3*4 + uint64(sz) + resoff%4

	sh.type_ = SHT_NOTE
	if alloc {
		sh.flags = SHF_ALLOC
	}
	sh.addralign = 4
	sh.addr = startva + resoff - n
	sh.off = resoff - n
	sh.size = n - resoff%4

	return int(n)
}

// NetBSD Signature (as per sys/exec_elf.h)
const (
	ELF_NOTE_NETBSD_NAMESZ  = 7
	ELF_NOTE_NETBSD_DESCSZ  = 4
	ELF_NOTE_NETBSD_TAG     = 1
	ELF_NOTE_NETBSD_VERSION = 599000000 /* NetBSD 5.99 */
)

var ELF_NOTE_NETBSD_NAME = []byte("NetBSD\x00")

func elfnetbsdsig(sh *ElfShdr, startva uint64, resoff uint64) int {
	n := int(Rnd(ELF_NOTE_NETBSD_NAMESZ, 4) + Rnd(ELF_NOTE_NETBSD_DESCSZ, 4))
	return elfnote(sh, startva, resoff, n, true)
}

// OpenBSD Signature
const (
	ELF_NOTE_OPENBSD_NAMESZ  = 8
	ELF_NOTE_OPENBSD_DESCSZ  = 4
	ELF_NOTE_OPENBSD_TAG     = 1
	ELF_NOTE_OPENBSD_VERSION = 0
)

var ELF_NOTE_OPENBSD_NAME = []byte("OpenBSD\x00")

func elfopenbsdsig(sh *ElfShdr, startva uint64, resoff uint64) int {
	n := ELF_NOTE_OPENBSD_NAMESZ + ELF_NOTE_OPENBSD_DESCSZ
	return elfnote(sh, startva, resoff, n, true)
}

// Build info note
const (
	ELF_NOTE_BUILDINFO_NAMESZ = 4
	ELF_NOTE_BUILDINFO_TAG    = 3
)

var ELF_NOTE_BUILDINFO_NAME = []byte("GNU\x00")

func elfbuildinfo(sh *ElfShdr, startva uint64, resoff uint64) int {
	n := int(ELF_NOTE_BUILDINFO_NAMESZ + Rnd(int64(len(buildinfo)), 4))
	return elfnote(sh, startva, resoff, n, true)
}

// Go specific notes
const (
	ELF_NOTE_GOPKGLIST_TAG = 1
	ELF_NOTE_GOABIHASH_TAG = 2
	ELF_NOTE_GODEPS_TAG    = 3
	ELF_NOTE_GOBUILDID_TAG = 4
)

var ELF_NOTE_GO_NAME = []byte("Go\x00\x00")

var elfverneed int

type Elfaux struct {
	next *Elfaux
	num  int
	vers string
}

type Elflib struct {
	next *Elflib
	aux  *Elfaux
	file string
}

func addelflib(list **Elflib, file string, vers string) *Elfaux {
	var lib *Elflib

	for lib = *list; lib != nil; lib = lib.next {
		if lib.file == file {
			goto havelib
		}
	}
	lib = new(Elflib)
	lib.next = *list
	lib.file = file
	*list = lib

havelib:
	for aux := lib.aux; aux != nil; aux = aux.next {
		if aux.vers == vers {
			return aux
		}
	}
	aux := new(Elfaux)
	aux.next = lib.aux
	aux.vers = vers
	lib.aux = aux

	return aux
}

func phsh(ph *ElfPhdr, sh *ElfShdr) {
	ph.vaddr = sh.addr
	ph.paddr = ph.vaddr
	ph.off = sh.off
	ph.filesz = sh.size
	ph.memsz = sh.size
	ph.align = sh.addralign
}

func ELF32_R_SYM(info uint32) uint32 {
	return info >> 8
}

func ELF32_R_TYPE(info uint32) uint32 {
	return uint32(uint8(info))
}

func ELF32_R_INFO(sym uint32, type_ uint32) uint32 {
	return sym<<8 | type_
}

func ELF32_ST_BIND(info uint8) uint8 {
	return info >> 4
}

func ELF32_ST_TYPE(info uint8) uint8 {
	return info & 0xf
}

func ELF32_ST_INFO(bind uint8, type_ uint8) uint8 {
	return bind<<4 | type_&0xf
}

func ELF32_ST_VISIBILITY(oth uint8) uint8 {
	return oth & 3
}

func ELF64_R_SYM(info uint64) uint32 {
	return uint32(info >> 32)
}

func ELF64_R_TYPE(info uint64) uint32 {
	return uint32(info)
}

func ELF64_R_INFO(sym uint32, type_ uint32) uint64 {
	return uint64(sym)<<32 | uint64(type_)
}

func ELF64_ST_BIND(info uint8) uint8 {
	return info >> 4
}

func ELF64_ST_TYPE(info uint8) uint8 {
	return info & 0xf
}

func ELF64_ST_INFO(bind uint8, type_ uint8) uint8 {
	return bind<<4 | type_&0xf
}

func ELF64_ST_VISIBILITY(oth uint8) uint8 {
	return oth & 3
}
