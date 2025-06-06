# 版权 @2023 凹语言 作者。保留所有权利。

default:

ci-test-all:
	@echo "== fmt examples/... =="
	go run ../../main.go fmt ./...

	@echo "== examples test begin =="

	@echo "# loop forever"

	@echo "#cd ./brainfuck && make"
	@echo "#cd ./expr && make"
	@echo "#cd ./hello && make"
	@echo "#cd ./misc && make"
	@echo "#cd ./pkg && make"
	@echo "#cd ./prime && make"
	@echo "#cd ./reftoptr && make"
	@echo "#cd ./surface && make"

	@echo "# snake"
	@echo "# cd ./snake && make publish"

	go run ../../main.go run copy.wa
	go run ../../main.go run eq.wa
	go run ../../main.go run interface_named.wa
	go run ../../main.go run short-var.wa
	go run ../../main.go run strbytes.wa
	go run ../../main.go run struct.wa

	@echo "== examples ok =="

clean:

