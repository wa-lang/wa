package ast

// A ParserMode value is a set of flags (or 0).
// They control the amount of source code parsed and other optional
// parser functionality.
type ParserMode uint

const (
	PackageClauseOnly ParserMode       = 1 << iota // stop parsing after package clause
	ImportsOnly                                    // stop parsing after import declarations
	ParseComments                                  // parse comments and add them to AST
	Trace                                          // print a trace of parsed productions
	DeclarationErrors                              // report declaration errors
	SpuriousErrors                                 // same as AllErrors, for backward-compatibility
	AllErrors         = SpuriousErrors             // report all errors (not just the first 10 on different lines)
)
