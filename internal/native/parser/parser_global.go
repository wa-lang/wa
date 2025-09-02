// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package parser

// global $f64: f64 = 12.34567
// global $name = "wa native assembly language"

// global $f32: 20 = f32(12.5)

// global $info: 1024 = {
//     5: "abc",    # 从第5字节开始 `abc\0`
//     9: i32(123), # 从第9字节开始
// }

func (p *parser) parseGlobal() {
	p.next()

	// TODO
}
