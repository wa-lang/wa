# Copyright (C) 2026 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

WA:=go run ../../../../../main.go

ast:
	$(WA) ast -arch=loong64 native-env-linux-la64.wa.s > env-wa.out.s

ast-zh:
	$(WA) ast -arch=loong64 native-env-linux-la64.wz.s > env-wz.out.s

clean:
