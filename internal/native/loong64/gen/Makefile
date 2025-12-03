# Copyright (C) 2025 武汉凹语言科技有限公司
# SPDX-License-Identifier: AGPL-3.0-or-later

# 生成龙芯英文版指令集表格
# 1. 下载英文版到当前目录 https://loongson.github.io/LoongArch-Documentation/LoongArch-Vol1-EN.pdf
# 2. 先执行 `go run a_gen.go`, 生成 `a_out.go` 文件
# 3. 再执行 `go run . -pkg=loong64 -out=../a_out.go`, 生成 `../a_out.go` 文件

default:
	go run . -pkg=loong64 -out=../a_out.go

gen2:
	go run .

gen1:
	go run a_gen.go
	go fmt ./...

clean:
