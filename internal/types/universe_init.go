// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package types

// 因为 wa 和 wz 的 error 接口有相互依赖, 必须控制初始化的顺序

func init() {
	initWa()
	initWz()

	// error/错误 双方可见, 这是中英文技术平权必须的
	waDef(wzUniverseError)
	wzDef(waUniverseError)
}
