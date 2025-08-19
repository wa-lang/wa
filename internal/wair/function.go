// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wair

/**************************************
本文件包含了 function 对象的功能
**************************************/

//-------------------------------------

// 开始
func (f *Function) StartBody() {
	f.Body = &Block{}
	f.Body.Init()
}
