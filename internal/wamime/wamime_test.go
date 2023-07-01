// 版权 @2023 凹语言 作者。保留所有权利。

package wamime

import (
	"testing"
)

func TestGetCodeMime(t *testing.T) {
	for i, tx := range tests {
		got := GetCodeMime(tx.filename, []byte(tx.code))
		expect := tx.mime
		if got != expect {
			t.Fatalf("%d: expect =%q, got = %q", i, expect, got)
		}
	}
}

var tests = []struct {
	filename string
	code     string
	mime     string
}{
	{"-", "", ""},
	{"prog.wa", "", "wa"},
	{"prog.wz", "", "wz"},

	{"x.wa", "#", "wa"},
	{"", "#wa:syntax=wx", "wx"},

	{
		"",
		`// 版权 @2019 凹语言 作者。保留所有权利。

#wa:syntax=wa

import "fmt"
import "runtime"

global year: i32 = 2023

func main {
	println("你好，凹语言！", runtime.WAOS)
	println(add(40, 2), year)

	fmt.Println("1+1 =", 1+1)
}

func add(a: i32, b: i32) => i32 {
	return a+b
}

`,
		"wa",
	},

	{
		"",
		`// 版权 @2019 凹语言 作者。保留所有权利。

#wa:syntax=wz

引于 "书"

【启】：
  书·说："你好，凹语言中文版！"
。
`,
		"wz",
	},
}
