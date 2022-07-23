package llmetadata_test

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/llir/llvm/asm"
	"github.com/llir/llvm/internal/osutil"
)

func TestModule(t *testing.T) {
	golden := []struct {
		path string
	}{
		// Specialized metadata nodes.
		{path: "../../testdata/llvm/test/Assembler/debug-info.ll"},
		// Coreutils.
		{path: "../../testdata/coreutils/test/cat.ll"},
	}
	hasTestdata := osutil.Exists("../../testdata/llvm")
	for _, g := range golden {
		if filepath.HasPrefix(g.path, "../../testdata") && !hasTestdata {
			// Skip test cases from the llir/testdata submodule if not downloaded.
			// Users may add this submodule using git clone --recursive.
			continue
		}
		log.Printf("=== [ %s ] ===", g.path)
		m, err := asm.ParseFile(g.path)
		if err != nil {
			t.Errorf("unable to parse %q into AST; %+v", g.path, err)
			continue
		}
		path := g.path
		if osutil.Exists(g.path + ".golden") {
			path = g.path + ".golden"
		}
		buf, err := ioutil.ReadFile(path)
		if err != nil {
			t.Errorf("unable to read %q; %+v", path, err)
			continue
		}
		want := string(buf)
		got := m.String()
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("module %q mismatch (-want +got):\n%s", path, diff)
			continue
		}
	}
}
