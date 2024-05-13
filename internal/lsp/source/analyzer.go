// 版权 @2024 凹语言 作者。保留所有权利。

package source

import "wa-lang.org/wa/internal/lsp/protocol"

// Analyzer represents a go/analysis analyzer with some boolean properties
// that let the user know how to use the analyzer.
type Analyzer struct {
	// Analyzer *analysis.Analyzer

	// Enabled reports whether the analyzer is enabled. This value can be
	// configured per-analysis in user settings. For staticcheck analyzers,
	// the value of the Staticcheck setting overrides this field.
	Enabled bool

	// Fix is the name of the suggested fix name used to invoke the suggested
	// fixes for the analyzer. It is non-empty if we expect this analyzer to
	// provide its fix separately from its diagnostics. That is, we should apply
	// the analyzer's suggested fixes through a Command, not a TextEdit.
	Fix string

	// ActionKind is the kind of code action this analyzer produces. If
	// unspecified the type defaults to quickfix.
	ActionKind []protocol.CodeActionKind

	// Severity is the severity set for diagnostics reported by this
	// analyzer. If left unset it defaults to Warning.
	Severity protocol.DiagnosticSeverity
}
