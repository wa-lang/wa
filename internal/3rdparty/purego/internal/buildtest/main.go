// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Ebitengine Authors

//go:build !windows

package main

import (
	_ "wa-lang.org/wa/internal/3rdparty/purego"
)

import "C"

// This file tests that build Cgo and purego at the same time succeeds to build (#189).
func main() {
}
