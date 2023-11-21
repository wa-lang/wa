# 版权 @2023 凹语言 作者。保留所有权利。

default:

ci-test-all:
	@echo "== fmt examples/... =="
	go run ../../main.go fmt ./...

	@echo "== examples test begin =="

	# loop forever

	#cd ./brainfuck && make
	#cd ./expr && make
	#cd ./hello && make
	#cd ./misc && make
	#cd ./pkg && make
	#cd ./prime && make
	#cd ./reftoptr && make
	#cd ./surface && make

	# snake
	cd ./snake && make publish

	go run ../../main.go run copy.wa
	go run ../../main.go run eq.wa
	go run ../../main.go run interface_named.wa
	go run ../../main.go run short-var.wa
	go run ../../main.go run strbytes.wa
	go run ../../main.go run struct.wa

	@echo "== examples ok =="

clean:

