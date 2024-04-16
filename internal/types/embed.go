// 版权 @2024 凹语言 作者。保留所有权利。

package types

// #wa:embed path
// #wa:embed path limit
// #wa:embed path offset limit
// const s = "if-empty"
type embedInfo struct {
	FilePath string
	Offset   int
	Limit    int
}
