# 版权 @2023 凹语言 作者。保留所有权利。

default:

ci-test-all:
	@echo "== fmt examples/... =="
	go run ../../main.go fmt ./...

	@echo "== examples test begin =="

	# loop forever
	# cd ./arduino && make
	# cd ./arduino-wat && make

	cd ./brainfuck && make
	cd ./expr && make
	cd ./hello && make
	cd ./misc && make
	cd ./pkg && make
	cd ./prime && make
	cd ./reftoptr && make

	# snake
	cd ./snake && make publish

	go run ../../main.go copy.wa
	go run ../../main.go eq.wa
	go run ../../main.go interface_named.wa
	go run ../../main.go short-var.wa
	go run ../../main.go strbytes.wa
	go run ../../main.go struct.wa

	@echo "== examples ok =="

clean:

