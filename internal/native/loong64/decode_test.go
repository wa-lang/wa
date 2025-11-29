// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package loong64

import (
	"encoding/binary"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func _TestDecode(t *testing.T) {
	input := filepath.Join("testdata", "plan9cases.txt")
	data, err := os.ReadFile(input)
	if err != nil {
		t.Fatal(err)
	}
	all := string(data)
	for strings.Contains(all, "\t\t") {
		all = strings.Replace(all, "\t\t", "\t", -1)
	}
	for _, line := range strings.Split(all, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		f := strings.SplitN(line, "\t", 2)
		i := strings.Index(f[0], "|")
		if i < 0 {
			t.Errorf("parsing %q: missing | separator", f[0])
			continue
		}
		if i%2 != 0 {
			t.Errorf("parsing %q: misaligned | separator", f[0])
		}
		code, err := hex.DecodeString(f[0][:i] + f[0][i+1:])
		if err != nil {
			t.Errorf("parsing %q: %v", f[0], err)
			continue
		}

		x := binary.LittleEndian.Uint32(code)
		as, arg, err := Decode(x)
		if err != nil {
			t.Errorf("parsing %x: %s", code, err)
			continue
		}

		out := AsmSyntax(as, "", arg)
		if asm := f[1]; asm != out || len(asm) != len(out) {
			t.Errorf("Decode(%s) = %s want %s", f[0], out, asm)
		}
	}
}
