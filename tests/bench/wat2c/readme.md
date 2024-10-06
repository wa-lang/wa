# wat2c 翻译后C代码性能对比

注意，这个fib时递归实现，执行性能比较低。

`fib(46)` 递归实现, `-O1`优化：

```
$ make wat2c
go build -o fib_go_native.exe fib_go_native.go
clang -O1 -o fib_c_native.exe _fib_c_native.c
wa wat2c -o fib_wat2c_native.c fib_wat.txt && clang -O1 -o fib_wat2c_native.exe fib_wat2c_main.c
time ./fib_go_native.exe
1836311903
        8.81 real         8.00 user         0.06 sys
time ./fib_c_native.exe
fib(46) = 1836311903
        6.50 real         5.24 user         0.05 sys
time ./fib_wat2c_native.exe
fib(46) = 1836311903
        6.79 real         5.11 user         0.05 sys
```

wat转译到C代码的执行性能和本地C版本持平, 比Go快30%.
