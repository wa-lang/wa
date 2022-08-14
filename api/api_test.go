// 版权 @2022 凹语言 作者。保留所有权利。

package api_test

import (
	"fmt"
	"os/exec"
	"runtime"
	"testing"
	"testing/fstest"

	"github.com/wa-lang/wa/api"
	"github.com/wa-lang/wa/internal/config"
)

var tClangPath string

func init() {
	clangExe := "clang"
	if runtime.GOOS == "windows" {
		clangExe += ".exe"
	}
	tClangPath, _ = exec.LookPath(clangExe)
}

func _TestRunVFS(t *testing.T) {
	if tClangPath == "" {
		t.Skip("clang not found")
	}

	var expect = "你好, 凹语言!"
	var code = fmt.Sprintf(`fn main() { println("%s") }`, expect)

	vfs := &config.PkgVFS{
		App: fstest.MapFS{
			"src/main.wa": &fstest.MapFile{
				Data: []byte(code),
			},
		},
	}

	output, err := api.RunVFS(vfs, "myapp")
	if err != nil {
		t.Fatal(err)
	}

	if got := string(output); got != expect {
		t.Fatalf("got %q, want %q", got, expect)
	}
}
