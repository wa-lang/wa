# wat2c 翻译后C代码性能对比

注意，这个fib时递归实现，执行性能比较低。

`fib(46)` 递归实现, `-O1`优化：

```
$ make wat2c
go build -o fib_go_native.exe fib_go_native.go
clang -O0 -o fib_c_native_O0.exe _fib_c_native.c
wa wat2c -o fib_wat2c_native.c fib_wat.txt && clang -O0 -o fib_wat2c_native_O0.exe fib_wat2c_main.c
clang -O1 -o fib_c_native_O1.exe _fib_c_native.c
wa wat2c -o fib_wat2c_native.c fib_wat.txt && clang -O1 -o fib_wat2c_native_O1.exe fib_wat2c_main.c
time ./fib_go_native.exe
fib(46) = 1836311903
        8.31 real         7.97 user         0.05 sys
time ./fib_c_native_O0.exe
fib(46) = 1836311903
        9.53 real         9.44 user         0.02 sys
time ./fib_wat2c_native_O0.exe
fib(46) = 1836311903
       30.80 real        30.56 user         0.06 sys
time ./fib_c_native_O1.exe
fib(46) = 1836311903
        5.08 real         4.92 user         0.01 sys
time ./fib_wat2c_native_O1.exe
fib(46) = 1836311903
        4.85 real         4.81 user         0.01 sys
```

wat转译到C代码在`-O1`优化的执行性能和本地C版本持平, 比Go快20-30%.
