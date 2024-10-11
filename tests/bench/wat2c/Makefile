# 版权 @2024 凹语言 作者。保留所有权利。

wat2c:
	go build -o fib_go_native.exe fib_go_native.go
	
	clang -O0 -o fib_c_native_O0.exe _fib_c_native.c
	wa wat2c -o fib_wat2c_native.c fib_wat.txt && clang -O0 -o fib_wat2c_native_O0.exe fib_wat2c_main.c

	clang -O1 -o fib_c_native_O1.exe _fib_c_native.c
	wa wat2c -o fib_wat2c_native.c fib_wat.txt && clang -O1 -o fib_wat2c_native_O1.exe fib_wat2c_main.c

	clang -O3 -o fib_c_native_O3.exe _fib_c_native.c
	wa wat2c -o fib_wat2c_native.c fib_wat.txt && clang -O3 -o fib_wat2c_native_O3.exe fib_wat2c_main.c

	time ./fib_go_native.exe
	
	time ./fib_c_native_O0.exe
	time ./fib_wat2c_native_O0.exe

	time ./fib_c_native_O1.exe
	time ./fib_wat2c_native_O1.exe

	time ./fib_c_native_O3.exe
	time ./fib_wat2c_native_O3.exe

clean:
