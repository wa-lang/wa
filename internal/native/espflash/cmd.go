// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package espflash

type CmdOpcode byte

const (
	CmdFlashBegin      CmdOpcode = 0x02
	CmdFlashData       CmdOpcode = 0x03
	CmdFlashEnd        CmdOpcode = 0x04
	CmdMemBegin        CmdOpcode = 0x05
	CmdMemEnd          CmdOpcode = 0x06
	CmdMemData         CmdOpcode = 0x07
	CmdSync            CmdOpcode = 0x08
	CmdWriteReg        CmdOpcode = 0x09
	CmdReadReg         CmdOpcode = 0x0A
	CmdSPISetParams    CmdOpcode = 0x0B
	CmdSPIAttach       CmdOpcode = 0x0D
	CmdChangeBaudrate  CmdOpcode = 0x0F
	CmdFlashDeflBegin  CmdOpcode = 0x10
	CmdFlashDeflData   CmdOpcode = 0x11
	CmdFlashDeflEnd    CmdOpcode = 0x12
	CmdSPIFlashMD5     CmdOpcode = 0x13
	CmdGetSecurityInfo CmdOpcode = 0x14

	// Stub only
	CmdEraseFlash  CmdOpcode = 0xd0
	CmdEraseRegion CmdOpcode = 0xd1
	CmdReadFlash   CmdOpcode = 0xd2
	CmdRunUserCode CmdOpcode = 0xd3
)
