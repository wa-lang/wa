// 版权 @2024 凹语言 作者。保留所有权利。

package loader

import "os"

func dirPathExists(path string) bool {
	fi, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return fi.IsDir()
}
