// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package version

import (
	"bytes"
	"os"
	"testing"
)

func TestVersion(t *testing.T) {
	v := tReadVersion(t)
	if v != Version {
		t.Fatal("invalid version")
	}
}

func tReadVersion(t *testing.T) string {
	d, err := os.ReadFile("../../waroot/VERSION")
	if err != nil {
		t.Fatal(err)
	}
	return string(bytes.TrimSpace(d))
}
