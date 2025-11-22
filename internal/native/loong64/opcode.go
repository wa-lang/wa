// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package loong64

// TODO: 指令格式注释

// 指令格式类型9种
type OpFormatType int

const (
	_ OpFormatType = iota + 1
	OpFormatType_2R
	OpFormatType_3R
	OpFormatType_4R
	OpFormatType_2RI8
	OpFormatType_2RI12
	OpFormatType_2RI14
	OpFormatType_2RI16
	OpFormatType_1RI21
	OpFormatType_I26
)
