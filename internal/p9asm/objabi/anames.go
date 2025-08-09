// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package objabi

import "fmt"

// 指令的名字
// 中间有空隙忽略
var Anames = []string{
	AXXX:      "XXX",
	ACALL:     "CALL",
	ACHECKNIL: "CHECKNIL",
	ADATA:     "DATA",
	ADUFFCOPY: "DUFFCOPY",
	ADUFFZERO: "DUFFZERO",
	AEND:      "END",
	AFUNCDATA: "FUNCDATA",
	AGLOBL:    "GLOBL",
	AJMP:      "JMP",
	ANOP:      "NOP",
	APCDATA:   "PCDATA",
	ARET:      "RET",
	ATEXT:     "TEXT",
	ATYPE:     "TYPE",
	AUNDEF:    "UNDEF",
	AUSEFIELD: "USEFIELD",
	AVARDEF:   "VARDEF",
	AVARKILL:  "VARKILL",
}

func (h *HeadType) Set(name string) error {
	switch name {
	case "darwin", "ios":
		*h = Hdarwin
	case "elf":
		*h = Helf
	case "linux", "android":
		*h = Hlinux
	case "windows":
		*h = Hwindows
	default:
		return fmt.Errorf("invalid headtype: %q", name)
	}
	return nil
}

func (v HeadType) String() string {
	switch v {
	case Hdarwin:
		return "darwin"
	case Helf:
		return "elf"
	case Hlinux:
		return "linux"
	case Hwindows:
		return "windows"
	}
	return fmt.Sprintf("HeadType(%d)", int(v))
}
