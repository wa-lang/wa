default:
	clang -S -emit-llvm _builtin.c
	clang -Woverride-module _builtin.ll main.ll
	./a.out || echo $$?

builtin-ll:
	clang -S -emit-llvm _builtin.c

clean:
