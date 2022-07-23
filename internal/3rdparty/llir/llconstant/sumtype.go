package llconstant

// === [ constant.Constant ] ===================================================

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*Int) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*Float) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*Null) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*NoneToken) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*Struct) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*Array) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*CharArray) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*Vector) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*ZeroInitializer) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*Undef) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*Poison) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*BlockAddress) IsConstant() {}

// --- [ Unary expressions ] ---------------------------------------------------

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*ExprFNeg) IsConstant() {}

// --- [ Binary expressions ] --------------------------------------------------

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*ExprAdd) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*ExprFAdd) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*ExprSub) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*ExprFSub) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*ExprMul) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*ExprFMul) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*ExprUDiv) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*ExprSDiv) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*ExprFDiv) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*ExprURem) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*ExprSRem) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*ExprFRem) IsConstant() {}

// --- [ Bitwise expressions ] -------------------------------------------------

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*ExprShl) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*ExprLShr) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*ExprAShr) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*ExprAnd) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*ExprOr) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*ExprXor) IsConstant() {}

// --- [ Vector expressions ] --------------------------------------------------

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*ExprExtractElement) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*ExprInsertElement) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*ExprShuffleVector) IsConstant() {}

// --- [ Aggregate expressions ] -----------------------------------------------

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*ExprExtractValue) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*ExprInsertValue) IsConstant() {}

// --- [ Memory expressions ] --------------------------------------------------

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*ExprGetElementPtr) IsConstant() {}

// --- [ Conversion expressions ] ----------------------------------------------

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*ExprTrunc) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*ExprZExt) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*ExprSExt) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*ExprFPTrunc) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*ExprFPExt) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*ExprFPToUI) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*ExprFPToSI) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*ExprUIToFP) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*ExprSIToFP) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*ExprPtrToInt) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*ExprIntToPtr) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*ExprBitCast) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*ExprAddrSpaceCast) IsConstant() {}

// --- [ Other expressions ] ---------------------------------------------------

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*ExprICmp) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*ExprFCmp) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*ExprSelect) IsConstant() {}

// === [ constant.Expression ] =================================================

// --- [ Unary expressions ] ---------------------------------------------------

// IsExpression ensures that only constants expressions can be assigned to the
// constant.Expression interface.
func (*ExprFNeg) IsExpression() {}

// --- [ Binary expressions ] --------------------------------------------------

// IsExpression ensures that only constants expressions can be assigned to the
// constant.Expression interface.
func (*ExprAdd) IsExpression() {}

// IsExpression ensures that only constants expressions can be assigned to the
// constant.Expression interface.
func (*ExprFAdd) IsExpression() {}

// IsExpression ensures that only constants expressions can be assigned to the
// constant.Expression interface.
func (*ExprSub) IsExpression() {}

// IsExpression ensures that only constants expressions can be assigned to the
// constant.Expression interface.
func (*ExprFSub) IsExpression() {}

// IsExpression ensures that only constants expressions can be assigned to the
// constant.Expression interface.
func (*ExprMul) IsExpression() {}

// IsExpression ensures that only constants expressions can be assigned to the
// constant.Expression interface.
func (*ExprFMul) IsExpression() {}

// IsExpression ensures that only constants expressions can be assigned to the
// constant.Expression interface.
func (*ExprUDiv) IsExpression() {}

// IsExpression ensures that only constants expressions can be assigned to the
// constant.Expression interface.
func (*ExprSDiv) IsExpression() {}

// IsExpression ensures that only constants expressions can be assigned to the
// constant.Expression interface.
func (*ExprFDiv) IsExpression() {}

// IsExpression ensures that only constants expressions can be assigned to the
// constant.Expression interface.
func (*ExprURem) IsExpression() {}

// IsExpression ensures that only constants expressions can be assigned to the
// constant.Expression interface.
func (*ExprSRem) IsExpression() {}

// IsExpression ensures that only constants expressions can be assigned to the
// constant.Expression interface.
func (*ExprFRem) IsExpression() {}

// --- [ Bitwise expressions ] -------------------------------------------------

// IsExpression ensures that only constants expressions can be assigned to the
// constant.Expression interface.
func (*ExprShl) IsExpression() {}

// IsExpression ensures that only constants expressions can be assigned to the
// constant.Expression interface.
func (*ExprLShr) IsExpression() {}

// IsExpression ensures that only constants expressions can be assigned to the
// constant.Expression interface.
func (*ExprAShr) IsExpression() {}

// IsExpression ensures that only constants expressions can be assigned to the
// constant.Expression interface.
func (*ExprAnd) IsExpression() {}

// IsExpression ensures that only constants expressions can be assigned to the
// constant.Expression interface.
func (*ExprOr) IsExpression() {}

// IsExpression ensures that only constants expressions can be assigned to the
// constant.Expression interface.
func (*ExprXor) IsExpression() {}

// --- [ Vector expressions ] --------------------------------------------------

// IsExpression ensures that only constants expressions can be assigned to the
// constant.Expression interface.
func (*ExprExtractElement) IsExpression() {}

// IsExpression ensures that only constants expressions can be assigned to the
// constant.Expression interface.
func (*ExprInsertElement) IsExpression() {}

// IsExpression ensures that only constants expressions can be assigned to the
// constant.Expression interface.
func (*ExprShuffleVector) IsExpression() {}

// --- [ Aggregate expressions ] -----------------------------------------------

// IsExpression ensures that only constants expressions can be assigned to the
// constant.Expression interface.
func (*ExprExtractValue) IsExpression() {}

// IsExpression ensures that only constants expressions can be assigned to the
// constant.Expression interface.
func (*ExprInsertValue) IsExpression() {}

// --- [ Memory expressions ] --------------------------------------------------

// IsExpression ensures that only constants expressions can be assigned to the
// constant.Expression interface.
func (*ExprGetElementPtr) IsExpression() {}

// --- [ Conversion expressions ] ----------------------------------------------

// IsExpression ensures that only constants expressions can be assigned to the
// constant.Expression interface.
func (*ExprTrunc) IsExpression() {}

// IsExpression ensures that only constants expressions can be assigned to the
// constant.Expression interface.
func (*ExprZExt) IsExpression() {}

// IsExpression ensures that only constants expressions can be assigned to the
// constant.Expression interface.
func (*ExprSExt) IsExpression() {}

// IsExpression ensures that only constants expressions can be assigned to the
// constant.Expression interface.
func (*ExprFPTrunc) IsExpression() {}

// IsExpression ensures that only constants expressions can be assigned to the
// constant.Expression interface.
func (*ExprFPExt) IsExpression() {}

// IsExpression ensures that only constants expressions can be assigned to the
// constant.Expression interface.
func (*ExprFPToUI) IsExpression() {}

// IsExpression ensures that only constants expressions can be assigned to the
// constant.Expression interface.
func (*ExprFPToSI) IsExpression() {}

// IsExpression ensures that only constants expressions can be assigned to the
// constant.Expression interface.
func (*ExprUIToFP) IsExpression() {}

// IsExpression ensures that only constants expressions can be assigned to the
// constant.Expression interface.
func (*ExprSIToFP) IsExpression() {}

// IsExpression ensures that only constants expressions can be assigned to the
// constant.Expression interface.
func (*ExprPtrToInt) IsExpression() {}

// IsExpression ensures that only constants expressions can be assigned to the
// constant.Expression interface.
func (*ExprIntToPtr) IsExpression() {}

// IsExpression ensures that only constants expressions can be assigned to the
// constant.Expression interface.
func (*ExprBitCast) IsExpression() {}

// IsExpression ensures that only constants expressions can be assigned to the
// constant.Expression interface.
func (*ExprAddrSpaceCast) IsExpression() {}

// --- [ Other expressions ] ---------------------------------------------------

// IsExpression ensures that only constants expressions can be assigned to the
// constant.Expression interface.
func (*ExprICmp) IsExpression() {}

// IsExpression ensures that only constants expressions can be assigned to the
// constant.Expression interface.
func (*ExprFCmp) IsExpression() {}

// IsExpression ensures that only constants expressions can be assigned to the
// constant.Expression interface.
func (*ExprSelect) IsExpression() {}
