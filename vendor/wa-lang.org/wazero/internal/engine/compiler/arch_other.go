//go:build !amd64 && !arm64

package compiler

import (
	"fmt"
	"runtime"

	"wa-lang.org/wazero/internal/wazeroir"
)

// archContext is empty on an unsupported architecture.
type archContext struct{}

// newCompiler returns an unsupported error.
func newCompiler(ir *wazeroir.CompilationResult, _ bool) (compiler, error) {
	return nil, fmt.Errorf("unsupported GOARCH %s", runtime.GOARCH)
}
