// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

//go:build !inline_asm

package main

func I64Eqz(a int32) int32     { return 0 }
func I64Eq(a, b int32) int32   { return 0 }
func I64Ne(a, b int32) int32   { return 0 }
func I64Lt_s(a, b int32) int32 { return 0 }
func I64Lt_u(a, b int32) int32 { return 0 }
func I64Gt_s(a, b int32) int32 { return 0 }
func I64Gt_u(a, b int32) int32 { return 0 }
func I64Le_s(a, b int32) int32 { return 0 }
func I64Le_u(a, b int32) int32 { return 0 }
func I64Ge_s(a, b int32) int32 { return 0 }
func I64Ge_u(a, b int32) int32 { return 0 }
