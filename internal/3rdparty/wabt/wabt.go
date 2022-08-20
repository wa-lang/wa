// 版权 @2022 凹语言 作者。保留所有权利。

//go:build cgo
// +build cgo

package wabt

/*
#cgo LDFLAGS:
#cgo CPPFLAGS: -I./internal/wabt-1.0.29
#cgo CXXFLAGS: -std=c++17
*/
import "C"

// TDOO
