// 版权 @2024 凹语言 作者。保留所有权利。

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
