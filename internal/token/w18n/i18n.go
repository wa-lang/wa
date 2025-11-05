// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package w18n

// 凹语言国际化
// 基于wz映射, 对应 *.$Local.wx 后缀名
type W18n struct {
	wx Local
}

func New(m *Local) *W18n {
	return &W18n{wx: *m}
}

func (p *W18n) Name() string {
	return p.wx.Local
}
