# 版权 @2023 凹语言 作者。保留所有权利。

default:
	go run ../../main.go yacc -l -p=expr -c="copyright.txt" -o="y.wa" expr.y
	go run ../../main.go y.wa

dev:
	go run ../../main.go yacc -l -p=expr -c="copyright.txt" -o="y.wa.go" expr.y
	go run y.wa.go

clean:
