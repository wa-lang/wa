# 版权 @2024 凹语言 作者。保留所有权利。

wat2c:
	go build -o fib_go_native.exe fib_go_native.go
	clang -O1 -o fib_c_native.exe _fib_c_native.c
	wa wat2c -o fib_wat2c_native.c fib_wat.txt && clang -O1 -o fib_wat2c_native.exe fib_wat2c_main.c

	time ./fib_go_native.exe
	time ./fib_c_native.exe
	time ./fib_wat2c_native.exe

clean:
