package llconstant

// Assert that each constant implements the constant.Constant interface.
var (
	// Constant expressions.
	_ Constant = Expression(nil)

	// Constants.
	_ Constant = (*Int)(nil)
	_ Constant = (*Float)(nil)
	_ Constant = (*Null)(nil)
	_ Constant = (*NoneToken)(nil)
	_ Constant = (*Struct)(nil)
	_ Constant = (*Array)(nil)
	_ Constant = (*CharArray)(nil)
	_ Constant = (*Vector)(nil)
	_ Constant = (*ZeroInitializer)(nil)
	_ Constant = (*Undef)(nil)
	_ Constant = (*Poison)(nil)
	_ Constant = (*BlockAddress)(nil)
)

// Assert that each constant expression implements the constant.Expression interface.
var (
	// Binary expressions.
	_ Expression = (*ExprAdd)(nil)
	_ Expression = (*ExprFAdd)(nil)
	_ Expression = (*ExprSub)(nil)
	_ Expression = (*ExprFSub)(nil)
	_ Expression = (*ExprMul)(nil)
	_ Expression = (*ExprFMul)(nil)
	_ Expression = (*ExprUDiv)(nil)
	_ Expression = (*ExprSDiv)(nil)
	_ Expression = (*ExprFDiv)(nil)
	_ Expression = (*ExprURem)(nil)
	_ Expression = (*ExprSRem)(nil)
	_ Expression = (*ExprFRem)(nil)
	// Bitwise expressions.
	_ Expression = (*ExprShl)(nil)
	_ Expression = (*ExprLShr)(nil)
	_ Expression = (*ExprAShr)(nil)
	_ Expression = (*ExprAnd)(nil)
	_ Expression = (*ExprOr)(nil)
	_ Expression = (*ExprXor)(nil)
	// Vector expressions.
	_ Expression = (*ExprExtractElement)(nil)
	_ Expression = (*ExprInsertElement)(nil)
	_ Expression = (*ExprShuffleVector)(nil)
	// Aggregate expressions.
	_ Expression = (*ExprExtractValue)(nil)
	_ Expression = (*ExprInsertValue)(nil)
	// Memory expressions.
	_ Expression = (*ExprGetElementPtr)(nil)
	// Conversion expressions.
	_ Expression = (*ExprTrunc)(nil)
	_ Expression = (*ExprZExt)(nil)
	_ Expression = (*ExprSExt)(nil)
	_ Expression = (*ExprFPTrunc)(nil)
	_ Expression = (*ExprFPExt)(nil)
	_ Expression = (*ExprFPToUI)(nil)
	_ Expression = (*ExprFPToSI)(nil)
	_ Expression = (*ExprUIToFP)(nil)
	_ Expression = (*ExprSIToFP)(nil)
	_ Expression = (*ExprPtrToInt)(nil)
	_ Expression = (*ExprIntToPtr)(nil)
	_ Expression = (*ExprBitCast)(nil)
	_ Expression = (*ExprAddrSpaceCast)(nil)
	// Other expressions.
	_ Expression = (*ExprICmp)(nil)
	_ Expression = (*ExprFCmp)(nil)
	_ Expression = (*ExprSelect)(nil)
)
