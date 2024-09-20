// 版权 @2024 凹语言 作者。保留所有权利。

package watutil

import "wa-lang.org/wa/internal/wat/parser"

func WatStrip(path string, src []byte) (watBytes []byte, err error) {
	data, err := readSource(path, src)
	if err != nil {
		return nil, err
	}

	m, err := parser.ParseModule(path, src)
	if err != nil {
		return nil, err
	}

	// 删除未使用对象
	m = new_RemoveUnusedPass(m).DoPass()
	_ = m

	// TODO: m => string

	return data, nil
}
