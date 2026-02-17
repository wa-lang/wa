// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package pe

// PE文件格式
//
// +---------------------+---+
// |     DOS Header      |   |
// +---------------------+    > MS-DOS 2.0 Section for
// | MS-DOS Stub Program |   |  MS-DOC compatibility only
// +---------------------+---+
// |    PE Signature     |   |
// +---------------------+   |
// |  Image File Header  |    > PE Header (NT Headers)
// |    (COFF Header)    |   |
// +---------------------+   |
// |Image Optional Header|   |
// +---------------------+---+
// |       .text         |   |
// +---------------------+   |
// |      .rdata         |   |
// +---------------------+   |
// |      .idata         |    > Section Table
// +---------------------+   |
// |       .rsrc         |   |
// +---------------------+   |
// |      .reloc         |   |
// +---------------------+---+
// |        ...          |   |
// +---------------------+   |
// |        ...          |    > Other Sextions
// +---------------------+   |
// |        ...          |   |
// +---------------------+---+
//
